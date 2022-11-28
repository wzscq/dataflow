package flow

import (
    "time"
	"buoyancyinfo.com/matchflow/common"
	"encoding/json"
	"log"
	"strconv"
)

type splitExtraQuantityModel struct {
	ModelID string `json:"modelID"`
	Field string `json:"field"`
	Side string `json:"side"`
}

type splitExtraQuantityConfig struct {
	Models []splitExtraQuantityModel `json:"models"`
}

type nodeExecutorSplitExtraQuantity struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorSplitExtraQuantity)getNodeConfig()(*splitExtraQuantityConfig){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	log.Println(string(jsonStr))
	conf:=&splitExtraQuantityConfig{}
    if err := json.Unmarshal(jsonStr, conf); err != nil {
        log.Println(err)
		return nil
    }

	return conf
}

func (nodeExecutor *nodeExecutorSplitExtraQuantity)getFieldValue(
	row *map[string]interface{},
	field string)(float64,*common.CommonError){
	params:=map[string]interface{}{
		"field":field,
	}
	
	fieldVal,found:=(*row)[field]
	if !found {
		log.Printf("nodeExecutorSplitExtraQuantity getModelValue no field: %s!\n", field)
		return 0,common.CreateError(common.ResultNoModelField,params) 
	}

	switch fieldVal.(type) {
	case float64:
		fVal, _ := fieldVal.(float64)
		return fVal,nil
	case string:
		sVal, _ := fieldVal.(string)
		fVal, err := strconv.ParseFloat(sVal, 64)
		if err !=nil {
			log.Printf("nodeExecutorSplitExtraQuantity getModelValue can not convert value to float64: %s!\n", sVal)
			return 0,common.CreateError(common.ResultFieldTypeError,params) 
		}
		return fVal,nil
	default:
		log.Printf("nodeExecutorSplitExtraQuantity getModelValue not supported field type: %T!\n", fieldVal)
		return 0,common.CreateError(common.ResultFieldTypeError,params)
	}
}

func (nodeExecutor *nodeExecutorSplitExtraQuantity)getModelValue(
	list *[]map[string]interface{},
	field string)(float64,*common.CommonError){
	//找到对应模型的数据
	var sumVal float64
	for _,row:=range (*list){
		fieldVal,err:=nodeExecutor.getFieldValue(&row,field)
		if err !=nil {
			return sumVal,err
		}
		sumVal=sumVal+fieldVal
	}
	return sumVal,nil
}

func (nodeExecutor *nodeExecutorSplitExtraQuantity)getMinSumVal(
	models []splitExtraQuantityModel,
	modelDatas *[]modelDataItem)(float64,string,*common.CommonError){
	var minVal float64
	var minValModelID string
	modelCount:=0
	for _,modelConf:=range(models) {
		for _,modelData:=range(*modelDatas) {
			if modelConf.ModelID == *modelData.ModelID {
				modelVal,err:=nodeExecutor.getModelValue(modelData.List,modelConf.Field)
				if err != nil {
					return minVal,minValModelID,err
				}
				if modelCount==0 {
					minVal=modelVal
					minValModelID=modelConf.ModelID
				} else {
					if minVal>modelVal {
						minVal=modelVal
						minValModelID=modelConf.ModelID
					}
				}
				modelCount++
			}
		}
	}
	return minVal,minValModelID,nil
}

func (nodeExecutor *nodeExecutorSplitExtraQuantity)splitModelDataItem(
	minList,extraList,modelList *[]map[string]interface{},
	field string,
	minVal float64){
	var sumVal float64
	for index,row:=range(*modelList) {
		fieldVal,_:=nodeExecutor.getFieldValue(&row,field)
		if index==0 || sumVal < minVal {
			sumVal=sumVal+fieldVal
			*minList=append(*minList,row)
		} else {
			*extraList=append(*extraList,row)
		}
	}
}

func (nodeExecutor *nodeExecutorSplitExtraQuantity)splitDataItem(
	itemMinVal,itemExtra *flowDataItem,
	models []splitExtraQuantityModel,
	dataItem *flowDataItem, 
	minVal float64,
	minValModelID string){
	for _,modelConf:=range(models) {
		for _,modelData:=range(dataItem.Models) {
			if modelConf.ModelID == *modelData.ModelID {
				minValList:=[]map[string]interface{}{}
				extraList:=[]map[string]interface{}{}
				if minValModelID == modelConf.ModelID {
					itemMinVal.Models=append(itemMinVal.Models,modelDataItem{
						ModelID:modelData.ModelID,
						List:modelData.List,
					})
				} else {
					nodeExecutor.splitModelDataItem(&minValList,&extraList,modelData.List,modelConf.Field,minVal)
					itemMinVal.Models=append(itemMinVal.Models,modelDataItem{
						ModelID:modelData.ModelID,
						List:&minValList,
					})
					if len(extraList)>0 {
						itemExtra.Models=append(itemExtra.Models,modelDataItem{
							ModelID:modelData.ModelID,
							List:&extraList,
						})
					}
				}
				
			}
		}
	}
}

func (nodeExecutor *nodeExecutorSplitExtraQuantity)split(
	models []splitExtraQuantityModel,
	dataItem *flowDataItem,
	resultData *[]flowDataItem)(*common.CommonError){
	//遍历每个模型，获取字段值汇总的最小值
	minVal,minValModelID,err:=nodeExecutor.getMinSumVal(models,&dataItem.Models)
	if err != nil {
		return err
	}
	//遍历每个模型，按照之前获取的最小值将模型数据拆分到2个组中
	itemMinVal:=flowDataItem{
		Models:[]modelDataItem{},
	}
	itemExtra:=flowDataItem{
		Models:[]modelDataItem{},
	}

	nodeExecutor.splitDataItem(&itemMinVal,&itemExtra,models,dataItem,minVal,minValModelID)

	*resultData=append(*resultData,itemMinVal)
	if len(itemExtra.Models)>0 {
		*resultData=append(*resultData,itemExtra)
	}
	return nil
}

func (nodeExecutor *nodeExecutorSplitExtraQuantity)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorSplitExtraQuantity run start")
	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID, 
		UserID:req.UserID,
		AppDB:req.AppDB,
	}

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_SPLIT_EXTQUANTITY,
	}
	//加载节点配置
	conf:=nodeExecutor.getNodeConfig()
	if conf==nil {
		log.Printf("nodeExecutorSplitExtraQuantity run get node config error\n")
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if node.Input.Data==nil || len(*node.Input.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=node.Input
		return node.Input,nil
	}

	flowData:=[]flowDataItem{}
	//在同一个分组中取左边或右边数据汇总刚好大于另一边所有数据汇总，
	//如果一边有剩余记录则单独放入一个分组中
	//遍历数据分组
	for _,item:= range (*req.Data) {
		err:=nodeExecutor.split(conf.Models,&item,&flowData)
		if err!=nil {
			err.Params["nodeID"]=node.ID
			err.Params["nodeType"]=NODE_SPLIT_EXTQUANTITY
			return flowResult,err
		}
	}

	flowResult.Data=&flowData

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorSplitExtraQuantity run end")
	return flowResult,nil
}