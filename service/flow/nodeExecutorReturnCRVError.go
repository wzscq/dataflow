package flow

import (
	"log"
	"strconv"
	"dataflow/common"
	"encoding/json"
)

type CRVErrorConf struct {
	ErrorCode string `json:"errorCode"`
	Message string `json:"message"`
}

type nodeExecutorReturnCRVError struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorReturnCRVError)getNodeConf()(*CRVErrorConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
  	if err != nil {
    	log.Println(err)
		return nil
  	}

	conf:=&CRVErrorConf{}
  	if err := json.Unmarshal(jsonStr, conf); err != nil {
    	log.Println(err)
		return nil
  	}

	return conf
}

func (nodeExecutor *nodeExecutorReturnCRVError)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	flowResult:=&flowReqRsp{
		GoOn:true,
		Over:true,
		List:nil,
		Total:0,
	}
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_RETURN_CRVERROR,
	}

	//获取节点配置
	conf:=nodeExecutor.getNodeConf()
	if conf==nil {
		log.Printf("nodeExecutorReturnCRVError run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeConfigError,params)
	}

	//将ErrorCode从string转换为int
	errorCode,err:=strconv.Atoi(conf.ErrorCode)
	if err!=nil {
		log.Println(err)
		log.Printf("nodeExecutorReturnCRVError can not convert errorcode from string to int.\n")
		return node.Input,common.CreateError(common.ResultNodeConfigError,params)
	}

	retErr:=&common.CommonError{
		ErrorCode:errorCode,
		Message:conf.Message,
	}

	return flowResult,retErr
}