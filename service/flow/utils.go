package flow

import (
    "dataflow/common"
	"log"
	"strconv"
)

func getFieldValue(
	row map[string]interface{},
	field string)(float64,int){
	
	fieldVal,found:=row[field]
	if !found {
		log.Printf("BPModelIterator getFieldValue no field: %s!\n", field)
		return 0,common.ResultNoWriteoffField 
	}

	switch fieldVal.(type) {
	case float64:
		fVal, _ := fieldVal.(float64)
		return fVal,common.ResultSuccess
	case string:
		sVal, _ := fieldVal.(string)
		fVal, err := strconv.ParseFloat(sVal, 64)
		if err !=nil {
			log.Printf("BPModelIterator getFieldValue can not convert value to float64: %s!\n", sVal)
			return 0,common.ResultWriteoffValueTypeError
		}
		return fVal,common.ResultSuccess
	default:
		log.Printf("BPModelIterator getFieldValue not supported field type: %T!\n", fieldVal)
		return 0,common.ResultWriteoffValueTypeError
	}
}

func ternaryOperatorFloat64(b bool,tVal float64,fVal float64)(float64){
	if b {
		return tVal
	}

	return fVal
}