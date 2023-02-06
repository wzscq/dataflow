package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"os"
	"io/ioutil"
)

type nodeExecutorLog struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorLog)logInstanceNode(instance *flowInstance,node *instanceNode){
	//打印一下流的实例内容
	jsonStr, err := json.MarshalIndent(node, "", "    ")
    if err != nil {
        log.Println(err)
    } else {

		path := "apps/"+instance.AppDB+"/instances/"+instance.InstanceID
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		} else {
			fileName:=path+"/node_"+node.ID+".json"
			ioutil.WriteFile(fileName, jsonStr, 0644)
		}
	}
	jsonStr=nil
}

func (nodeExecutor *nodeExecutorLog)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=node.Input
	node.Input=nil
	preNode.Input=nil
	nodeExecutor.logInstanceNode(instance,node)

	return node.Output,nil
}