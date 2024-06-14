package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"strconv"
)

type matchTolerance struct {
	Left int
	Right int
}

type tolerance struct {
	Left string `json:"left"`
	Right string `json:"right"`
}

type matchStep struct {
	MatchType string `json:"matchType"`
	Tolerance tolerance `json:"tolerance"`
}

type matchModel struct {
	ModelID string `json:"modelID"`
	Field string `json:"field"`
	Side string `json:"side"`
}

type matchConfig struct {
	Models []matchModel `json:"models"`
	Steps []matchStep `json:"steps"`
}

type matchGroup struct {
	LeftRows []int
	RightRows []int
}

type matchValue struct {
	LeftModels []modelDataItem
	RightModels []modelDataItem
	LeftValues []int
	RightValues []int
}

type nodeExecutorMatch struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorMatch)getMatchConfig()(*matchConfig){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok {
		return nil
	}

	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	matchConf:=&matchConfig{}
    if err := json.Unmarshal(jsonStr, matchConf); err != nil {
        log.Println(err)
		return nil
    }

	return matchConf
}

func (nodeExecutor *nodeExecutorMatch)getTolerance(t tolerance)(matchTolerance){
	
	leftVal,_:=strconv.ParseFloat(t.Left,64)
	rightVal,_:=strconv.ParseFloat(t.Right,64)
	
	//将容差转换为整数，这里因为主要考虑货币类型，因此容差只允许2为小数
	var f float64 = 100
	return matchTolerance{
		Left:int(leftVal*f),
		Right:int(rightVal*f),
	}
}

func (nodeExecutor *nodeExecutorMatch)getFieldValue(
	list *[]map[string]interface{},
	field string)(*[]int,int){
	
	valArray:= make([]int,len(*list))
	for index,row:=range (*list){
		fieldVal,found:=row[field]
		if !found {
			log.Printf("nodeExecutorMatch getFieldValue no field: %s!\n", field)
			return nil,common.ResultNoMatchField 
		}

		switch fieldVal.(type) {
		case float64:
			fVal, _ := fieldVal.(float64)
			valArray[index]=int(fVal*100)
		case string:
			sVal, _ := fieldVal.(string)
			fVal, err := strconv.ParseFloat(sVal, 64)
			if err !=nil {
				log.Printf("nodeExecutorMatch getFieldValue can not convert value to float64: %s!\n", sVal)
				return nil,common.ResultMatchValueToFloat64Error
			}
			valArray[index]=int(fVal*100)
		default:
			log.Printf("nodeExecutorMatch getFieldValue not supported field type: %T!\n", fieldVal)
			return nil,common.ResultMatchFieldTypeError
		}
	}
	
	return &valArray,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorMatch)getModelRow(
	groupModels *[]modelDataItem,
	rowIdx int)(string,map[string]interface{}){

	for _,modelItem:=range(*groupModels){
		rowCount:=len(*modelItem.List)
		if rowIdx>=rowCount {
			rowIdx=rowIdx-rowCount
		} else {
			return *modelItem.ModelID,(*modelItem.List)[rowIdx]
		}
	}

	return "",nil
}

func (nodeExecutor *nodeExecutorMatch)getGroupModel(
	groupedModels *[]modelDataItem,
	modelID string)(*modelDataItem){
	
	for index,_:=range(*groupedModels){
		if *((*groupedModels)[index].ModelID) == modelID {
			return &((*groupedModels)[index])
		}
	}
	return nil
}

func (nodeExecutor *nodeExecutorMatch)createGroupModel(modelID string)(*modelDataItem){
	return &modelDataItem{
		ModelID:&modelID,
		List:&[]map[string]interface{}{},
	}
}

func (nodeExecutor *nodeExecutorMatch)updateDataItem(
	rows *[]int,
	groupModels *[]modelDataItem,
	groupedModels *[]modelDataItem){
	//把记录从原来的list移动到新的list中
	for _,row:=range (*rows) {
		modelID,rowData:=nodeExecutor.getModelRow(groupModels,row)
		if rowData !=nil {
			groupModel:=nodeExecutor.getGroupModel(groupedModels,modelID)
			if groupModel==nil {
				groupModel=nodeExecutor.createGroupModel(modelID)
				(*groupedModels)=append((*groupedModels),*groupModel)
			}
			*(groupModel.List)=append(*(groupModel.List),rowData)
		}
	}
}

