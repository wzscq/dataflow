package esi

import (
	"dataflow/common"
	"dataflow/data"
	"log"
	//"github.com/rs/xid"
	"fmt"
)

const (
	CC_BATCH_NUMBER="import_batch_number"
	CC_IMPORT_FILE="import_file_name"
	CC_SHEET_NAME="sheet_name"
)

type listDataRowHandler struct {
	AppDB string
	ModelID string
	UserID string
	UserRoles string
	FileName string
	ESIModel *ESIModel
	Count int
	InputRowData *map[string]interface{}
	Result []map[string]interface{}
	BatchID string
}

func getDataRowHandler(
	appDB,modelID,userID,userRoles,fileName string,
	esiModel *ESIModel,
	inputRowData *map[string]interface{} )(*listDataRowHandler){
	return &listDataRowHandler{
		AppDB:appDB,
		ModelID:modelID,
		UserID:userID,
		FileName:fileName,
		UserRoles:userRoles,
		Count:0,
		ESIModel:esiModel,
		Result:[]map[string]interface{}{},
		InputRowData:inputRowData,
	}
}

func (dataHandler *listDataRowHandler)onInit()(*common.CommonError){
	dataHandler.BatchID=dataHandler.getBatchID()
	return nil
}

func (dataHandler *listDataRowHandler)onOver(commit bool)(*common.CommonError){
	return nil
}

func (dataHandler *listDataRowHandler)getRowID(batchID string)(string){
	return fmt.Sprintf("%s%05d",batchID,dataHandler.Count)
}

func (dataHandler *listDataRowHandler)getBatchID()(string){
	guid := GetBatchID()  //xid.New().String()
	return guid
}

func (dataHandler *listDataRowHandler)updateRowData(row *map[string]interface{},sheetName string){
	(*row)[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE
	//从表单输入的数据写入导入记录对应字段上
	for fidx,_:=range(dataHandler.ESIModel.Fields) {
		esiField:=&dataHandler.ESIModel.Fields[fidx]
		if esiField.Source==DATA_SOURCE_INPUT {
			(*row)[esiField.Field]=(*dataHandler.InputRowData)[esiField.Field]
		}
	}
	//生成导入文件名+ID+批次号
	batchID:=dataHandler.BatchID
	(*row)[CC_BATCH_NUMBER]=batchID
	(*row)[CC_IMPORT_FILE]=dataHandler.FileName
	(*row)[CC_SHEET_NAME]=sheetName
	if dataHandler.ESIModel.Options.GenerateRowID == true {
		(*row)[data.CC_ID]=dataHandler.getRowID(batchID)
	}
}

func (dataHandler *listDataRowHandler)handleRow(row map[string]interface{},sheetName string)(*common.CommonError){
	log.Println("listDataRowHandler handleRow start")
	dataHandler.Count++
	dataHandler.updateRowData(&row,sheetName)
	dataHandler.Result=append(dataHandler.Result,row)	
	log.Println("listDataRowHandler handleRow end")
	return nil
}



