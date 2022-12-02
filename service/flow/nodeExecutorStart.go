package flow

import (
    "time"
	"dataflow/common"
)

type nodeExecutorStart struct {

}

func (nodeExecutor *nodeExecutorStart)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=node.Input

	return node.Output,nil
}