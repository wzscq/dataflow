package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"github.com/dop251/goja"
)

type funcScript struct {
	Name string `json:"name"`
	Content string `json:"content"`
}

type dataTransformField struct {
	Field string `json:"field"`
	FuncScript funcScript `json:"funcScript"`
}

type transformModel struct {
	ModelID string `json:"modelID"`
	Fields []dataTransformField `json:"fields"`
}

type dataTransformConf struct {
	OutputType string `json:"outputType"`
	Models []transformModel  `json:"models"`
}

type nodeExecutorDataTransform struct {
	NodeConf node
	//使用modelID、field索引转换函数列表
	TransformFuncs map[string]map[string]goja.Callable
	JSRuntime *goja.Runtime
}

const (
	OUTPUT_TYPE_ALL="all"
	OUTPUT_TYPE_MODIFIED="modified"
)

func (nodeExecutor *nodeExecutorDataTransform)getNodeConf()(*dataTransformConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	log.Println(string(jsonStr))
	conf:=&dataTransformConf{}
    if err := json.Unmarshal(jsonStr, conf); err != nil {
        log.Println(err)
		return nil
    }

	return conf
}

func (nodeExecutor *nodeExecutorDataTransform)copyDataItem(
	item,newItem *flowDataItem){
	newItem.VerifyResult=append(newItem.VerifyResult,item.VerifyResult...)
	newItem.Models=append(newItem.Models,item.Models...)
}

func (nodeExecutor *nodeExecutorDataTransform)createTransformFunction(name,body string)(goja.Callable,*common.CommonError){
	funcStr:="function "+name+"(field, row){"+body+"}"
	_, err:=nodeExecutor.JSRuntime.RunString(funcStr)
	if err!=nil {
		log.Println(err)
		params:=map[string]interface{}{
			"errorMessage":err,
		}
		return nil,common.CreateError(common.ResultCreateTransformFunctionError,params)
	}
	
	//
	fieldFunc, ok := goja.AssertFunction(nodeExecutor.JSRuntime.Get(name))
	if !ok {
		log.Println("AssertFunction return false")
		params:=map[string]interface{}{
			"errorMessage":"AssertFunction return false",
		}
		return nil,common.CreateError(common.ResultCreateTransformFunctionError,params)
	}

	return fieldFunc,nil
}

func (nodeExecutor *nodeExecutorDataTransform)getTransformFunction(
	modelID string,
	transformFieldCfg *dataTransformField)(goja.Callable, *common.CommonError){
	//首先判断goja的runtime是否创建，如果没有创建则创建新的runtime
	if nodeExecutor.JSRuntime == nil {
		nodeExecutor.JSRuntime=goja.New()
	}

	if nodeExecutor.TransformFuncs== nil {
		nodeExecutor.TransformFuncs=map[string]map[string]goja.Callable{}
	}

	//查看当前model、field的函数是否已经存在，如果不存在则创建函数对象
	modelFuncs,ok:=nodeExecutor.TransformFuncs[modelID]
	if !ok {
		modelFuncs=map[string]goja.Callable{}
		nodeExecutor.TransformFuncs[modelID]=modelFuncs
	}

	fieldFunc,ok:=modelFuncs[transformFieldCfg.Field]
	if !ok {
		funcName:=modelID+"_"+transformFieldCfg.Field
		var err *common.CommonError
		fieldFunc,err=nodeExecutor.createTransformFunction(funcName,transformFieldCfg.FuncScript.Content)
		if err!=nil {
			return nil,err
		}
		nodeExecutor.TransformFuncs[modelID][transformFieldCfg.Field]=fieldFunc
	}

	return fieldFunc,nil
}

func (nodeExecutor *nodeExecutorDataTransform)TransformModelField(
	modelID string,
	transformFieldCfg *dataTransformField,
	rowData *map[string]interface{})(*common.CommonError){
	transFunc,err:=nodeExecutor.getTransformFunction(modelID,transformFieldCfg)
	if err!=nil {
		err.Params["modelID"]=modelID
		err.Params["field"]=transformFieldCfg.Field
		return err
	}

	fieldPara:=nodeExecutor.JSRuntime.ToValue(transformFieldCfg.Field)
	rowPara:=nodeExecutor.JSRuntime.ToValue(*rowData)
	res, funcErr := transFunc(goja.Undefined(), fieldPara,rowPara)
	if funcErr != nil {
		log.Println(funcErr)
		params:=map[string]interface{}{
			"modelID":modelID,
			"field":transformFieldCfg.Field,
			"errorMessage":funcErr.Error(),
		}
		return common.CreateError(common.ResultExecuteTransformFunctionError,params)
	}
	
	(*rowData)[transformFieldCfg.Field]=res.Export()

	return nil
}

