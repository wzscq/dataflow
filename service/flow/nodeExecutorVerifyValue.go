package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"strconv"
)

const (
	AGGREGATION_AVG="avg"
	AGGREGATION_COUNT="count"
	AGGREGATION_MAX="max"
	AGGREGATION_MIN="min"
)

type verifyValueItem struct {
	VerfiyID string `json:"verifyID"`
	Tolerance string `json:"tolerance"`
	ModelID string `json:"modelID"`
	Field string  `json:"field"`
	Aggregation string `json:"aggregation"`
	Value string `json:"value"`
}

type  verifyValueConf struct {
	Items []verifyValueItem `json:"items"`
}

type nodeExecutorVerifyValue struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorVerifyValue)getNodeConf()(*verifyValueConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println(string(jsonStr))
	conf:=&verifyValueConf{}
  if err := json.Unmarshal(jsonStr, conf); err != nil {
    log.Println(err)
		return nil
  }

	return conf
}

func (nodeExecutor *nodeExecutorVerifyValue)getModelValue(
	list *[]map[string]interface{},
	field,aggregation string)(float64,int){
	var fCount float64
	fCount=float64(len(*list))
	if aggregation == AGGREGATION_COUNT {
		return fCount,common.ResultSuccess
	} 
	//找到对应模型的数据
	var sumVal,maxVal,minVal float64
	for index,row:=range (*list){
		fieldVal,found:=row[field]
		if !found {
			log.Printf("nodeExecutorVerifyValue getModelValue no field: %s!\n", field)
			return sumVal,common.ResultNoMatchField 
		}
		var fVal float64
		switch fieldVal.(type) {
		case float64:
			fVal, _ = fieldVal.(float64)
		case string:
			sVal, _ := fieldVal.(string)
			var err error
			fVal,err= strconv.ParseFloat(sVal, 64)
			if err !=nil {
				log.Printf("nodeExecutorVerifyValue getModelValue can not convert value to float64: %s!\n", sVal)
				return sumVal,common.ResultMatchValueToFloat64Error
			}
		default:
			log.Printf("nodeExecutorVerifyValue getModelValue not supported field type: %T!\n", fieldVal)
			return sumVal,common.ResultMatchFieldTypeError
		}

		sumVal=sumVal+fVal
		if index==0 {
			maxVal=fVal
			minVal=fVal
		} else {
			if maxVal<fVal {
				maxVal=fVal
			}

			if minVal>fVal {
				minVal=fVal
			}
		}
	}

	if aggregation == AGGREGATION_AVG {
		return sumVal/fCount,common.ResultSuccess
	} 

	if aggregation == AGGREGATION_MAX {
		return maxVal,common.ResultSuccess
	}
	
	if aggregation == AGGREGATION_MIN {
		return minVal,common.ResultSuccess
	}

	return sumVal,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorVerifyValue)copyDataItem(
	item,newItem *flowDataItem){
	
	newItem.VerifyResult=item.VerifyResult
	newItem.Models=item.Models
}

func (nodeExecutor *nodeExecutorVerifyValue)getVerifyResult(
	fieldValue,value,tolerance float64)(string){
		diff:=fieldValue-value
	
		if diff<tolerance*-1 {
			return "-1";
		}

		if diff>tolerance {
			return "1";
		}

		return "0"
}

func (nodeExecutor *nodeExecutorVerifyValue)verifyItem(
	verifyConf *verifyValueItem,
	dataItem *flowDataItem)(int){

	var fieldValue float64
	for _,modelDataItem:= range (dataItem.Models) {
		if *modelDataItem.ModelID == verifyConf.ModelID && modelDataItem.List!=nil {
			//获取模型数据
			var errorCode int
			fieldValue,errorCode=nodeExecutor.getModelValue(modelDataItem.List,verifyConf.Field,verifyConf.Aggregation)
			if errorCode != common.ResultSuccess {
				return errorCode
			}
			
			var value float64
			value,_=strconv.ParseFloat(verifyConf.Value,64)
			verifyItem:=verifyResultItem{
				VerfiyID:verifyConf.VerfiyID,
				VerfiyType:"valueVerify",
				Result:"0",
				Message:"",
			}
			//判断左右表数据是否相等
			tolerance,_:=strconv.ParseFloat(verifyConf.Tolerance,64)
			verifyItem.Result=nodeExecutor.getVerifyResult(fieldValue,value,tolerance)
			dataItem.VerifyResult=append(dataItem.VerifyResult,verifyItem)
			break;
		}
	}

	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorVerifyValue)verify(
	verifyConf *verifyValueConf,
	dataItem *flowDataItem)(int){
	for _,verifyItem:=range(verifyConf.Items){
		result:=nodeExecutor.verifyItem(&verifyItem,dataItem)
		if common.ResultSuccess != result {
			return result
		}
	}
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorVerifyValue)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorVerifyValue run start")
	
	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID, 
		UserID:req.UserID,
		AppDB:req.AppDB,
	}

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_VERIFY_VALUE,
	}
	//加载节点配置
	nodeConf:=nodeExecutor.getNodeConf()
	if nodeConf==nil {
		log.Printf("nodeExecutorVerifyValue run get node config error\n")
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if node.Input.Data==nil || len(*node.Input.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=node.Input
		return node.Input,nil
	}
	
	//tolerance,_:=strconv.ParseFloat(nodeConf.Tolerance,64)
	//执行校验逻辑
	resultData:=make([]flowDataItem,len(*req.Data))
	//遍历每个数据分组
	for index,item:= range (*req.Data) {
		//复制原始数据
		nodeExecutor.copyDataItem(&item,&(resultData[index]))
		errorCode:=nodeExecutor.verify(nodeConf,&(resultData[index]))
		if errorCode != common.ResultSuccess {
			return flowResult,common.CreateError(errorCode,params)
		}
	}

	flowResult.Data=&resultData

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorVerifyValue run end")
	return flowResult,nil
}