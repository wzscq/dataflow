package flow

import (
	"dataflow/data"
	"dataflow/common"
)

const (
	NODE_START = "start"
	NODE_QUERY = "query"
	NODE_REQUEST_QUERY = "requestQuery"
	NODE_RELATED_QUERY = "relatedQuery"
	NODE_GROUP = "fieldGroup"
	NODE_NUMERIC_GROUP = "numericGroup"
	NODE_MATCH = "match"
	NODE_VERIFY_MATCH = "verifyMatch"
	NODE_VERIFY_VALUE = "verifyValue"
	NODE_NUMBER_COMPARE = "numberCompare"
	NODE_FILTER = "filter"
	NODE_END = "end"
	NODE_SAVE_MATCHED = "saveMatched"
	NODE_SAVE_NOTMATCHED = "saveNotMatched"
	NODE_RETURN_CRVRESULT= "retrunCRVResult"
	NODE_LOG = "log"
	NODE_DEBUG = "debug"
	NODE_DELAY = "delay"
	NODE_EB_PROCESSING = "ebProcessing"
	NODE_SAVE = "save"
	NODE_DATA_TRANSFER = "dataTransfer"
	NODE_DATA_TRANSFORM = "dataTransform"
	NODE_GROUP_TRANSFORM = "groupTransform"
	NODE_SPLIT_EXTQUANTITY= "splitExtraQuantity"
	NODE_CREATE_MATCH_RESULT="createMatchResult"
	NODE_ESI = "esi"
	NODE_CRV_REQUEST = "CRVRequest"
	NODE_CRV_FORM = "CRVForm"
	NODE_TASK = "taskInfo"
	NODE_EXPORT_EXCEL = "exportExcel"
	NODE_FLOW = "callFlow"
	NODE_FLOW_ASYNC = "callFlowAsync"
	NODE_CALL_EXTERNAL_API= "callExternalAPI"
)

type nodeExecutor interface {
	run(instance *flowInstance,node,preNode *instanceNode)(*flowReqRsp,*common.CommonError)
}

func getNodeExecutor(
	node *node,
	dataRepo data.DataRepository,
	mqtt *common.MqttConf,
	redis *common.RedisConf)(nodeExecutor){
	if node.Type ==NODE_START {
		return &nodeExecutorStart{}
	} else if node.Type == NODE_END {
		return &nodeExecutorEnd{}
	} else if node.Type == NODE_QUERY {
		return &nodeExecutorQuery{
			DataRepository:dataRepo, 
			NodeConf:*node,
		}
	} else if node.Type == NODE_REQUEST_QUERY {
		return &nodeExecutorRequestQuery{
			DataRepository:dataRepo, 
			NodeConf:*node,
		}
	} else if node.Type == NODE_RELATED_QUERY {
		return &nodeExecutorRelatedQuery{
			DataRepository:dataRepo, 
			NodeConf:*node,
		}
	} else if node.Type == NODE_GROUP {
		return &nodeExecutorGroup{
			NodeConf:*node,
		}
	} else if node.Type == NODE_NUMERIC_GROUP {
		return &nodeExecutorNumericGroup{
			NodeConf:*node,
		}
	} else if node.Type == NODE_MATCH {
		return &nodeExecutorMatch{
			NodeConf:*node,
		}
	} else if node.Type == NODE_VERIFY_MATCH {
		return &nodeExecutorVerifyMatch{
			NodeConf:*node,
		}
	} else if node.Type == NODE_VERIFY_VALUE {
		return &nodeExecutorVerifyValue{
			NodeConf:*node,
		}
	}else if node.Type == NODE_FILTER {
		return &nodeExecutorFilter{
			NodeConf:*node,
		}
	} else if node.Type == NODE_SAVE_MATCHED {
		return &nodeExecutorSaveMatched{
			DataRepository:dataRepo,
			NodeConf:*node,
		}
	} else if node.Type == NODE_SAVE_NOTMATCHED {
		return &nodeExecutorSaveNotMatched{
			DataRepository:dataRepo,
			NodeConf:*node,
		}
	} else if node.Type == NODE_LOG {
		return &nodeExecutorLog{
			NodeConf:*node,
		}
	} else if node.Type == NODE_DEBUG {
		return &nodeExecutorDebug{
			NodeConf:*node,
			Mqtt:mqtt,
		}
	} else if node.Type == NODE_DELAY {
		return &nodeExecutorDelay{
			NodeConf:*node,
		}
	} else if node.Type == NODE_EB_PROCESSING {
		return &nodeExecutorEBProcessing{
			NodeConf:*node,
		}
	} else if node.Type == NODE_SAVE {
		return &nodeExecutorSave{
			DataRepository:dataRepo,
			NodeConf:*node,
		}
	} else if node.Type == NODE_DATA_TRANSFER {
		return &nodeExecutorDataTransfer{
			NodeConf:*node,
		}
	} else if node.Type == NODE_DATA_TRANSFORM {
		return &nodeExecutorDataTransform{
			NodeConf:*node,
		}
	} else if node.Type == NODE_NUMBER_COMPARE {
		return &nodeExecutorNumberCompare{
			NodeConf:*node,
		}
	} else if node.Type == NODE_GROUP_TRANSFORM {
		return &nodeExecutorGroupTransform{
			NodeConf:*node,
		}
	} else if node.Type == NODE_SPLIT_EXTQUANTITY {
		return &nodeExecutorSplitExtraQuantity{
			NodeConf:*node,
		}
	} else if node.Type == NODE_CREATE_MATCH_RESULT {
		return &nodeExecutorCreateMatchResult{
			NodeConf:*node,
		}
	} else if node.Type == NODE_ESI {
		return &nodeExecutorESI{
			NodeConf:*node,
			DataRepository:dataRepo,
		}
	} else if node.Type == NODE_CRV_REQUEST {
		return &nodeExecutorCRVRequest{
			NodeConf:*node,
		}
	} else if node.Type == NODE_FLOW {
		return &nodeExecutorFlow{
			NodeConf:*node,
			DataRepository:dataRepo,
			Mqtt:mqtt,
			Redis:redis,
		}
	} else if node.Type == NODE_RETURN_CRVRESULT {
		return &nodeExecutorReturnCRVResult{
			NodeConf:*node,
		}
	} else if node.Type == NODE_CRV_FORM {
		return &nodeExecutorCRVForm{
			NodeConf:*node,
		}
	} else if node.Type == NODE_FLOW_ASYNC {
		return &nodeExecutorFlowAsync{
			NodeConf:*node,
			Mqtt:mqtt,
		}
	} else if node.Type == NODE_TASK {
		return &nodeExecutorTask{
			NodeConf:*node,
			Mqtt:mqtt,
			Redis:redis,
			DataRepository:dataRepo,
		}
	} else if node.Type == NODE_EXPORT_EXCEL {
		return &nodeExecutorExportExcel{
			NodeConf:*node,
		}
	} else if node.Type == NODE_CALL_EXTERNAL_API {
		return &nodeExecutorCallExternalAPI{
			NodeConf:*node,
		}
	}
	return nil
}