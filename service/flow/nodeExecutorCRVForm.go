package flow

import (
  "time"
	"dataflow/common"
	"encoding/json"
	"log"
)

type crvFormConf struct {
	Url string `json:"url,omitempty"` 
	Location string `json:"location,omitempty"`
	Title string  `json:"title,omitempty"`
	Key string  `json:"key,omitempty"`
	Width int `json:"width,omitempty"`
	Height int  `json:"height,omitempty"`
}

type nodeExecutorCRVForm struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorCRVForm)getNodeConf()(*crvFormConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
  if err != nil {
    log.Println(err)
		return nil
  }
	
	conf:=&crvFormConf{}
  if err := json.Unmarshal(jsonStr, conf); err != nil {
    log.Println(err)
		return nil
  }

	return conf
}

func (nodeExecutor *nodeExecutorCRVForm)runStage0(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_CRV_FORM,
	}

	formConf:=nodeExecutor.getNodeConf()
	if formConf==nil {
		return nil,common.CreateError(common.ResultNodeConfigError,params)
	}

	stage:=1
	operation:=map[string]interface{}{
		"id":"", 
		"type":"open",
		"params":map[string]interface{}{
				"url":formConf.Url,
				"location":formConf.Location,
				"title":formConf.Title,
				"key":formConf.Key,
				"width":formConf.Width,
				"height":formConf.Height,
			},
		"input":map[string]interface{}{
				"flowID":instance.FlowID,
				"flowInstanceID":&(instance.InstanceID),
				"stage":&stage,
			},
    "description":"",
	}
	
	result:=&flowReqRsp{
		Operation:&operation,
	}

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=false
	node.EndTime=&endTime
	node.Output=result

	return result,nil
}

func (nodeExecutor *nodeExecutorCRVForm)runStage1(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
	req:=node.Input
	flowResult:=&flowReqRsp{
		GoOn:true,
		List:nil,
		Total:0,
		Data:req.Data,
	}
	if flowResult.Data==nil {
		flowResult.Data=&[]flowDataItem{}
	}

	if len(*flowResult.Data)==0 {
		emptyItem:=flowDataItem{
			VerifyResult:[]verifyResultItem{},
			Models:[]modelDataItem{},
		}
		(*flowResult.Data)=append(*flowResult.Data,emptyItem)
	}
	
	//将录入的数据放入data[0]的models中
	inputModelItem:=modelDataItem{
		ModelID:req.ModelID,
		List:req.List,
	}
	(*flowResult.Data)[0].Models=append((*flowResult.Data)[0].Models,inputModelItem)

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult

	return flowResult,nil
}

func (nodeExecutor *nodeExecutorCRVForm)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
	req:=node.Input

	if req.Stage==nil || *(req.Stage) == 0 {
		return nodeExecutor.runStage0(instance,node,preNode)
	}

	return nodeExecutor.runStage1(instance,node,preNode)
}