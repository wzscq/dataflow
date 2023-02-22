package flow

import (
  "time"
	"dataflow/common"
	"dataflow/data"
	"encoding/json"
	"log"
)

type nodeFlowConf struct {
	FlowID string `json:"flowID"`
}

type nodeExecutorFlow struct {
	NodeConf node
	DataRepository data.DataRepository
	Mqtt *common.MqttConf
}

func (nodeExecutor *nodeExecutorFlow)getNodeConf()(*nodeFlowConf){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok {
		return nil
	}
	jsonStr, err := json.Marshal(mapData)
  if err != nil {
    log.Println(err)
		return nil
  }
	conf:=nodeFlowConf{}
  if err := json.Unmarshal(jsonStr, &conf); err != nil {
    log.Println(err)
		return nil
  }

	return &conf
}

func (nodeExecutor *nodeExecutorFlow)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_FLOW,
	}

	req:=node.Input

	//加载节点配置
	nodeConf:=nodeExecutor.getNodeConf()
	if nodeConf==nil {
		log.Printf("nodeExecutorFlow run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeFilterConfigError,params)
	}

	//调用一个外部的流
	//目前流的执行本身是没有返回值的，所以现在可以不用考虑流的返回值的处理
	//调用时，直接使用当前节点的Input作为子流程的请求参数
	//创建流,子流程不支持直接传递流的配置
	flowInstance,errorCode:=createInstance(
		req.AppDB,
		nodeConf.FlowID,
		req.UserID,
		nil,
		instance.DebugID,
		instance.TaskID,
		instance.TaskStep,
		nil,nil)
		
	if errorCode!=common.ResultSuccess {
		return node.Input,common.CreateError(errorCode,params)
	}
	//执行流
	result,err:=flowInstance.push(nodeExecutor.DataRepository,req,nodeExecutor.Mqtt)
	if err!=nil {
		return result,err
	}

	//将流的返回值作为节点的输出,将flow的标识替换回主流程的标识
	if result!=nil {
		result.FlowID=req.FlowID
		result.FlowInstanceID=req.FlowInstanceID
		result.Over=false //子流程结束，主流程继续执行

		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=result
		return result,nil
	}

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=req
	return req,nil
}