func (nodeExecutor *nodeExecutorMatch)updateResult(
	matchConf *matchConfig,
	matchedGroup *matchGroup,
	matchValue *matchValue,
	result *[]flowDataItem){
	
	//这里左表和右表的数据可能来自不同的表
	groupModels:=[]modelDataItem{}
	if matchedGroup.LeftRows!=nil && len(matchedGroup.LeftRows) >0 {
		nodeExecutor.updateDataItem(&matchedGroup.LeftRows,&(matchValue.LeftModels),&groupModels)
	}

	if matchedGroup.RightRows!=nil && len(matchedGroup.RightRows) >0  {
		nodeExecutor.updateDataItem(&matchedGroup.RightRows,&(matchValue.RightModels),&groupModels)
	}

	newGroup:=flowDataItem{
		Models:groupModels,
	}
	(*result)=append((*result),newGroup)
}

func (nodeExecutor *nodeExecutorMatch)removeGroupedIndex(
	groupedIndex,allIndex *[]int){
	if len(*allIndex) <= len(*groupedIndex) {
		*allIndex=[]int{}
		return
	}

	remainIndex:=make([]int,len(*allIndex)-len(*groupedIndex))
	index:=0
	grouped:=false
	for _,rowIndex:=range (*allIndex) {
		grouped=false
		for _,groupidx:=range (*groupedIndex) {
			//log.Printf("removeGroupedIndex rowIndex: %d,groupidx:%d",rowIndex,groupidx)
			if rowIndex==groupidx {
				grouped=true
			}
		}
		if !grouped {
			remainIndex[index]=rowIndex
			index+=1
		} 
	}
	*allIndex=remainIndex
}

func (nodeExecutor *nodeExecutorMatch)updateForMatchGroup(matchedGroup,forMatch *matchGroup){
	nodeExecutor.removeGroupedIndex(&matchedGroup.LeftRows,&forMatch.LeftRows)
	nodeExecutor.removeGroupedIndex(&matchedGroup.RightRows,&forMatch.RightRows)
}

func (nodeExecutor *nodeExecutorMatch)matchStep(
	matchConf *matchConfig,
	step *matchStep,
	matchValue *matchValue,
	forMatch *matchGroup,
	result *[]flowDataItem)(int){
	
	matchExecutor:=getMatchExecutor(step.MatchType)
	if matchExecutor==nil {
		return common.ResultNotSupportedMatchType
	}
	//数值计算均使用整数，这里将数字转换为整数参与后续计算
	tolerance:=nodeExecutor.getTolerance(step.Tolerance)
	//循环获取分组
	for {
		if len(forMatch.LeftRows)==0 || len(forMatch.RightRows)==0 {
			break
		}
		matchedGroup:=matchExecutor.getMatchGroup(matchValue,forMatch,tolerance)
		if matchedGroup != nil {
			nodeExecutor.updateResult(matchConf,matchedGroup,matchValue,result)
			nodeExecutor.updateForMatchGroup(matchedGroup,forMatch)
		} else {
			break
		}
	}
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorMatch)getMatchValue(
	matchConf *matchConfig,
	item *flowDataItem)(*matchValue,int){

	matchValue:=matchValue{
		LeftValues:[]int{},
		RightValues:[]int{},
		LeftModels:[]modelDataItem{},
		RightModels:[]modelDataItem{},
	}
	for _,modelData:=range (item.Models){
		//log.Println("getMatchValue 1:" ,*modelData.ModelID,matchConf.Left.ModelID)
		for _,matchModle:=range (matchConf.Models) {
			if *modelData.ModelID == matchModle.ModelID && modelData.List != nil {
				matchValues,errorCode:=nodeExecutor.getFieldValue(modelData.List,matchModle.Field)
				if errorCode!=common.ResultSuccess{
					return nil,errorCode
				}

				if matchModle.Side == MODEL_SIDE_LEFT {
					matchValue.LeftValues=append(matchValue.LeftValues,(*matchValues)...)
					matchValue.LeftModels=append(matchValue.LeftModels,modelData)
				}

				if matchModle.Side == MODEL_SIDE_RIGHT {
					matchValue.RightValues=append(matchValue.RightValues,(*matchValues)...)
					matchValue.RightModels=append(matchValue.RightModels,modelData)
				}
			}
		}
	}
	
	return &matchValue,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorMatch)getInitMatchGroup(matchValue *matchValue)(*matchGroup){
	//用于记录左右表待比对数据的索引数组，初始化时所有数据参与比对，每次比对后将比对成功的索引删除
	forMatch:=matchGroup{
		LeftRows:nil,
		RightRows:nil, 
	}
	
	if matchValue.LeftValues!=nil {
		forMatch.LeftRows=make([]int,len(matchValue.LeftValues))
		for index,_:= range (matchValue.LeftValues) {
			forMatch.LeftRows[index]=index
		}
	}

	if matchValue.RightValues!=nil {
		forMatch.RightRows=make([]int,len(matchValue.RightValues))
		for index,_:= range (matchValue.RightValues) {
			forMatch.RightRows[index]=index
		}
	}
	return &forMatch
}

