package flow

import (
  "time"
	"log"
	"dataflow/common"
	"encoding/json"
)

type CRVResultConf struct {
	ModelID string `json:"modelID"`
}

type nodeExecutorReturnCRVResult struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorReturnCRVResult)getNodeConf()(*CRVResultConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
  if err != nil {
    log.Println(err)
		return nil
  }

	conf:=&CRVResultConf{}
  if err := json.Unmarshal(jsonStr, conf); err != nil {
    log.Println(err)
		return nil
  }

	return conf
}

func (nodeExecutor *nodeExecutorReturnCRVResult)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	req:=node.Input
	flowResult:=&flowReqRsp{
		GoOn:true,
		Over:true,
		List:nil,
		Total:0,
	}
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_RETURN_CRVRESULT,
	}

	//获取节点配置
	conf:=nodeExecutor.getNodeConf()
	if conf==nil {
		log.Printf("nodeExecutorReturnCRVResult run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeConfigError,params)
	}

	flowResult.ModelID=&conf.ModelID

	//获取返回数据
	if req.Data==nil || len(*req.Data)==0 {
		log.Printf("nodeExecutorReturnCRVResult run end with empty input data.\n")
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=flowResult
		return flowResult,nil
	}

	//从第一个data item中获取对应model的数据
	for _,modelDataItem:=range (*req.Data)[0].Models {
		if *modelDataItem.ModelID == conf.ModelID {
			flowResult.List=modelDataItem.List
			flowResult.Total=modelDataItem.Total
			break
		}
	}
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult

	return flowResult,nil
}