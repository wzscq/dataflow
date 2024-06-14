package flow

import (
    "time"
	"dataflow/common"
	"dataflow/data"
	"log"
	"database/sql"
)

type nodeExecutorSave struct {
	NodeConf node
	DataRepository data.DataRepository
}

func (nodeExecutor *nodeExecutorSave)saveModels(
	dataItem *flowDataItem,
	tx *sql.Tx,
	instance *flowInstance)(int){

	for _,modelData:=range (dataItem.Models) {
		if modelData.List!=nil && len(*modelData.List)>0 {
			log.Printf("nodeExecutorSave save model:%s",*modelData.ModelID)
			saver:=data.Save{
				List:modelData.List,
				AppDB:instance.AppDB,
				UserID:instance.UserID,
				ModelID:*modelData.ModelID,
			}
			_,errorCode:=saver.SaveList(nodeExecutor.DataRepository,tx)
			if errorCode != common.ResultSuccess {
				log.Printf("nodeExecutorSave saveModels error:%d",errorCode)
				return errorCode
			}
		}
	}

	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorSave)save(
	dataItem *flowDataItem,
	instance *flowInstance)(int){
	
	log.Println("start nodeExecutorSave save")
	
	//开启事务
	tx,err:= nodeExecutor.DataRepository.Begin()
	if err != nil {
		log.Println(err)
		return common.ResultSQLError
	}

	//将分组号更新到左右表，同时更新左右表数据的匹配状态
	errorCode:=nodeExecutor.saveModels(dataItem,tx,instance)
	if errorCode!=common.ResultSuccess {
		tx.Rollback()
		log.Println("end nodeExecutorSave save with error")
		return errorCode
	}
	
	//提交事务
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return common.ResultSQLError
	}
	log.Println("end nodeExecutorSave save")
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorSave)run(
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
		"nodeType":NODE_SAVE,
	}

	//一个分组作为一个独立的事务进行保存
	for _,item:= range (*req.Data) {
		err:=nodeExecutor.save(&item,instance)
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