func (nodeExecutor *nodeExecutorMatch)match(
	matchConf *matchConfig,
	item *flowDataItem,
	result *[]flowDataItem)(int){

	//获取待比对数据	
	matchValue,errorCode:=nodeExecutor.getMatchValue(matchConf,item)
	if errorCode!=common.ResultSuccess {
		return errorCode
	}

	//获取待比对数据索引
	forMatch:=nodeExecutor.getInitMatchGroup(matchValue)

	//没有需要匹配的数据
	if len(forMatch.LeftRows)==0 || len(forMatch.RightRows)==0 {
		nodeExecutor.updateResult(matchConf,forMatch,matchValue,result)
		return common.ResultSuccess
	}
	
	//迭代处理每个step
	for _,step:=range (matchConf.Steps) {
		errorCode:=nodeExecutor.matchStep(matchConf,&step,matchValue,forMatch,result)
		if errorCode!=common.ResultSuccess {
			return errorCode
		}
		//没有需要匹配的数据,如果有剩余的数据，将剩余数据放到单独的分组
		if len(forMatch.LeftRows)==0 || len(forMatch.RightRows)==0 {
			break
		}
	}

	//所有比对步骤完成，如果仍有剩余的左右表数据，则将数据放入单独的分组中
	if len(forMatch.LeftRows)>0 {
		matchGroup:=matchGroup{
			LeftRows:forMatch.LeftRows,
			RightRows:[]int{},
		}
		nodeExecutor.updateResult(matchConf,&matchGroup,matchValue,result)
	}

	if len(forMatch.RightRows)>0 {
		matchGroup:=matchGroup{
			LeftRows:[]int{},
			RightRows:forMatch.RightRows,
		}
		nodeExecutor.updateResult(matchConf,&matchGroup,matchValue,result)
	}

	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorMatch)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

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
		Fields:req.Fields,
		Pagination:req.Pagination,
		Operation:req.Operation,
		SelectAll:req.SelectAll,
		GoOn:true,
	}

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_MATCH,
	}

	//加载节点配置
	matchConf:=nodeExecutor.getMatchConfig()
	if matchConf==nil {
		log.Printf("nodeExecutorMatch run get node config error\n")
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if node.Input.Data==nil || len(*node.Input.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=node.Input
		return node.Input,nil
	}

	resultData:=&[]flowDataItem{}
	//这里的分组操作是在数据已经分组的基础上再次分组，分组数据不能跨原来的分组
	//遍历每个数据分组
	for _,item:= range (*req.Data) {
		errorCode:=nodeExecutor.match(matchConf,&item,resultData)
		if errorCode!=common.ResultSuccess {
			return flowResult,common.CreateError(errorCode,params)
		}
	}

	flowResult.Data=resultData
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	return flowResult,nil
}