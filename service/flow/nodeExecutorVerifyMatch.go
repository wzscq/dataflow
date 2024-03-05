package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
)

/**
目前支持两种类型的校验规则
1、存在性判断，检查左右表是否都有数据
2、数值判断，检查同一组中不同模型对应字段的汇总数据是否满足给定操作符要求，支持的操作符：
	=：指定字段值相等，如果给了容差，差值在容差范围内也符合要求
	<: 前面的表字段值<后面表的字段值
	>：前面的字段值>后面表的字段值
	数值判断规则可以配置多个，对不同的字段做校验，当一个字段的校验失败时仍然会执行其它数值判断的校验
校验结果中仅记录校验失败节点和失败信息
**/
type existRuleItem struct {
	ModelID string `json:"modelID"`
	Message string `json:"message"`
}

type matchVerifyConf struct {
	VerfiyID string `json:"verifyID"`
	FailureResult string `json:"failureResult"`
    SuccessfulResult string `json:"successfulResult"`
	Rules *[]existRuleItem `json:"rules"`
}

type nodeExecutorVerifyMatch struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorVerifyMatch)getVerifyConf()(*matchVerifyConf){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok || mapData["rules"]==nil {
		return nil
	}
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	conf:=matchVerifyConf{}
    if err := json.Unmarshal(jsonStr, &conf); err != nil {
        log.Println(err)
		return nil
    }

	return &conf
}

func (nodeExecutor *nodeExecutorVerifyMatch)copyDataItem(
	item,newItem *flowDataItem){
	newItem.VerifyResult=append(newItem.VerifyResult,item.VerifyResult...)
	newItem.Models=append(newItem.Models,item.Models...)
}

func (nodeExecutor *nodeExecutorVerifyMatch)verify(verifyConf *matchVerifyConf,dataItem *flowDataItem){
	verifyItem:=verifyResultItem{
		VerfiyID:verifyConf.VerfiyID,
		VerfiyType:"matchVerify",
		Result:verifyConf.SuccessfulResult,
		Message:"",
	}
	log.Println("start verify")
	for _,ruleItem:= range (*verifyConf.Rules) {
		hasRuleModel:=false
		for _,modelItem:= range (dataItem.Models) {
			if *modelItem.ModelID == ruleItem.ModelID {
				hasRuleModel=true
				if modelItem.List == nil || len(*modelItem.List)==0 {
					verifyItem.Result=verifyConf.FailureResult
					verifyItem.Message=ruleItem.Message
				}
			}
		}
		if hasRuleModel==false {
			verifyItem.Result=verifyConf.FailureResult
			verifyItem.Message=ruleItem.Message
		}
	}
	log.Println("end verify")
	log.Println(verifyItem)
	dataItem.VerifyResult=append(dataItem.VerifyResult,verifyItem)
}

func (nodeExecutor *nodeExecutorVerifyMatch)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID,
		FlowInstanceID:req.FlowInstanceID,
		Stage:req.Stage,
		DebugID:req.DebugID,
		UserRoles:req.UserRoles,
		GlobalFilterData:req.GlobalFilterData,
		UserID:req.UserID,
		AppDB:req.AppDB,
		Token:req.Token,
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
		SelectAll:req.SelectAll,
		GoOn:true,
	}
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_VERIFY_MATCH,
	}
	//加载节点配置
	verifyConf:=nodeExecutor.getVerifyConf()
	if verifyConf==nil {
		log.Printf("nodeExecutorVerifyExist run get node config error\n")
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if node.Input.Data==nil || len(*node.Input.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=node.Input
		return node.Input,nil
	}

	//执行校验逻辑
	resultData:=make([]flowDataItem,len(*req.Data))
	//这里的分组操作是在数据已经分组的基础上再次分组，分组数据不能跨原来的分组
	//遍历每个数据分组
	for index,item:= range (*req.Data) {
		//复制原始数据
		nodeExecutor.copyDataItem(&item,&(resultData[index]))
		nodeExecutor.verify(verifyConf,&(resultData[index]))
	}

	flowResult.Data=&resultData
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	return flowResult,nil
}