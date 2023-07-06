package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
)

type filterItem struct {
	VerifyID string `json:"verifyID,omitempty"`
	Result string `json:"result,omitempty"`
}

type filterConf struct {
	Filter []filterItem `json:"filter,omitempty"`
}

type nodeExecutorFilter struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorFilter)getFilterConf()(*filterConf){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok {
		return nil
	}
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	conf:=filterConf{}
    if err := json.Unmarshal(jsonStr, &conf); err != nil {
        log.Println(err)
		return nil
    }

	return &conf
}

func (nodeExecutor *nodeExecutorFilter)filter(
	dataItem *flowDataItem,
	filterConf *filterConf)(bool){

	//log.Println("filter",filterConf.Filter,dataItem.VerifyResult)

	for _,filterItem:=range filterConf.Filter {
		for _,verifyItem:= range dataItem.VerifyResult {
			//log.Println("filter",verifyItem.VerfiyID,filterItem.VerifyID,verifyItem.Result,filterItem.Result)
			if verifyItem.VerfiyID==filterItem.VerifyID && verifyItem.Result==filterItem.Result {
				return true
			}
		}
	}

	return false
}

func (nodeExecutor *nodeExecutorFilter)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_FILTER,
	}
	//加载节点配置
	filterConf:=nodeExecutor.getFilterConf()
	if filterConf==nil {
		log.Printf("nodeExecutorFilter run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeFilterConfigError,params)
	}

	if node.Input.Data==nil || len(*node.Input.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=node.Input
		return node.Input,nil
	}

	//执行校验逻辑
	resultData:=make([]flowDataItem,len(*node.Input.Data))
	
	resultCount:=0
	for _,item:= range (*node.Input.Data) {
		if nodeExecutor.filter(&item,filterConf) {
			resultData[resultCount]=item
			resultCount+=1
		}
	}

	resultData=resultData[0:resultCount]
	
	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID,
		FlowInstanceID:req.FlowInstanceID,
		Stage:req.Stage,
		DebugID:req.DebugID,
		UserRoles:req.UserRoles,
		UserID:req.UserID,
		AppDB:req.AppDB,
		FlowConf:req.FlowConf,
		ModelID:req.ModelID,
		ViewID:req.ViewID,
		FilterData:req.FilterData,
		Filter:req.Filter,
		List:req.List,
		Total:req.Total,
		SelectedRowKeys:req.SelectedRowKeys,
		Pagination:req.Pagination,
		Operation:req.Operation,
		GoOn:true,
	}
	flowResult.Data=&resultData

	//如果没有数据就不再往下执行了
	if len(resultData)==0 {
		flowResult.GoOn=false
	}
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	return flowResult,nil
}