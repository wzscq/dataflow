package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"github.com/dop251/goja"
	"github.com/rs/xid"
)

type groupTransformConf struct {
	FuncScript funcScript `json:"funcScript"`
}

type nodeExecutorGroupTransform struct {
	NodeConf node
	TransformFunc goja.Callable
	JSRuntime *goja.Runtime
}

func (nodeExecutor *nodeExecutorGroupTransform)getBatchNumber()(string){
	return xid.New().String()
}

func (nodeExecutor *nodeExecutorGroupTransform)convertToFlowDataItem(
	data interface{})(*flowDataItem,*common.CommonError){
	jsonStr, err := json.Marshal(data)
    if err != nil {
        log.Println(err)
		return nil,nil
    }
	//log.Println(string(jsonStr))
	dataItem:=&flowDataItem{}
    if err := json.Unmarshal(jsonStr, dataItem); err != nil {
        log.Println(err)
		return nil,nil
    }

	return dataItem,nil
}

func (nodeExecutor *nodeExecutorGroupTransform)getNodeConf()(*groupTransformConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	conf:=&groupTransformConf{}
    if err := json.Unmarshal(jsonStr, conf); err != nil {
        log.Println(err)
		return nil
    }

	return conf
}


func (nodeExecutor *nodeExecutorGroupTransform)createTransformFunction(name,body string)(goja.Callable,*common.CommonError){
	funcStr:="function "+name+"(groupItem){"+body+"}"
	_, err:=nodeExecutor.JSRuntime.RunString(funcStr)
	if err!=nil {
		log.Println(err)
		params:=map[string]interface{}{
			"errorMessage":err.Error(),
		}
		return nil,common.CreateError(common.ResultCreateTransformFunctionError,params)
	}
	
	//
	transFunc, ok := goja.AssertFunction(nodeExecutor.JSRuntime.Get(name))
	if !ok {
		log.Println("AssertFunction return false")
		params:=map[string]interface{}{
			"errorMessage":"AssertFunction return false",
		}
		return nil,common.CreateError(common.ResultCreateTransformFunctionError,params)
	}

	return transFunc,nil
}

func (nodeExecutor *nodeExecutorGroupTransform)getTransformFunction(
	transformCfg *groupTransformConf)(goja.Callable, *common.CommonError){
	//????????????goja???runtime????????????????????????????????????????????????runtime
	if nodeExecutor.JSRuntime == nil {
		nodeExecutor.JSRuntime=goja.New()
		nodeExecutor.JSRuntime.SetFieldNameMapper(goja.TagFieldNameMapper("json",true))
		nodeExecutor.JSRuntime.Set("g_BatchNumber",nodeExecutor.getBatchNumber())
		nodeExecutor.JSRuntime.Set("g_Index",1)
	}

	if nodeExecutor.TransformFunc== nil {
		funcName:="groupTransform"
		transformFunc,err:=nodeExecutor.createTransformFunction(funcName,transformCfg.FuncScript.Content)
		if err!=nil {
			return nil,err
		}
		nodeExecutor.TransformFunc=transformFunc
	}

	return nodeExecutor.TransformFunc,nil
}

func (nodeExecutor *nodeExecutorGroupTransform)Transform(
	nodeCfg *groupTransformConf,
	dataItem *flowDataItem)(*flowDataItem,*common.CommonError){
	
	transFunc,err:=nodeExecutor.getTransformFunction(nodeCfg)
	if err!=nil {
		return nil,err
	}

	groupData:=nodeExecutor.JSRuntime.ToValue(*dataItem)
	res, funcErr := transFunc(goja.Undefined(), groupData)
	if funcErr != nil {
		log.Println(funcErr)
		params:=map[string]interface{}{
			"errorMessage":funcErr.Error(),
		}
		return nil,common.CreateError(common.ResultExecuteTransformFunctionError,params)
	}

	return nodeExecutor.convertToFlowDataItem(res.Export())
}

func (nodeExecutor *nodeExecutorGroupTransform)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorGroupTransform run start")
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_GROUP_TRANSFORM,
	}
	//??????????????????
	nodeCfg:=nodeExecutor.getNodeConf()
	if nodeCfg==nil {
		log.Printf("nodeExecutorGroupTransform run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeConfigError,params)
	}

	if node.Input.Data==nil || len(*node.Input.Data)==0 {
		log.Printf("nodeExecutorGroupTransform run end with empty input data.\n")
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=node.Input
		return node.Input,nil
	}

	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID, 
		UserID:req.UserID,
		AppDB:req.AppDB,
		GoOn:true,
	}

	//?????????????????????????????????????????????????????????????????????????????????
	//??????????????????
	resultData:=make([]flowDataItem,len(*req.Data))
	//????????????????????????????????????????????????????????????????????????????????????????????????????????????
	//????????????????????????
	for index,item:= range (*req.Data) {
		resultItem,err:=nodeExecutor.Transform(nodeCfg,&item)
		if err !=nil {
			err.Params["nodeID"]=node.ID
			err.Params["nodeType"]=NODE_GROUP_TRANSFORM
			return node.Input,err
		}

		if resultItem!=nil {
			resultData[index]=*resultItem
		}
	}
	
	//log.Println("resultData:",resultData)
	flowResult.Data=&resultData

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorGroupTransform run end")
	return flowResult,nil
}