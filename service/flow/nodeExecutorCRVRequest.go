package flow

import (
  "time"
	"dataflow/common"
	"log"
	"encoding/json"
	"dataflow/data"
)

type testData struct {
	ModelID *string `json:"modelID"`
	Filter *map[string]interface{} `json:"filter"`
	List *[]map[string]interface{} `json:"list"`
	GlobalFilterData *map[string]interface{} `json:"globalFilterData"`
	FilterData *[]data.FilterDataItem `json:"filterData"`
	SelectedRowKeys *[]string `json:"selectedRowKeys"`
	Pagination *data.Pagination `json:"pagination"`
	SelectAll bool `json:"selectedAll"`
	Token string `json:"token"`
}

type nodeExecutorCRVRequest struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorCRVRequest)loadTestData()(*testData){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	testDataMap,ok:=mapData["testData"]
	if !ok {
		return nil
	}

	jsonStr, err := json.Marshal(testDataMap)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println(string(jsonStr))
	testData:=&testData{}
  if err := json.Unmarshal(jsonStr, testData); err != nil {
    log.Println(err)
		return nil
  }

	return testData
}

func (nodeExecutor *nodeExecutorCRVRequest)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorCRVRequest run start")

	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID,
		FlowInstanceID:req.FlowInstanceID,
		Stage:req.Stage,
		DebugID:req.DebugID,
		UserRoles:req.UserRoles,
		UserID:req.UserID,
		GlobalFilterData:req.GlobalFilterData,
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

	//如果是调试模式，择用测试数据填充List
	if instance.DebugID!=nil && len(*instance.DebugID)>0 {
		testData:=nodeExecutor.loadTestData()
		if testData != nil {
			flowResult.ModelID=testData.ModelID
			flowResult.Filter=testData.Filter
			flowResult.List=testData.List
			flowResult.GlobalFilterData=testData.GlobalFilterData
			flowResult.FilterData=testData.FilterData
			flowResult.SelectedRowKeys=testData.SelectedRowKeys
			flowResult.Pagination=testData.Pagination
			flowResult.SelectAll=testData.SelectAll
			flowResult.Token=testData.Token
		}
	}

	/*params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_CRV_REQUEST,
	}*/

	//将接口传入的数据放入data中
	modelID:=*flowResult.ModelID
	modelDatas:=[]modelDataItem{
		modelDataItem{
			ModelID:&modelID,
			Filter:flowResult.Filter,
			List:flowResult.List,
			GlobalFilterData:flowResult.GlobalFilterData,
			FilterData:flowResult.FilterData,
			SelectedRowKeys:flowResult.SelectedRowKeys,
			Pagination:flowResult.Pagination,
			SelectAll:flowResult.SelectAll,
		},
	}

	//如果之前的查询中已经存在数据项，则将当前查询的项目和并到之前查询数据的第0个项目上
	data:=req.Data
	if data ==nil || len(*data)==0 {
		data=&[]flowDataItem{
			flowDataItem{
				Models:modelDatas,
			},
		}
	} else {
		(*data)[0].Models=append((*data)[0].Models,modelDatas...)
	}

	flowResult.Data=data
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult

	return flowResult,nil
}