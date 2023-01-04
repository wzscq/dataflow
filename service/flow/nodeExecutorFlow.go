package flow

import (
    "time"
	"dataflow/common"
)

type nodeFlowConf struct {
	FlowID string `json:"flowID"`
}

type nodeExecutorFlow struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorFlow)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	//调用一个外部的流
	//目前流的执行本身是没有返回值的，所以现在可以不用考虑流的返回值的处理
	//后续如果增加了流的返回值节点，可以通过返回值节点来实现流的返回值
	//情况下只需要考虑调用流就可以了
	
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=node.Input

	return node.Output,nil
}