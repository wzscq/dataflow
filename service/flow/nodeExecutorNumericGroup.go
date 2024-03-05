package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"strconv"
)

type numericGroupModel struct {
	ModelID string `json:"modelID"`
	Field string  `json:"field"`
}

type numericGroupConf struct {
	Models []numericGroupModel `json:"models"`
	Tolerance string `json:"tolerance"`
}

type nodeExecutorNumericGroup struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorNumericGroup)getNodeConf()(*numericGroupConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	log.Println(string(jsonStr))
	groupConf:=&numericGroupConf{}
    if err := json.Unmarshal(jsonStr, groupConf); err != nil {
        log.Println(err)
		return nil
    }

	return groupConf
}

func (nodeExecutor *nodeExecutorNumericGroup)getRowKey(
	groupModel *numericGroupModel,
	row map[string]interface{})(float64,int){

	fieldVal,found:=row[groupModel.Field]
	if !found {
		log.Printf("nodeExecutorNumericGroup getRowKey no group key field: %s!\n", groupModel.Field)
		return 0,common.ResultNoKeyFieldForGroup 
	}
		
	switch fieldVal.(type) {
	case string:
		sVal, _ := fieldVal.(string)
		fVal,err:= strconv.ParseFloat(sVal, 64)
		if err!=nil {
			log.Println("nodeExecutorNumericGroup getRowKey key value can not convert to float64")
			log.Println("key value is :"+sVal)
			return 0,common.ResultNotSupportedFieldType
		}
		return fVal,common.ResultSuccess	
	default:
		log.Printf("nodeExecutorNumericGroup getRowKey not supported field type: %T!\n", fieldVal)
		return 0,common.ResultNotSupportedFieldType
	}
}

func (nodeExecutor *nodeExecutorNumericGroup)getGroupDataItem(
	rowKey,tolerance float64,
	modelID string,
	groupResult map[float64]map[string]modelDataItem)(*modelDataItem){
	for key,dataItem:=range(groupResult){
		diff:=key - rowKey
		if diff <0 {
			diff=diff*-1
		}

		if diff <= tolerance && diff >= tolerance*-1 {
			modelData,found:=dataItem[modelID]
			if !found {
				return nil
			}
			return &modelData
		}
	}
	
	return nil
}

func (nodeExecutor *nodeExecutorNumericGroup)addGroupDataItem(
	rowKey float64,
	modelID string,
	row map[string]interface{},
	groupResult map[float64]map[string]modelDataItem){

	//创建一个新的dataItem
	modelData:=modelDataItem{
		ModelID:&modelID,
		List:&[]map[string]interface{}{
			row,
		},
	}

	keyMap,found:=groupResult[rowKey]
	if !found {
		//创建新的key节点
		keyMap:=map[string]modelDataItem{}
		keyMap[modelID]=modelData
		groupResult[rowKey]=keyMap
	} else {
		keyMap[modelID]=modelData
	}
}

func (nodeExecutor *nodeExecutorNumericGroup)addDataItemRow(
	modelData *modelDataItem,
	row map[string]interface{},
){
	*(modelData.List)=append(*(modelData.List),row)
}

func (nodeExecutor *nodeExecutorNumericGroup)groupModel(
	groupModel *numericGroupModel,
	modelData *modelDataItem,
	tolerance float64,
	groupResult map[float64]map[string]modelDataItem)(int){
	
	for _, row := range (*modelData.List) {
		rowKey,errorCode:=nodeExecutor.getRowKey(groupModel,row)
		if errorCode != common.ResultSuccess {
			return errorCode
		}

		modelData:=nodeExecutor.getGroupDataItem(rowKey,tolerance,groupModel.ModelID,groupResult)
		//log.Printf("nodeExecutorGroup groupModel getGroupDateItem rowKey:%s,ModelID:%s \n",rowKey,model.ModelID)
		if modelData == nil {
			//log.Println("nodeExecutorGroup groupModel getGroupDateItem nil")
			nodeExecutor.addGroupDataItem(rowKey,groupModel.ModelID,row,groupResult)
		} else {
			//log.Println("nodeExecutorGroup groupModel addDateItemRow")
			nodeExecutor.addDataItemRow(modelData,row)
		}
	}
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorNumericGroup)group(
	nodeConf *numericGroupConf,
	dataItem *flowDataItem,
	tolerance float64)(*[]flowDataItem,int){
	//对于数值字段的分组处理逻辑，
	//遍历每个model，将指定分组字段的值取出，循环结果map中的所有key，
	//如果这个值和key的差异在容差范围内则将数据分入这个key对应的组中，
	//如果未找到匹配的KEY，则将新的值作为key，创建新的组放入map中
	//map结构 [fieldValueKey][model]modelItem
	groupResult:=map[float64]map[string]modelDataItem{}
	for _, groupModel := range (nodeConf.Models) {
		//注意这里如果没有对应的model不报错，直接跳过处理，
		//主要考虑有些场景可能允许部分表有数据
		for _,modelDataItem:=range (dataItem.Models) {
			if *modelDataItem.ModelID == groupModel.ModelID && modelDataItem.List !=nil && len(*modelDataItem.List)>0 {
				errorCode:=nodeExecutor.groupModel(&groupModel,&modelDataItem,tolerance,groupResult)
				if errorCode!=common.ResultSuccess {
					return nil,errorCode
				}
			}
		}
	}

	//将groupResult map转换为数组放入result中
	result:=[]flowDataItem{}
	for _, group := range groupResult {
		modelDatas:=[]modelDataItem{}
		for _,modelData:= range (group) {
			modelDatas=append(modelDatas,modelData)
		} 
		
		flowDataItem:=flowDataItem{
			VerifyResult:dataItem.VerifyResult,
			Models:modelDatas,
		}
	
		result=append(result,flowDataItem)
	}

	return &result,common.ResultSuccess
}


func (nodeExecutor *nodeExecutorNumericGroup)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorNumericGroup run start")
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_NUMERIC_GROUP,
	}
	//加载节点配置
	nodeConf:=nodeExecutor.getNodeConf()
	if nodeConf==nil {
		log.Printf("nodeExecutorNumericGroup run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeGroupConfigError,params)
	}

	tolerance,_:=strconv.ParseFloat(nodeConf.Tolerance,64)
	
	if node.Input.Data==nil || len(*node.Input.Data)==0 {
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
	flowData:=[]flowDataItem{}
	//按字段分组逻辑，这里不考虑容差，将容差字段分组单独做一个节点类型处理
	//这里的分组操作是在数据已经分组的基础上再次分组，分组数据不能跨原来的分组
	//遍历数据分组
	for _,item:= range (*req.Data) {
		result,errorCode:=nodeExecutor.group(nodeConf,&item,tolerance)
		if errorCode!=common.ResultSuccess {
			return flowResult,common.CreateError(errorCode,params)
		}
		flowData=append(flowData,(*result)...)
	}

	flowResult.Data=&flowData

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorNumericGroup run end")
	return flowResult,nil
}