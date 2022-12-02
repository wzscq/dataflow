package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
)

type groupField struct {
	Field string `json:"field"`
}

type groupModel struct {
	ModelID string `json:"modelID"`
	Fields []groupField  `json:"fields"`
}

type nodeExecutorGroup struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorGroup)getGroupModels()([]groupModel){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok || mapData["models"]==nil {
		return nil
	}
	jsonStr, err := json.Marshal(mapData["models"])
    if err != nil {
        log.Println(err)
		return nil
    }
	models:=[]groupModel{}
    if err := json.Unmarshal(jsonStr, &models); err != nil {
        log.Println(err)
		return nil
    }

	return models
}

func (nodeExecutor *nodeExecutorGroup)getGroupDataItem(
	rowKey,modelID string,
	groupResult map[string]map[string]modelDataItem)(*modelDataItem){
	keyMap,found:=groupResult[rowKey]
	if !found {
		return nil
	}
	
	modelData,found:=keyMap[modelID]
	if !found {
		return nil
	}

	return &modelData
}

func (nodeExecutor *nodeExecutorGroup)getRowKey(
	groupModel *groupModel,
	row map[string]interface{})(string,int){

	rowKey:=""
	for _,field:= range (groupModel.Fields) {
		fieldVal,found:=row[field.Field]
		if !found {
			log.Printf("nodeExecutorGroup getRowKey no group key field: %s!\n", field.Field)
			return "",common.ResultNoKeyFieldForGroup 
		}
		
		switch fieldVal.(type) {
		case string:
			sVal, _ := fieldVal.(string)
			rowKey=rowKey+sVal+"##"
		case nil:
			sVal:="null"
			rowKey=rowKey+sVal+"##"
		default:
			log.Printf("nodeExecutorGroup getRowKey not supported field type: %T!\n", fieldVal)
			return "",common.ResultNotSupportedFieldType
		}
	}

	return rowKey,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorGroup)addGroupDataItem(
	rowKey,modelID string,
	row map[string]interface{},
	groupResult map[string]map[string]modelDataItem){

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

func (nodeExecutor *nodeExecutorGroup)addDataItemRow(
	modelData *modelDataItem,
	row map[string]interface{},
){
	*(modelData.List)=append(*(modelData.List),row)
}

func (nodeExecutor *nodeExecutorGroup)groupModel(
	groupModel *groupModel,
	modelData *modelDataItem,
	groupResult map[string]map[string]modelDataItem)(int){
	
	for _, row := range (*modelData.List) {
		rowKey,errorCode:=nodeExecutor.getRowKey(groupModel,row)
		if errorCode != common.ResultSuccess {
			return errorCode
		}

		modelData:=nodeExecutor.getGroupDataItem(rowKey,groupModel.ModelID,groupResult)
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

func (nodeExecutor *nodeExecutorGroup)group(
	groupModels []groupModel,
	item *flowDataItem)(*[]flowDataItem,int){

	//遍历每个model，将字段值拼接后作为key放入map中
	//map结构 [fieldsValueKey][model]modelItem
	groupResult:=map[string]map[string]modelDataItem{}
	for _, groupModel := range (groupModels) {
		//注意这里如果没有对应的model不报错，直接跳过处理，
		//主要考虑有些场景可能允许部分表有数据
		for _,modelDataItem:=range ((*item).Models) {
			if *modelDataItem.ModelID == groupModel.ModelID && modelDataItem.List !=nil && len(*modelDataItem.List)>0 {
				errorCode:=nodeExecutor.groupModel(&groupModel,&modelDataItem,groupResult)
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
			VerifyResult:item.VerifyResult,
			Models:modelDatas,
		}
	
		result=append(result,flowDataItem)
	}

	return &result,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorGroup)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorGroup run start")
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_GROUP,
	}
	//加载节点配置
	models:=nodeExecutor.getGroupModels()
	if models==nil {
		log.Printf("nodeExecutorGroup run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeGroupConfigError,params)
	}

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
		UserID:req.UserID,
		AppDB:req.AppDB,
	}
	flowData:=[]flowDataItem{}
	//按字段分组逻辑，这里不考虑容差，将容差字段分组单独做一个节点类型处理
	//这里的分组操作是在数据已经分组的基础上再次分组，分组数据不能跨原来的分组
	//遍历数据分组
	for _,item:= range (*req.Data) {
		result,errorCode:=nodeExecutor.group(models,&item)
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
	log.Println("nodeExecutorGroup run end")
	return flowResult,nil
}