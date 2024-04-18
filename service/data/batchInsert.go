package data

import (
	"dataflow/common"
	"database/sql"
	"log"
	"strings"
	"fmt"
	"time"
)

type BatchInsert struct {
	ModelID         string    `json:"modelID"`
	AppDB           string    `json:"appDB"`
	UserID string `json:"userID"`
	List *[]map[string]interface{} `json:"list"` 
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

func (insert *BatchInsert) getValues(fields *[]string, list *[]map[string]interface{},commonFieldsValue string) (*[]string) {
	values:=[]string{}
	for _,data:=range *list {
		value:=insert.getRowValues(fields,data)
		if value==nil {
			return nil
		}
		values=append(values,"("+*value+","+commonFieldsValue+")")
	}
	return &values
}

func (insert *BatchInsert) Insert(dataRepository DataRepository, tx *sql.Tx) (int) {
	//获取所有待删数据查询条件
	if insert.List == nil || len(*insert.List) == 0 {
		return common.ResultSuccess
	}

	commonFields,commonFieldsValue:=GetCreateCommonFieldsValues(insert.UserID)
	fields:=insert.getFields((*insert.List)[0])
	values:=insert.getValues(&fields,insert.List,commonFieldsValue)
	valuesStr:=strings.Join(*values,",")
	fieldsStr:=strings.Join(fields,",")
	sql := "insert into " + insert.AppDB + "." + insert.ModelID + "("+fieldsStr+","+commonFields+") values "+valuesStr+";"
	_, rowCount, err := dataRepository.execWithTx(sql, tx)
	if err != nil {
		log.Printf("Insert error: %s\n", err)
		return common.ResultSQLError
	}
	result := map[string]interface{}{}
	result["count"] = rowCount
	result["modelID"] = insert.ModelID

	return common.ResultSuccess
}
