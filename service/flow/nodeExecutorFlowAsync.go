package flow

import (
  "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type nodeFlowAsyncConf struct {
	FlowID string `json:"flowID"`
}

type nodeExecutorFlowAsync struct {
	NodeConf node
	Mqtt *common.MqttConf
}

func (nodeExecutor *nodeExecutorFlowAsync) connectHandler(client mqtt.Client){
	log.Println("connectHandler connect status: ",client.IsConnected())
}

func (nodeExecutor *nodeExecutorFlowAsync) connectLostHandler(client mqtt.Client, err error){
	log.Println("connectLostHandler connect status: ",client.IsConnected(),err)
}

func (nodeExecutor *nodeExecutorFlowAsync) messagePublishHandler(client mqtt.Client, msg mqtt.Message){
	log.Println("messagePublishHandler topic: ",msg.Topic())
}

func (nodeExecutor *nodeExecutorFlowAsync)getNodeConf()(*nodeFlowAsyncConf){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok {
		return nil
	}
	jsonStr, err := json.Marshal(mapData)
  if err != nil {
    log.Println(err)
		return nil
  }
	conf:=nodeFlowAsyncConf{}
  if err := json.Unmarshal(jsonStr, &conf); err != nil {
    log.Println(err)
		return nil
  }

	return &conf
}

//将输入对象转换为json，然后再反序列化回来，实现对象的深度复制
func (nodeExecutor *nodeExecutorFlowAsync)copyFlowReqRsp(input *flowReqRsp)(*flowReqRsp,int){
	jsonStr, err := json.Marshal(input)
  if err != nil {
    log.Println(err)
		return nil,common.ResultJsonEncodeError
  }
	output:=&flowReqRsp{}
  if err := json.Unmarshal(jsonStr, output); err != nil {
    log.Println(err)
		return nil,common.ResultJsonDecodeError
  }

	return output,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorFlowAsync)getMqttClient(instanceID string)(*mqtt.Client){
	broker := nodeExecutor.Mqtt.Broker //"121.36.192.249"
	port := nodeExecutor.Mqtt.Port //1983
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d",broker,port))
	opts.SetClientID("flow_node_flowAsync_"+instanceID)
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

func (nodeExecutor *nodeExecutorFlowAsync)getRequestContent(
	instance *flowInstance,
	conf *nodeFlowAsyncConf,
	input *flowReqRsp)(string,int){
		output,commonErr:=nodeExecutor.copyFlowReqRsp(input)
		if commonErr!=common.ResultSuccess {
			return "",commonErr
		}

		output.TaskID=instance.TaskID
		output.TaskStep=instance.TaskStep+1
		output.FlowID=conf.FlowID
		output.DebugID=instance.DebugID
		output.FlowInstanceID=nil

		jsonStr, err := json.Marshal(output)
  	if err != nil {
    	log.Println(err)
			return "",common.ResultJsonEncodeError
  	}

		return string(jsonStr),common.ResultSuccess
}

func (nodeExecutor *nodeExecutorFlowAsync)sendRequest(reqStr string,instance *flowInstance){
	client:=nodeExecutor.getMqttClient(instance.InstanceID)
	if client == nil {
		return
	}
	defer (*client).Disconnect(250)

	token:=(*client).Publish(nodeExecutor.Mqtt.StartFlowTopic,0,false,reqStr)
	token.Wait()
}

func (nodeExecutor *nodeExecutorFlowAsync)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_FLOW_ASYNC,
	}

	req:=node.Input

	//加载节点配置
	nodeConf:=nodeExecutor.getNodeConf()
	if nodeConf==nil {
		log.Printf("nodeExecutorFlow run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeFilterConfigError,params)
	}

	//组装新的请求对象
	reqContent,commonErr:=nodeExecutor.getRequestContent(instance,nodeConf,req)
	if commonErr != common.ResultSuccess {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=req
		return req,common.CreateError(commonErr,params)
	}

	//发送请求到MQTT
	nodeExecutor.sendRequest(reqContent,instance)
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=req
	return req,nil
}