func (nodeExecutor *nodeExecutorDataTransform)TransformModel(
	transModelCfg *transformModel,
	modelData *modelDataItem,
	outputType string)(*common.CommonError){
	
	log.Println("****************************************************")
	log.Println(outputType)
	var keepFields map[string]interface{}
	//如果只保留修改的列，则去掉未修改的列，这里注意内部使用的字段需要保留，包括：id,version,_save_type
	if outputType==OUTPUT_TYPE_MODIFIED {
			log.Println(keepFields)
		  keepFields=map[string]interface{}{
				"id":"id",
				"version":"version",
				"_save_type":"_save_type",
			}
			for _,transModelField:=range(transModelCfg.Fields){
				keepFields[transModelField.Field]=transModelField.Field
			}
			log.Println(keepFields)
			log.Println("****************************************************")
	}

	for index,_:=range(*modelData.List){
		for _,transModelField:=range(transModelCfg.Fields){
			err:=nodeExecutor.TransformModelField(*modelData.ModelID,&transModelField,&(*modelData.List)[index])
			if err!=nil {
				return err
			}
		}

		log.Println(keepFields)
		if keepFields!=nil {
			rowData:=map[string]interface{}{}
			for field,_:=range(keepFields){
				rowData[field]=(*modelData.List)[index][field]
			}
			(*modelData.List)[index]=rowData
		}
	}

	return nil
}

func (nodeExecutor *nodeExecutorDataTransform)TransformModefied(
	transformConf *dataTransformConf,
	dataItem *flowDataItem)(*flowDataItem,*common.CommonError){
	var modifiedDataItem *flowDataItem
	for _,modelItem:= range (dataItem.Models) {
		for _,transModel:=range(transformConf.Models){
			if *modelItem.ModelID == transModel.ModelID {
				if modelItem.List != nil && len(*modelItem.List)>=0 {
					log.Println("TransformModefied modelid:",transModel.ModelID)
					if modifiedDataItem==nil {
						modifiedDataItem=&flowDataItem{
							Models:[]modelDataItem{},
						}
					}
					log.Println(transformConf.OutputType)
					err:=nodeExecutor.TransformModel(&transModel,&modelItem,transformConf.OutputType)
					if err!=nil {
						return nil,err
					}
					modifiedDataItem.Models=append(modifiedDataItem.Models,modelItem)
					log.Println("TransformModefied modelid:",modifiedDataItem.Models)
				}
			}
		}
	}
	return modifiedDataItem,nil
}

func (nodeExecutor *nodeExecutorDataTransform)Transform(
	transformConf *dataTransformConf,
	dataItem *flowDataItem)(*common.CommonError){
	for index,modelItem:= range (dataItem.Models) {
		for _,transModel:=range(transformConf.Models){
			if *modelItem.ModelID == transModel.ModelID {
				if modelItem.List != nil && len(*modelItem.List)>=0 {
					log.Println(transformConf.OutputType)
					err:=nodeExecutor.TransformModel(&transModel,&(dataItem.Models[index]),transformConf.OutputType)
					if err!=nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (nodeExecutor *nodeExecutorDataTransform)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorDataTransform run start")
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_DATA_TRANSFORM,
	}
	//加载节点配置
	nodeCfg:=nodeExecutor.getNodeConf()
	if nodeCfg==nil {
		log.Printf("nodeExecutorDataTransform run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeConfigError,params)
	}

	if node.Input.Data==nil || len(*node.Input.Data)==0 {
		log.Printf("nodeExecutorDataTransform run end with empty input data.\n")
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=node.Input
		return node.Input,nil
	}

	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID,
		FlowInstanceID:req.FlowInstanceID,
		Stage:req.Stage,
		DebugID:req.DebugID,
		UserRoles:req.UserRoles,
		GlobalFilterData:req.GlobalFilterData,
		UserID:req.UserID,
		AppDB:req.AppDB,
		Token:req.Token,
		FlowConf:req.FlowConf,
		ModelID:req.ModelID,
		ViewID:req.ViewID,
		FilterData:req.FilterData,
		Filter:req.Filter,
		List:req.List,
		Total:req.Total,
		SelectedRowKeys:req.SelectedRowKeys,
		Pagination:req.Pagination,
		Operation:req.Operation,
		SelectAll:req.SelectAll,
		GoOn:true,
	}

	log.Println(nodeCfg.OutputType)
	var resultData []flowDataItem
	if nodeCfg.OutputType==OUTPUT_TYPE_ALL {
		//遍历所有数据模型，对于配置了转换逻辑的模型字段进行处理
		//执行校验逻辑
		resultData=make([]flowDataItem,len(*req.Data))
		//这里的分组操作是在数据已经分组的基础上再次分组，分组数据不能跨原来的分组
		//遍历每个数据分组
		for index,item:= range (*req.Data) {
			//复制原始数据
			nodeExecutor.copyDataItem(&item,&(resultData[index]))
			err:=nodeExecutor.Transform(nodeCfg,&(resultData[index]))
			if err !=nil {
				err.Params["nodeID"]=node.ID
				err.Params["nodeType"]=NODE_DATA_TRANSFORM
				return node.Input,err
			}
		}
	} else {
		resultData=[]flowDataItem{}
		for _,item:= range (*req.Data) {
			dataItem,err:=nodeExecutor.TransformModefied(nodeCfg,&item)
			if err !=nil {
				err.Params["nodeID"]=node.ID
				err.Params["nodeType"]=NODE_DATA_TRANSFORM
				return node.Input,err
			}
			if dataItem!=nil {
				log.Println("append dataItem",dataItem)
				resultData=append(resultData,*dataItem)
			}
		}
	}

	log.Println("resultData:",resultData)
	flowResult.Data=&resultData

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorDataTransform run end")
	return flowResult,nil
}