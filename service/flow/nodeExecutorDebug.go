package flow

import (
    "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type sendNode struct {
	ID string `json:"id"`
	Completed bool `json:"completed"`
	StartTime string `json:"startTime"`
	EndTime *string `json:"endTime,omitempty"`
	UserID string `json:"userID,omitempty"`
	NodeType string `json:"nodeType,omitempty"`
	Input *flowReqRsp `json:"input,omitempty"`
	Output *flowReqRsp `json:"output,omitempty"`
	Priority int64 `json:"priority,omitempty"`
	FlowID string `json:"flowID"`
}

type nodeExecutorDebug struct {
	NodeConf node
	Mqtt *common.MqttConf
}

func (nodeExecutor *nodeExecutorDebug) connectHandler(client mqtt.Client){
	log.Println("connectHandler connect status: ",client.IsConnected())
}

func (nodeExecutor *nodeExecutorDebug) connectLostHandler(client mqtt.Client, err error){
	log.Println("connectLostHandler connect status: ",client.IsConnected(),err)
}

func (nodeExecutor *nodeExecutorDebug) messagePublishHandler(client mqtt.Client, msg mqtt.Message){
	log.Println("messagePublishHandler topic: ",msg.Topic())
}

func (nodeExecutor *nodeExecutorDebug)getMqttClient(debugID string)(*mqtt.Client){
	broker := nodeExecutor.Mqtt.Broker //"121.36.192.249"
	port := nodeExecutor.Mqtt.Port //1983
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d",broker,port))
	opts.SetClientID("flow_node_debug_"+debugID)
	opts.SetUsername(nodeExecutor.Mqtt.User)
	opts.SetPassword(nodeExecutor.Mqtt.Password)
	opts.SetDefaultPublishHandler(nodeExecutor.messagePublishHandler)
	opts.OnConnect = nodeExecutor.connectHandler
	opts.OnConnectionLost = nodeExecutor.connectLostHandler
	client:=mqtt.NewClient(opts)
	if token:=client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error)
		return nil
	}
	return &client
}

func (nodeExecutor *nodeExecutorDebug)SendDebugMessage(instance *flowInstance,node *sendNode){
	log.Println("SendDebugMessage topic: flowdebug/"+(*instance.DebugID))
	//打印一下流的实例内容
	jsonStr, err := json.MarshalIndent(node, "", "    ")
    if err != nil {
        log.Println(err)
		return
    }
	client:=nodeExecutor.getMqttClient(*instance.DebugID)
	if client == nil {
		return
	}
	defer (*client).Disconnect(250)

	token:=(*client).Publish("flowdebug/"+(*instance.DebugID),0,false,string(jsonStr))
	token.Wait()
}

func (nodeExecutor *nodeExecutorDebug)getSendData(data *flowReqRsp)(*flowReqRsp){
	sendData:=&flowReqRsp{
		FlowID:data.FlowID,
		UserID:data.UserID,
		AppDB:data.AppDB,
		DebugID:data.DebugID,
		Data:&[]flowDataItem{},
	}
	if data.Data!=nil && len(*data.Data)>0 {
		sendData.Total=len(*data.Data)
		for index,dataItem:=range(*data.Data){
			if index<5 {
				sendDataItem:=flowDataItem{
					VerifyResult:dataItem.VerifyResult,
					Models:[]modelDataItem{},
				}

				for _,modelData:=range(dataItem.Models) {
					if modelData.List !=nil {
						sendModelData:=modelDataItem{
							ModelID:modelData.ModelID,
							List:&[]map[string]interface{}{},
							Total:len(*modelData.List),
						}
						for rowIdx,dataRow:=range(*modelData.List){
							if rowIdx<100 {
								(*sendModelData.List)=append((*sendModelData.List),dataRow)
							} else {
								break
							}
						}
						sendDataItem.Models=append(sendDataItem.Models,sendModelData)
					}
				}

				(*sendData.Data)=append((*sendData.Data),sendDataItem)
			} else {
				break
			}
		}
	}

	return sendData
}

func (nodeExecutor *nodeExecutorDebug)getSendObject(node *instanceNode,flowID string)(*sendNode){
	sendNode:=&sendNode{
		ID:node.ID,
		Completed:node.Completed,
		StartTime:node.StartTime,
		EndTime:node.EndTime,
		UserID:node.UserID,
		NodeType:node.NodeType,
		FlowID:flowID,
	}

	if node.Input!=nil {
		sendNode.Input=nodeExecutor.getSendData(node.Input)
	}

	if node.Output!=nil {
		sendNode.Output=nodeExecutor.getSendData(node.Output)
	}

	return sendNode
}

func (nodeExecutor *nodeExecutorDebug)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	if instance.DebugID!=nil && len(*instance.DebugID)>0 {
		sendObj:=nodeExecutor.getSendObject(node,instance.FlowID)
		nodeExecutor.SendDebugMessage(instance,sendObj)
	}

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=node.Input
	node.Input=nil
	return node.Output,nil
}