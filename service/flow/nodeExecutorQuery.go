package flow

import (
  "time"
	"dataflow/common"
	"dataflow/data"
	"encoding/json"
	"log"
)

type queryModel struct {
	ModelID string `json:"modelID"`
	Fields *[]data.Field  `json:"fields"`
	Filter *map[string]interface{} `json:"filter"`
	Sorter *[]data.Sorter  `json:"sorter"`
}

type nodeExecutorQuery struct {
	DataRepository data.DataRepository
	NodeConf node
}

func (nodeExecutor *nodeExecutorQuery)getQueryModels()([]queryModel){
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

func (nodeExecutor *nodeExecutorQuery)getQeruyFields(confFields *[]data.Field)([]data.Field){
	fields:=append(*confFields,
		data.Field{
			Field:"id",
		},
		data.Field{
			Field:"version",
		})
	return fields
}

func (nodeExecutor *nodeExecutorQuery)queryModel(appDB string,modelQueryConf queryModel)(*modelDataItem,int){
	//过滤字段中默认增加id、version字段
	fields:=nodeExecutor.getQeruyFields(modelQueryConf.Fields)
	query:=&data.Query{
		ModelID:modelQueryConf.ModelID,
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
		ModelID:&modelQueryConf.ModelID,
		List:&result.List,
		Total:result.Total,
	}
	return modelData,errorCode
}

func (nodeExecutor *nodeExecutorQuery)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
	log.Println("nodeExecutorQuery run start")

	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID, 
		UserID:req.UserID,
		AppDB:req.AppDB,
		GoOn:true,
	}

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_QUERY,
	}

	//根据配置，从数据表中读取数据，将读取的数据传入下一步操作
	//如何避免数据量过大，目前临时通过取数配置中增加限制条件来控制数据量，
	//可以基于过滤条件分多次对不同的数据进行处理
	models:=nodeExecutor.getQueryModels()
	if models==nil {
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	modelDatas:=[]modelDataItem{}
	for _, model := range (models) {
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
	log.Println("nodeExecutorQuery run end")
	return flowResult,nil
}