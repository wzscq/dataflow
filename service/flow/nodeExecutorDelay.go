package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"strconv"
)

type nodeDelayConf struct {
	Seconds string `json:"seconds"`
}

type nodeExecutorDelay struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorDelay)getNodeConf()(*nodeDelayConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	log.Println(string(jsonStr))
	conf:=&nodeDelayConf{}
    if err := json.Unmarshal(jsonStr, conf); err != nil {
        log.Println(err)
		return nil
    }

	return conf
}

func (nodeExecutor *nodeExecutorDelay)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_DELAY,
	}
	//加载节点配置
	nodeConf:=nodeExecutor.getNodeConf()
	if nodeConf==nil {
		log.Printf("nodeExecutorDelay run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeGroupConfigError,params)
	}

	seconds,_:=strconv.ParseInt(nodeConf.Seconds,0,64)

	time.Sleep(time.Duration(seconds)*time.Second)	
		
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=node.Input

	return node.Output,nil
}