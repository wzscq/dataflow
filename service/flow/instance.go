package flow

import (
	"log"
	"dataflow/common"
	"dataflow/data"
	"sort"
	
)

type flowInstance struct {
	 AppDB string `json:"appDB"`
	 FlowID string `json:"flowID"`
	 InstanceID string `json:"instanceID"`
	 UserID string  `json:"UserID"`
	 FlowConf *flowConf `json:"flowConf,omitempty"`
	 CompletedNodes []*instanceNode `json:"completedNodes"`
	 WaitingNodes []*instanceNode `json:"waitingNodes"`
	 Completed bool `json:"completed"`
	 StartTime string `json:"startTime"`
	 EndTime *string `json:"endTime,omitempty"`
	 DebugID *string `json:"debugID,omitempty"`
	 InstanceRepository FlowInstanceRepository
}

func (flow *flowInstance)getCurrentNode(flowRep* flowReqRsp)(*instanceNode){
	nodeCount:=len(flow.WaitingNodes)
	if nodeCount>0 {
		currentNode:=flow.WaitingNodes[nodeCount-1]
		//将Node的Input替换成最新的请求参数
		flowRep.Data=currentNode.Input.Data
		currentNode.Input=flowRep
		return currentNode
	}
	return nil
}

func (flow *flowInstance)sortInstanceNodes(instanceNodes *[]*instanceNode){
	sort.Slice(*instanceNodes, func (i, j int) bool {
		return (*instanceNodes)[i].Priority < (*instanceNodes)[j].Priority
	})
}

func (flow *flowInstance)getStartNode(flowRep* flowReqRsp)(*instanceNode){
	fllowingNodes:=[]*instanceNode{}
	for _, nodeItem := range (flow.FlowConf.Nodes) {
		if nodeItem.Type == NODE_START {
			priority:=flow.getNodePriority(&nodeItem)
			instanceNode:=createInstanceNode(nodeItem.ID,nodeItem.Type,priority,flowRep)
			fllowingNodes=append(fllowingNodes,instanceNode)
		}
	}
	//按照优先级对节点排序，数字高的先执行，排在数组的后面
	flow.sortInstanceNodes(&fllowingNodes)

	flow.WaitingNodes=append(flow.WaitingNodes,fllowingNodes...)
	nodeCount:=len(flow.WaitingNodes)
	if nodeCount>0 {
		return flow.WaitingNodes[nodeCount-1]
	}
	return nil
}

func  (flow *flowInstance)getNodePriority(node *node)(int64){
	dataMap,ok:=node.Data.(map[string]interface{})
	if !ok {
		log.Println("getNodePriority node data is not a map")
		return 0
	}

	priorityItem,ok:=dataMap["priority"]
	if !ok {
		log.Println("getNodePriority node data map not contain priority")
		log.Println(dataMap)
		return 0
	}

	priority,ok:=priorityItem.(float64)
	if !ok {
		log.Println("getNodePriority node data priority is not a int64 value")
		log.Printf("%T",priorityItem)
		return 0
	}
	log.Printf("getNodePriority priority:%d",int64(priority))
	return int64(priority)
}

func (flow *flowInstance)getFllowingNodes(id string,flowRep* flowReqRsp)(*[]*instanceNode){
	fllowingNodes:=[]*instanceNode{}
	for _, edgeItem := range (flow.FlowConf.Edges) {
		if edgeItem.Source == id {
			nodeItem:=flow.getNodeConfig(edgeItem.Target)
			priority:=flow.getNodePriority(nodeItem)
			instanceNode:=createInstanceNode(nodeItem.ID,nodeItem.Type,priority,flowRep)
			fllowingNodes=append(fllowingNodes,instanceNode)
		}
	}
	//按照优先级排序
	if len(fllowingNodes)>1 {

		//jsonStr, _:= json.MarshalIndent(fllowingNodes, "", "    ")
		//log.Println("getFllowingNodes before sort:",string(jsonStr))
		flow.sortInstanceNodes(&fllowingNodes)
		//jsonStr, _= json.MarshalIndent(fllowingNodes, "", "    ")
		//log.Println("getFllowingNodes after sort:",string(jsonStr))
	}

	return &fllowingNodes
}

