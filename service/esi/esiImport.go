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

const (
	SHEETSELECTOR_OPTIONAL_YES="yes"
	SHEETSELECTOR_OPTIONAL_NO="no"
)

const (
	SHEETSELECTOR_TYPE_INDEX="index"
	SHEETSELECTOR_TYPE_NAME="name"
)

type DataRowHandler interface {
	handleRow(row map[string]interface{},sheetName string)(*common.CommonError)
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

type SheetSelector struct {
	Type string `json:"type"`
	Value string `json:"value"`
	Optional string `json:"optional"`
}

type ESIModel struct {
	ModelID string `json:"modelID"`
	Fields []ESIModelField `json:"fields"`
	Options ESIOption `json:"options"`
	SheetSelectors []SheetSelector `json:"sheets"`
	FileModel string `json:"fileModel"` 
	FileField string `json:"fileField"`
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

	err=parseBase64File(esi.FileName,esi.FileContent,contentHandler,dataRowHandler,esiModel.SheetSelectors)
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

