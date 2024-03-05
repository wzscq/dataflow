package flow

import (
    "time"
	"dataflow/common"
	"dataflow/data"
	"encoding/json"
	"log"
)

type relatedFiled struct {
	Field string `json:"field"`
	RelatedField string `json:"relatedField"`
}

type relatedQueryModel struct {
	ModelID string `json:"modelID"`
	Fields *[]data.Field  `json:"fields"`
	Filter *map[string]interface{} `json:"filter"`
	RelatedModel string `json:"relatedModel"`
	RelatedFileds []relatedFiled `json:"relatedFields"`
	Sorter *[]data.Sorter  `json:"sorter"`
}

type nodeExecutorRelatedQuery struct {
	DataRepository data.DataRepository
	NodeConf node
}

func (nodeExecutor *nodeExecutorRelatedQuery)getQeruyFields(
	confFields *[]data.Field)([]data.Field){
	fields:=append(*confFields,
		data.Field{
			Field:"id",
		},
		data.Field{
			Field:"version",
		})
	return fields
}

func (nodeExecutor *nodeExecutorRelatedQuery)getQueryModels()([]relatedQueryModel){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok || mapData["models"]==nil {
		return nil
	}
	//models,ok=mapData["models"].([]interface{})
	//if !ok {
	//	return nil
	//}
	jsonStr, err := json.Marshal(mapData["models"])
    if err != nil {
        log.Println(err)
		return nil
    }
	models:=[]relatedQueryModel{}
    if err := json.Unmarshal(jsonStr, &models); err != nil {
        log.Println(err)
		return nil
    }

	return models
}

func (nodeExecutor *nodeExecutorRelatedQuery)updateModelsFilter(
	input *flowReqRsp,
	model *relatedQueryModel)(int){
	data.ProcessFilter(model.Filter,nil,nil,input.UserID,input.UserRoles,input.AppDB,nodeExecutor.DataRepository)
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorRelatedQuery)queryModel(
	appDB string,modelQueryConf *relatedQueryModel)(*modelDataItem,int){
	//过滤字段中默认增加id、version字段
	modelID:=modelQueryConf.ModelID
	fields:=nodeExecutor.getQeruyFields(modelQueryConf.Fields)
	query:=&data.Query{
		ModelID:modelID,
		Filter:modelQueryConf.Filter,
		Fields:&fields,
		AppDB:appDB,
		Sorter:modelQueryConf.Sorter,
		Pagination:&data.Pagination{
			Current:1,
			PageSize:100000,
		},
	}
	result,errorCode:=query.Execute(nodeExecutor.DataRepository,false)
	
	modelData:=&modelDataItem{
		ModelID:&modelID,
		List:&result.List,
		Total:result.Total,
	}
	return modelData,errorCode
}

func (nodeExecutor *nodeExecutorRelatedQuery)mergeFilter(
	modelFilter *map[string]interface{},
	realtedFilter *map[string]interface{},)(*map[string]interface{}){
	
	if modelFilter == nil {
		return realtedFilter
	}

	filter:=&map[string]interface{}{
		data.Op_and:[]interface{}{
			*modelFilter,
			*realtedFilter,
		},
	}
	return filter
}

func (nodeExecutor *nodeExecutorRelatedQuery)getRelatedModelData(
	model *relatedQueryModel,
	dateItem *flowDataItem)(*[]map[string]interface{}){
	//如果关联模型的ID为空的话就直接返回空
	if len(model.RelatedModel)<=0 {
		return nil
	}
	
	//查找对应模型的数据
	for _,modelData := range(dateItem.Models) {
		if *modelData.ModelID == model.RelatedModel {
			return modelData.List
		}
	}
	return nil
}

func (nodeExecutor *nodeExecutorRelatedQuery)getStringValue(val interface{})(string){
	switch val.(type) {
	case string:
		sVal, _ := val.(string)
		return sVal
	default:
		log.Printf("nodeExecutorRelatedQuery getStringValue not supported value type %T!\n", val)
		return ""
	}
}

