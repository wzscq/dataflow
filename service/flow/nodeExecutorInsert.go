package flow

import (
    "time"
	"dataflow/common"
	"dataflow/data"
	"log"
	"database/sql"
	"encoding/json"
	"strconv"
)

type InsertParams struct {
	SQLMaxLen string `json:"sqlMaxLen"`
}

type nodeExecutorInsert struct {
	NodeConf node
	DataRepository data.DataRepository
}

func (nodeExecutor *nodeExecutorInsert)getNodeConfig()(*InsertParams){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok {
		return nil
	}

	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	insertParams:=&InsertParams{}
    if err := json.Unmarshal(jsonStr, insertParams); err != nil {
        log.Println(err)
		return nil
    }

	return insertParams
}

func (nodeExecutor *nodeExecutorInsert)insertModels(
	dataItem *flowDataItem,
	tx *sql.Tx,
	instance *flowInstance,
	insertConf *InsertParams)(int){

	for _,modelData:=range (dataItem.Models) {
		if modelData.List!=nil && len(*modelData.List)>0 {
			log.Printf("nodeExecutorInsert insert model:%s",*modelData.ModelID)
			var sqlMaxLen int64 = 32000
			if(insertConf.SQLMaxLen!=""){
				sqlMaxLen,_=strconv.ParseInt(insertConf.SQLMaxLen,0,32)
			}
			insert:=data.BatchInsert{
				List:modelData.List,
				AppDB:instance.AppDB,
				UserID:instance.UserID,
				ModelID:*modelData.ModelID,
				SQLMaxLen:int(sqlMaxLen),
			}
			errorCode:=insert.Insert(nodeExecutor.DataRepository,tx)
			if errorCode != common.ResultSuccess {
				log.Printf("nodeExecutorInsert insertModels error:%d",errorCode)
				return errorCode
			}
		}
	}

	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorInsert)insert(
	dataItem *flowDataItem,
	instance *flowInstance,
	insertConf *InsertParams)(int){
	
	log.Println("start nodeExecutorInsert insert")
	
	//开启事务
	tx,err:= nodeExecutor.DataRepository.Begin()
	if err != nil {
		log.Println(err)
		return common.ResultSQLError
	}

	//将分组号更新到左右表，同时更新左右表数据的匹配状态
	errorCode:=nodeExecutor.insertModels(dataItem,tx,instance,insertConf)
	if errorCode!=common.ResultSuccess {
		tx.Rollback()
		log.Println("end nodeExecutorInsert insert with error")
		return errorCode
	}
	
	//提交事务
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return common.ResultSQLError
	}
	log.Println("end nodeExecutorInsert insert")
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorInsert)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

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
		Fields:req.Fields,
		Pagination:req.Pagination,
		Operation:req.Operation,
		SelectAll:req.SelectAll,
		GoOn:true,
	}

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_INSERT,
	}

	conf:=nodeExecutor.getNodeConfig()
	if conf==nil {
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	//一个分组作为一个独立的事务进行保存
	for _,item:= range (*req.Data) {
		err:=nodeExecutor.insert(&item,instance,conf)
		if err!=common.ResultSuccess {
			return flowResult,common.CreateError(err,params)
		}
	}
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	return flowResult,nil
}