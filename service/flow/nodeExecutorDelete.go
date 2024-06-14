package flow

import (
    "time"
	"dataflow/common"
	"dataflow/data"
	"log"
	"database/sql"
)

type nodeExecutorDelete struct {
	NodeConf node
	DataRepository data.DataRepository
}

func (nodeExecutor *nodeExecutorDelete)deleteModels(
	dataItem *flowDataItem,
	tx *sql.Tx,
	instance *flowInstance,
	input *flowReqRsp)(int){

	for _,modelData:=range (dataItem.Models) {
		data.ProcessFilter(modelData.Filter,nil,input.GlobalFilterData,input.UserID,input.UserRoles,input.AppDB,nodeExecutor.DataRepository)
		delete:=data.BatchDelete{
			ModelID:*modelData.ModelID,
			SelectedRowKeys:modelData.SelectedRowKeys,
			AppDB:instance.AppDB,
			Filter:modelData.Filter,
			SelectAll:modelData.SelectAll,
		}
		_,errorCode:=delete.Delete(nodeExecutor.DataRepository,tx)
		if errorCode != common.ResultSuccess {
			log.Printf("nodeExecutorDelete deleteModels error:%d",errorCode)
			return errorCode
		}
	}

	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorDelete)delete(
	dataItem *flowDataItem,
	instance *flowInstance,
	input *flowReqRsp)(int){
	
	log.Println("start nodeExecutorDelete delete")
	
	//开启事务
	tx,err:= nodeExecutor.DataRepository.Begin()
	if err != nil {
		log.Println(err)
		return common.ResultSQLError
	}

	//将分组号更新到左右表，同时更新左右表数据的匹配状态
	errorCode:=nodeExecutor.deleteModels(dataItem,tx,instance,input)
	if errorCode!=common.ResultSuccess {
		tx.Rollback()
		log.Println("end nodeExecutorDelete delete with error")
		return errorCode
	}
	
	//提交事务
	if err := tx.Commit(); err != nil {
		log.Println(err)
		return common.ResultSQLError
	}
	log.Println("end nodeExecutorDelete delete")
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorDelete)run(
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
		"nodeType":NODE_DELETE,
	}

	//一个分组作为一个独立的事务进行保存
	for _,item:= range (*req.Data) {
		err:=nodeExecutor.delete(&item,instance,req)
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