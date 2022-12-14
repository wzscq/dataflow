package esi

import (
	"dataflow/common"
	"log"
)

const (
	EXCEL_RANGE_TABLE="table"
	EXCEL_RANGE_RIGHTVAL="rightVal"
	EXCEL_RANGE_AUTO="auto"
)

const (
	DATA_SOURCE_INPUT="input"
)

const (
	END_ROW_YES="yes"
	END_ROW_NO="no"
	END_ROW_AUTO="auto"
)

const (
	EMPTY_VALUE_YES="yes"
	EMPTY_VALUE_NO="no"
)

const (
	PREVENT_SAME_FILE_YES="yes"
	PREVENT_SAME_FILE_NO="no"
)

type DataRowHandler interface {
	handleRow(row map[string]interface{})(*common.CommonError)
	onInit()(*common.CommonError)
	onOver(commit bool)(*common.CommonError)
}

type ESIModelField struct {
	Field string `json:"field"`
	LabelRegexp string `json:"labelRegexp"`
	ExcelRangeType string `json:"excelRangeType"`
	EndRow string `json:"endRow"`
	EmptyValue string `json:"emptyValue"`
	DetectedRangeType string `json:"detectedRangeType"`
	Source string `json:"source"`
}

type ESIOption struct {
	GenerateRowID bool `json:"generateRowID"` 
	MaxHeaderRow int `json:"maxHeaderRow"` 
	PreventSameFile bool `json:"preventSameFile"`
}

type ESIModel struct {
	ModelID string `json:"modelID"`
	Fields []ESIModelField `json:"fields"`
	Options ESIOption `json:"options"`
}

type ESImport struct {
	AppDB string
	ModelID string
	UserID string
	UserRoles string
	FileName string
	FileContent string
	InputRowData *map[string]interface{}
}

func (esi *ESImport)DoImport(esiModel *ESIModel)(map[string]interface{},*common.CommonError){
	log.Println("ESImport doImport start")

	dataRowHandler:=getDataRowHandler(
		esi.AppDB,
		esi.ModelID,
		esi.UserID,
		esi.UserRoles,
		esi.FileName,
		esiModel,
		esi.InputRowData)

	contentHandler:=getContentHandler(esiModel)

	err:=dataRowHandler.onInit()
	if err!=nil {
		return nil,err
	}

	err=parseBase64File(esi.FileName,esi.FileContent,contentHandler,dataRowHandler)
	if err!=nil {
		dataRowHandler.onOver(false)
		return nil,err
	}

	err=dataRowHandler.onOver(true)
	if err!=nil {
		return nil,err
	}

	result:=map[string]interface{}{
		"count":dataRowHandler.Count,
		"list":dataRowHandler.Result,
	}
	log.Println("esiImport doImport end")
	return result,nil
}

