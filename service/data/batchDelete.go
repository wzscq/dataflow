package data

import (
	"dataflow/common"
	"database/sql"
)

type BatchDelete struct {
	ModelID         string    `json:"modelID"`
	SelectedRowKeys *[]string `json:"selectedRowKeys"`
	AppDB           string    `json:"appDB"`
	Filter     *map[string]interface{} `json:"filter"`
	SelectAll        bool     `json:"selectedAll"`
	Fields     *[]Field       `json:"fields"`
}

func (delete *BatchDelete) getWhere() (string, int) {
	var filter *map[string]interface{}
	if delete.SelectAll == false {
		if delete.SelectedRowKeys == nil || len(*delete.SelectedRowKeys) == 0 {
			return "1==2", common.ResultSuccess
		}

		filter=&map[string]interface{}{
			"id": map[string]interface{}{
				Op_in: *delete.SelectedRowKeys,
			},
		}
	} else {
		filter=delete.Filter
	}

	return FilterToSQLWhere(filter)
}

func (delete *BatchDelete) Delete(dataRepository DataRepository, tx *sql.Tx) (*map[string]interface{}, int) {
	//获取所有待删数据查询条件
	sWhere,errorCode:=delete.getWhere()
	if errorCode != common.ResultSuccess {
		return nil, errorCode
	}

	sql := "delete from " + delete.AppDB + "." + delete.ModelID + " where " + sWhere

	_, rowCount, err := dataRepository.execWithTx(sql, tx)
	if err != nil {
		return nil, common.ResultSQLError
	}
	result := map[string]interface{}{}
	result["count"] = rowCount
	result["modelID"] = delete.ModelID

	return &result, errorCode
}
