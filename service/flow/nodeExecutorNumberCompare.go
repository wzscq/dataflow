package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"strconv"
)

type numberCompareModel struct {
	ModelID string `json:"modelID"`
	Field string  `json:"field"`
	Side string `json:"side"`
}

type numberCompareConf struct {
	VerfiyID string `json:"verifyID"`
	Models []numberCompareModel `json:"models"`
	Tolerance string `json:"tolerance"`
}

type nodeExecutorNumberCompare struct {
	NodeConf node
}

const (
	MODEL_SIDE_LEFT = "left"
	MODEL_SIDE_RIGHT = "right"
)

func (nodeExecutor *nodeExecutorNumberCompare)getNodeConf()(*numberCompareConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	log.Println(string(jsonStr))
	conf:=&numberCompareConf{}
    if err := json.Unmarshal(jsonStr, conf); err != nil {
        log.Println(err)
		return nil
    }

	return conf
}

func (nodeExecutor *nodeExecutorNumberCompare)getModelValue(
	list *[]map[string]interface{},
	field string)(float64,int){
	//找到对应模型的数据
	var sumVal float64
	for _,row:=range (*list){
		fieldVal,found:=row[field]
		if !found {
			log.Printf("nodeExecutorVerifyValue getModelValue no field: %s!\n", field)
			return sumVal,common.ResultNoMatchField 
		}

		switch fieldVal.(type) {
		case float64:
			fVal, _ := fieldVal.(float64)
			sumVal=sumVal+fVal
		case string:
			sVal, _ := fieldVal.(string)
			fVal, err := strconv.ParseFloat(sVal, 64)
			if err !=nil {
				log.Printf("nodeExecutorVerifyValue getModelValue can not convert value to float64: %s!\n", sVal)
				return sumVal,common.ResultMatchValueToFloat64Error
			}
			sumVal=sumVal+fVal
		default:
			log.Printf("nodeExecutorVerifyValue getModelValue not supported field type: %T!\n", fieldVal)
			return sumVal,common.ResultMatchFieldTypeError
		}
	}
	return sumVal,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorNumberCompare)copyDataItem(
	item,newItem *flowDataItem){
	
	newItem.VerifyResult=append(newItem.VerifyResult,item.VerifyResult...)
	newItem.Models=append(newItem.Models,item.Models...)
}

func (nodeExecutor *nodeExecutorNumberCompare)compare(
	verifyConf *numberCompareConf,
	dataItem *flowDataItem,
	tolerance float64)(int){
	verifyItem:=verifyResultItem{
		VerfiyID:verifyConf.VerfiyID,
		VerfiyType:"numberCompare",
		Result:"0",
		Message:"",
	}
	var leftValue,rightValue float64
	for _,verifyModel:= range (verifyConf.Models) {
		for _,modelDataItem:= range (dataItem.Models) {
			if *modelDataItem.ModelID == verifyModel.ModelID && modelDataItem.List!=nil {
				//获取模型数据
				modelValue,errorCode:=nodeExecutor.getModelValue(modelDataItem.List,verifyModel.Field)
				if errorCode != common.ResultSuccess {
					return errorCode
				}
				if verifyModel.Side == MODEL_SIDE_LEFT {
					leftValue=leftValue+modelValue
				}
				if verifyModel.Side == MODEL_SIDE_RIGHT {
					rightValue=rightValue+modelValue
				}
			}
		}
	}
	log.Printf("nodeExecutorNumberCompare compare left:%.f,right:%.f,tolerance:%.f",leftValue,rightValue,tolerance)
	//判断左右表数据是否相等
	diff:=leftValue - rightValue
	//如果左右表数值差异超过了容差范围，则校验失败
	if diff > tolerance {
		verifyItem.Result = "1"
	}
	
	if diff < tolerance*-1 {
		verifyItem.Result="-1"
	}
	dataItem.VerifyResult=append(dataItem.VerifyResult,verifyItem)
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorNumberCompare)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorNumberCompare run start")
	
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

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_NUMBER_COMPARE,
	}
	//加载节点配置
	nodeConf:=nodeExecutor.getNodeConf()
	if nodeConf==nil {
		log.Printf("nodeExecutorNumberCompare run get node config error\n")
		return flowResult,common.CreateError(common.ResultNodeGroupConfigError,params)
	}

	if node.Input.Data==nil || len(*node.Input.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=node.Input
		return node.Input,nil
	}
	
	tolerance,_:=strconv.ParseFloat(nodeConf.Tolerance,64)
	//执行校验逻辑
	resultData:=make([]flowDataItem,len(*req.Data))
	//这里的分组操作是在数据已经分组的基础上再次分组，分组数据不能跨原来的分组
	//遍历每个数据分组
	for index,item:= range (*req.Data) {
		//复制原始数据
		nodeExecutor.copyDataItem(&item,&(resultData[index]))
		errorCode:=nodeExecutor.compare(nodeConf,&(resultData[index]),tolerance)
		if errorCode != common.ResultSuccess {
			return flowResult,common.CreateError(errorCode,params)
		}
	}

	flowResult.Data=&resultData

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorNumberCompare run end")
	return flowResult,nil
}