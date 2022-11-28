package flow

import (
    "time"
	"buoyancyinfo.com/dataflow/common"
	"buoyancyinfo.com/dataflow/data"
	"encoding/json"
	"log"
)

type nodeExecutorRequestQuery struct {
	DataRepository data.DataRepository
	NodeConf node
}

func (nodeExecutor *nodeExecutorRequestQuery)getQeruyFields(
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

func (nodeExecutor *nodeExecutorRequestQuery)getQueryModels()([]queryModel){
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
	models:=[]queryModel{}
    if err := json.Unmarshal(jsonStr, &models); err != nil {
        log.Println(err)
		return nil
    }

	return models
}

func (nodeExecutor *nodeExecutorRequestQuery)queryModel(
	appDB string,modelQueryConf queryModel)(*modelDataItem,int){
	//过滤字段中默认增加id、version字段
	fields:=nodeExecutor.getQeruyFields(modelQueryConf.Fields)
	query:=&data.Query{
		ModelID:modelQueryConf.ModelID,
		Filter:modelQueryConf.Filter,
		Sorter:modelQueryConf.Sorter,
		Fields:&fields,
		AppDB:appDB,
		Pagination:&data.Pagination{
			Current:1,
			PageSize:100000,
		},
	}
	result,errorCode:=query.Execute(nodeExecutor.DataRepository,false)
	modelData:=&modelDataItem{
		ModelID:&modelQueryConf.ModelID,
		List:&result.List,
		Total:result.Total,
	}
	return modelData,errorCode
}

func (nodeExecutor *nodeExecutorRequestQuery)updateFilterBySelectedRowKeys(
	confFilter *map[string]interface{},
	selectedRowKeys *[]string)(*map[string]interface{},int){
		
	filter:=&map[string]interface{}{
			"id":map[string]interface{}{
				data.Op_in:*selectedRowKeys,
			},
		}
			
	if 	confFilter!=nil && len(*confFilter) > 0 {
		filter=&map[string]interface{}{
			data.Op_and:[]interface{}{
				*confFilter,
				*filter,
			},
		}
	}
			
	return filter,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorRequestQuery)updateFilterByQueryFilter(
	confFilter *map[string]interface{},
	queryFilter *map[string]interface{},)(*map[string]interface{},int){
	
	if confFilter == nil {
		return queryFilter,common.ResultSuccess
	}

	filter:=&map[string]interface{}{
		data.Op_and:[]interface{}{
			*confFilter,
			*queryFilter,
		},
	}
	return filter,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorRequestQuery)updateModels(
	input *flowReqRsp,
	models *[]queryModel)(int){
	
	for index,model := range (*models) {
		log.Printf("nodeExecutorRequestQuery updateModels %s,%s\n",model.ModelID,*input.ModelID)
		if model.ModelID==*input.ModelID {
			//这里勾选数据优先，如果前端传入了勾选数据，则使用and合并前端传入的勾选数据和配置的filter
			errorCode:=common.ResultSuccess
			if input.SelectedRowKeys!=nil && len(*input.SelectedRowKeys)>0 {
				log.Printf("nodeExecutorRequestQuery updateModels updateFilterBySelectedRowKeys \n")
				(*models)[index].Filter,errorCode=nodeExecutor.updateFilterBySelectedRowKeys(model.Filter,input.SelectedRowKeys)
				return errorCode
			}
			//如果前端传入了Filter，则使用and合并前端传入的Filter和配置的filter
			if input.Filter!=nil && len(*input.Filter)>0 {
				log.Printf("nodeExecutorRequestQuery updateModels updateFilterByQueryFilter \n")
				(*models)[index].Filter,errorCode=nodeExecutor.updateFilterByQueryFilter(model.Filter,input.Filter)
				return errorCode
			}
		}
	}
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorRequestQuery)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
	log.Println("nodeExecutorRequestQuery run start")

	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID, 
		UserID:req.UserID,
		AppDB:req.AppDB,
	}

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_REQUEST_QUERY,
	}

	//根据配置，从数据表中读取数据，将读取的数据传入下一步操作
	//如何避免数据量过大，目前临时通过取数配置中增加限制条件来控制数据量，
	//可以基于过滤条件分多次对不同的数据进行处理
	models:=nodeExecutor.getQueryModels()
	if models==nil {
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	//update request filter
	data.ProcessFilter(req.Filter,req.FilterData,req.UserID,req.UserRoles,req.AppDB,nodeExecutor.DataRepository)

	//根据页面请求传入的查询条件更新模型的查询条件
	if req.ModelID!=nil && (req.Filter!=nil || req.SelectedRowKeys!=nil) {
		log.Println("nodeExecutorRequestQuery updateModels")
		errCode:=nodeExecutor.updateModels(req,&models)
		if errCode!=common.ResultSuccess {
			return flowResult,common.CreateError(errCode,params)
		}
	}

	modelDatas:=[]modelDataItem{}
	for _, model := range (models) {
		//postJson,_:=json.Marshal(model)
		//log.Println(string(postJson))
		modelData,errCode:=nodeExecutor.queryModel(instance.AppDB,model)
		if errCode!=common.ResultSuccess {
			return flowResult,common.CreateError(errCode,params)
		}
		//结果放入data中
		modelDatas=append(modelDatas,*modelData)
	}

	//如果之前的查询中已经存在数据项，则将当前查询的项目和并到之前查询数据的第0个项目上
	data:=req.Data
	if data ==nil || len(*data)==0 {
		data=&[]flowDataItem{
			flowDataItem{
				Models:modelDatas,
			},
		}
	} else {
		(*data)[0].Models=append((*data)[0].Models,modelDatas...)
	}

	flowResult.Data=data
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorRequestQuery run end")
	return flowResult,nil
}