func (nodeExecutor *nodeExecutorRelatedQuery)getRelatedRows(
	rows *[]map[string]interface{},
	relatedFileds []relatedFiled)(*map[string]interface{},int){
	//循环获取每个关联行的数据,同时去重
	rfRows:=map[string]interface{}{}
	for _,row := range(*rows) {
		rfRowKey:=""
		rfRow:=map[string]interface{}{}
		for _,field := range(relatedFileds) {
			iVal,ok:=row[field.RelatedField]
			if !ok {
				return nil,common.ResultRelatedQueryNoRelatedField
			}
			sVal:=nodeExecutor.getStringValue(iVal)
			rfRowKey=rfRowKey+"#"+sVal
			rfRow[field.Field]=sVal
		}
		rfRows[rfRowKey]=rfRow
	}
	return &rfRows,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorRelatedQuery)convertRelatedRowsToFilter(
	rfRows *map[string]interface{})(*map[string]interface{}){
	//行间使用or
	orFilters:=[]interface{}{}
	//行内采用and
	for _,row := range(*rfRows) {
		orFilters=append(orFilters,row)
	}
	filter:=map[string]interface{}{
		data.Op_or:orFilters,
	}

	return &filter
}

func (nodeExecutor *nodeExecutorRelatedQuery)getRelatedFilter(
	model *relatedQueryModel,
	dataItem *flowDataItem)(*map[string]interface{},int){

	if len(model.RelatedModel)<=0 || model.RelatedFileds == nil || len(model.RelatedFileds)<=0 {
		return nil,common.ResultNodeConfigError
	}

	modelData:=nodeExecutor.getRelatedModelData(model,dataItem)
	if modelData==nil || len(*modelData)==0 {
		return nil,common.ResultSuccess
	}

	rfRows,errorCode:=nodeExecutor.getRelatedRows(modelData,model.RelatedFileds)
	if errorCode != common.ResultSuccess {
		return nil,errorCode
	}

	return nodeExecutor.convertRelatedRowsToFilter(rfRows),common.ResultSuccess
}

func (nodeExecutor *nodeExecutorRelatedQuery)dealDataItem(
	dataItem *flowDataItem,
	model *relatedQueryModel,
	appDB string,
	input *flowReqRsp)(int){

	relatedFilter,errorCode:=nodeExecutor.getRelatedFilter(model,dataItem)
	if errorCode!=common.ResultSuccess {
		return errorCode
	}
	//当关联查询不为空的时候才执行查询动作，
	//所以如果是用了关联查询的组件则但是又没有配置关联条件的情况下是不执行数据查询的
	if relatedFilter!=nil {
		model.Filter=nodeExecutor.mergeFilter(model.Filter,relatedFilter)
		nodeExecutor.updateModelsFilter(input,model)
		postJson,_:=json.Marshal(model)
		log.Println(string(postJson))
		//return flowResult,nil
		modelData,errCode:=nodeExecutor.queryModel(appDB,model)
		if errCode!=common.ResultSuccess {
			return errCode
		}
		//结果放入data中
		dataItem.Models=append(dataItem.Models,*modelData)
	}
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorRelatedQuery)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
	log.Println("nodeExecutorRelatedQuery run start")

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

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_RELATED_QUERY,
	}

	//根据配置，从数据表中读取数据，将读取的数据传入下一步操作
	//如何避免数据量过大，目前临时通过取数配置中增加限制条件来控制数据量，
	//可以基于过滤条件分多次对不同的数据进行处理
	models:=nodeExecutor.getQueryModels()
	if models==nil {
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if req.Data!=nil && len(*req.Data)>0 {
		for index,_:=range(*req.Data){
			for _, model := range (models) {
				errorCode:=nodeExecutor.dealDataItem(&((*req.Data)[index]),&model,instance.AppDB,req)
				if errorCode!=common.ResultSuccess {
					return flowResult,common.CreateError(errorCode,params)
				}
			}
		}
	}

	flowResult.Data=req.Data
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorRelatedQuery run end")
	return flowResult,nil
}