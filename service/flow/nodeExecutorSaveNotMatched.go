package flow

import (
    "time"
	"dataflow/common"
	"dataflow/data"
	"log"
	"database/sql"
)

type nodeExecutorSaveNotMatched struct {
	NodeConf node
	DataRepository data.DataRepository
}

func (nodeExecutor *nodeExecutorSaveNotMatched)saveModels(
	dataItem *flowDataItem,
	tx *sql.Tx,
	message string,
	instance *flowInstance)(int){

	for _,modelData:=range (dataItem.Models) {
		if modelData.List!=nil && len(*modelData.List)>0 {
			rows:=make([]map[string]interface{},len(*modelData.List))
			for index,row:=range (*modelData.List) {
				rows[index]=map[string]interface{}{}
				rows[index][data.CC_ID]=row[data.CC_ID]
				rows[index][data.CC_VERSION]=row[data.CC_VERSION]
				rows[index][data.SAVE_TYPE_COLUMN]=data.SAVE_UPDATE
				rows[index][CC_MATCH_STATUS]=MATCH_STATUS_FAILURE
				rows[index][CC_MATCH_MESSAGE]=message
			}

			saver:=data.Save{
				List:&rows,
				AppDB:instance.AppDB,
				UserID:instance.UserID,
				ModelID:*modelData.ModelID,
			}
			_,errorCode:=saver.SaveList(nodeExecutor.DataRepository,tx)
			if errorCode != common.ResultSuccess {
				return errorCode
			}
		}
	}

	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorSaveNotMatched)getMessage(dataItem *flowDataItem)(string){
	var message string
	for _,verifyItem:=range dataItem.VerifyResult {
		message=message+verifyItem.Message+"; "
	}
	return message
}

func (nodeExecutor *nodeExecutorSaveNotMatched)save(
	dataItem *flowDataItem,
	instance *flowInstance){
	
	log.Println("start nodeExecutorSaveNotMatched save")
	
	//开启事务
	tx,err:= nodeExecutor.DataRepository.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	message:=nodeExecutor.getMessage(dataItem)
	//将分组号更新到左右表，同时更新左右表数据的匹配状态
	errorCode:=nodeExecutor.saveModels(dataItem,tx,message,instance)
	if errorCode!=common.ResultSuccess {
		tx.Rollback()
		log.Println("end nodeExecutorSaveNotMatched save")
		return
	}
	
	//提交事务
	if err := tx.Commit(); err != nil {
		log.Println(err)
	}
	log.Println("end nodeExecutorSaveNotMatched save")
}

func (nodeExecutor *nodeExecutorSaveNotMatched)run(
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
		UserID:req.UserID,
		AppDB:req.AppDB,
	}

	//一个分组作为一个独立的事务进行保存
	for _,item:= range (*req.Data) {
		nodeExecutor.save(&item,instance)
	}
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	return flowResult,nil
}