func (flow *flowInstance)getNextNode(currentNode *instanceNode,flowRep* flowReqRsp)(*instanceNode){
	//如果当前节点返回要求继续执行，先检查当前节点是否有后续节点，如果有则优先执行后续节点
	if flowRep.GoOn==true {
		fllowingNodes:=flow.getFllowingNodes(currentNode.ID,flowRep)
		flow.WaitingNodes=append(flow.WaitingNodes,(*fllowingNodes)...)
	}
	
	nodeCount:=len(flow.WaitingNodes)
	if nodeCount>0 {
		return flow.WaitingNodes[nodeCount-1]
	}
	return nil
}

func (flow *flowInstance)addCompletedNode(node *instanceNode){
	//这里暂时不记录已经执行完成的node
	//flow.CompletedNodes=append(flow.CompletedNodes,node)
	nodeCount:=len(flow.WaitingNodes)
	flow.WaitingNodes=flow.WaitingNodes[0:nodeCount-1]
}

func (flow *flowInstance)updateCurrentNode(node *instanceNode){
	nodeCount:=len(flow.WaitingNodes)
	flow.WaitingNodes[nodeCount-1]=node
}

func (flow *flowInstance)getNodeConfig(id string)(node *node){
	for _, nodeItem := range (flow.FlowConf.Nodes) {
		if nodeItem.ID == id {
			return &nodeItem
		}
	}
	return nil
}

func (flow *flowInstance)runNode(
	dataRepo data.DataRepository,
	node,preNode *instanceNode,
	mqtt *common.MqttConf)(*flowReqRsp,*common.CommonError){
	//根据节点类型，找到对应的节点，然后执行节点
	log.Printf("flowInstance runNode id: %s \n",node.ID)
	nodeCfg:=flow.getNodeConfig(node.ID)
	if nodeCfg==nil {
		log.Println("can not find the node config with id: ",node.ID)
		params:=map[string]interface{}{
			"nodeID":node.ID,
		}		
		return nil,common.CreateError(common.ResultNoNodeOfGivenID,params)
	}
	log.Printf("flowInstance runNode id: %s type: %s \n",node.ID,nodeCfg.Type)
	executor:=getNodeExecutor(nodeCfg,dataRepo,mqtt)
	if executor==nil {
		log.Println("can not find the node executor with type: ",nodeCfg.Type)
		params:=map[string]interface{}{
			"nodeID":node.ID,
			"nodeType":nodeCfg.Type,
		}
		return nil,common.CreateError(common.ResultNoExecutorForNodeType,params)
	}
	return executor.run(flow,node,preNode)
}

func (flow *flowInstance)push(dataRepo data.DataRepository,flowRep* flowReqRsp,mqtt *common.MqttConf)(*flowReqRsp,*common.CommonError){
	log.Println("start flowInstance push")

	if flow.InstanceRepository!=nil {
		err:=flow.InstanceRepository.loadInstance(flow)
		if err!=nil {
			log.Println(err)
		}
	}

	currentNode:=flow.getCurrentNode(flowRep)
	if currentNode==nil {
		currentNode=flow.getStartNode(flowRep)
	}

	var preNode *instanceNode
	//循环执行所有同步的node
	for currentNode!=nil {
		log.Printf("flowInstance push currentNode id %s \n",currentNode.ID)
		result,err:=flow.runNode(dataRepo,currentNode,preNode,mqtt)
		if err!= nil {
			return nil,err
		}

		if result.Over == true {
			return result,nil
		}

		//如果执行完，就拿下一个节点继续执行
		if currentNode.Completed {
			flow.addCompletedNode(currentNode)
			//logInstanceNode(flow,currentNode)
			preNode=currentNode
			currentNode=flow.getNextNode(currentNode,result)
		} else {
			//如果没有执行完，说明这个节点是异步节点，直接将结果返回，待后续触发
			//更新节点状态
			flow.updateCurrentNode(currentNode)
			//保存执行调用栈
			err:=flow.InstanceRepository.saveInstance(flow)
			if err!=nil {
				log.Println(err)

			}
			return result,nil
		}
	}
	
	log.Println("end flowInstance push")
	return nil,nil
}
