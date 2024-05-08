package data

import (
	"dataflow/common"
	"database/sql"
	"log"
	"strings"
	"fmt"
)

type BatchInsert struct {
	ModelID         string    `json:"modelID"`
	AppDB           string    `json:"appDB"`
	UserID string `json:"userID"`
	List *[]map[string]interface{} `json:"list"` 
	SQLMaxLen int `json:"sqlMaxLen"`
}

func (insert *BatchInsert) getFields(data map[string]interface{}) ([]string) {
	fields:=[]string{}
	for k,_:=range data {
		if k==CC_VERSION || k==CC_CREATE_TIME || k==CC_CREATE_USER || k==CC_UPDATE_TIME || k==CC_UPDATE_USER || k==SAVE_TYPE_COLUMN {
			continue
		}
		fields=append(fields,k)
	}
	return fields
}

func (insert *BatchInsert) getRowValues(fields *[]string, row map[string]interface{}) (*string){
	value:=""
	for _,field:=range *fields {
		if value!="" {
			value+=","
		}

		v:=row[field]
		switch vt := v.(type) {
		case string:
			sVal, _ := v.(string)
			value+=("'"+sVal+"'")
		case int64:
			iVal, _ := v.(int64)
			sVal:=fmt.Sprintf("%d",iVal)
			value+=sVal
		case float64:
			fVal, _ := v.(float64)
			sVal:=fmt.Sprintf("%f",fVal)
			value+=sVal
		case nil:
			value+="null"
		default:
			log.Printf("getRowValues not supported value type %T!\n", vt)
			return nil
		}
	}
	return &value
}

func (insert *BatchInsert) getValues(
	fields *[]string, 
	list *[]map[string]interface{},
	commonFieldsValue string,
	startRow,valueLength int) (*[]string,int) {
	values:=[]string{}
	for startRow<len(*list)&&valueLength>0 {
		data:=(*list)[startRow]
		value:=insert.getRowValues(fields,data)
		if value==nil {
			return nil,startRow
		}
		rowStr:="("+*value+","+commonFieldsValue+")"
		valueLength=valueLength-len(rowStr)-1
		if valueLength>=0 {
			startRow++
			values=append(values,"("+*value+","+commonFieldsValue+")")
		}
	}
	return &values,startRow
}

func (insert *BatchInsert) Insert(dataRepository DataRepository, tx *sql.Tx) (int) {
	//获取所有待删数据查询条件
	if insert.List == nil || len(*insert.List) == 0 {
		return common.ResultSuccess
	}

	commonFields,commonFieldsValue:=GetCreateCommonFieldsValues(insert.UserID)
	fields:=insert.getFields((*insert.List)[0])
	fieldsStr:=strings.Join(fields,",")
	sqlPart1 := "insert into " + insert.AppDB + "." + insert.ModelID + "("+fieldsStr+","+commonFields+") values "
	valueLength:=insert.SQLMaxLen-len(sqlPart1)
	if valueLength<0 {
		log.Printf("Insert error: values length is negative.\n")
		return common.ResultSQLError
	}
	var startRow int = 0
	var values *[]string
	for startRow<len(*insert.List) {
		values,startRow=insert.getValues(&fields,insert.List,commonFieldsValue,startRow,valueLength)
		valuesStr:=strings.Join(*values,",")
		sql := sqlPart1+valuesStr+";"
		_, _, err := dataRepository.execWithTx(sql, tx)
		if err != nil {
			log.Printf("Insert error: %s\n", err)
			return common.ResultSQLError
		}
	}
	result := map[string]interface{}{}
	result["count"] = startRow
	result["modelID"] = insert.ModelID

	return common.ResultSuccess
}
