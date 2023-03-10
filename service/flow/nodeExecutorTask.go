package flow

import (
  "time"
	"dataflow/common"
	"dataflow/data"
	"dataflow/user"
	"encoding/json"
	"log"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type TaskConf struct {
	Name *string `json:"name,omitempty"`
	ExecutionStatus *string `json:"executionStatus,omitempty"`
	ExecutionProgress string `json:"executionProgress"`
	ResultStatus *string `json:"resultStatus,omitempty"`
	ErrorCode *string `json:"errorCode,omitempty"`
	Message *string `json:"errorCode,message"`
}

type StepConf struct {
	Name *string `json:"name,omitempty"`
	ExecutionStatus *string `json:"executionStatus,omitempty"`
	ExecutionProgress string `json:"executionProgress"`
	ResultStatus *string `json:"resultStatus,omitempty"`
	ErrorCode *string `json:"errorCode,omitempty"`
	Message *string `json:"errorCode,message"`
}

type nodeExecutorTaskConf struct {
	Task TaskConf `json:"task"`
	Step StepConf `json:"step"`
}

type nodeExecutorTask struct {
	NodeConf node
	Mqtt *common.MqttConf
	Redis *common.RedisConf
	DataRepository data.DataRepository
}

func (nodeExecutor *nodeExecutorTask) connectHandler(client mqtt.Client){
	log.Println("connectHandler connect status: ",client.IsConnected())
}

func (nodeExecutor *nodeExecutorTask) connectLostHandler(client mqtt.Client, err error){
	log.Println("connectLostHandler connect status: ",client.IsConnected(),err)
}

func (nodeExecutor *nodeExecutorTask) messagePublishHandler(client mqtt.Client, msg mqtt.Message){
	log.Println("messagePublishHandler topic: ",msg.Topic())
}

func (nodeExecutor *nodeExecutorTask)getMqttClient(instanceID string)(*mqtt.Client){
	broker := nodeExecutor.Mqtt.Broker //"121.36.192.249"
	port := nodeExecutor.Mqtt.Port //1983
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d",broker,port))
	opts.SetClientID("flow_node_task_"+instanceID)
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

func (nodeExecutor *nodeExecutorTask)getNodeConf()(*nodeExecutorTaskConf){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok {
		return nil
	}
	jsonStr, err := json.Marshal(mapData)
  if err != nil {
    log.Println(err)
		return nil
  }
	conf:=nodeExecutorTaskConf{}
  if err := json.Unmarshal(jsonStr, &conf); err != nil {
    log.Println(err)
		return nil
  }

	return &conf
}

func (nodeExecutor *nodeExecutorTask)getTask(instance *flowInstance)(*[]map[string]interface{},int){
	fieldType:="one2many"
	relatedModelID:="core_task_step"
	relatedField:="task_id"

	query:=&data.Query{
		ModelID:"core_task",
		Filter:&map[string]interface{}{
			"id":*instance.TaskID,
		},
		Fields:&[]data.Field{
			data.Field{
				Field:"id",
			},
			data.Field{
				Field:"version",
			},
			data.Field{
				Field:"steps",
				FieldType:&fieldType,
				RelatedModelID:&relatedModelID,
				RelatedField:&relatedField,
				Filter:&map[string]interface{}{
					"step":instance.TaskStep,
				},
				Fields:&[]data.Field{
					data.Field{
						Field:"id",
					},
					data.Field{
						Field:"version",
					},
				},
			},
		},
		AppDB:instance.AppDB,
	}
	result,errorCode:=query.Execute(nodeExecutor.DataRepository,false)
	return &result.List,errorCode
}

func (nodeExecutor *nodeExecutorTask)getTaskStepVersion(taskInfo *[]map[string]interface{})(int64,int64){
	var taskVersion int64 =-1
	var stepVersion int64 =-1
	if taskInfo != nil && len(*taskInfo)>0 {
		taskVersion,_=(*taskInfo)[0]["version"].(int64)
		steps,ok:=(*taskInfo)[0]["steps"]
		if ok {
			stepMap,ok:=steps.(map[string]interface{})
			if ok {
				stepList,ok:=stepMap["list"].([]interface{})
				if ok && len(stepList)>0 {
					stepRow,ok:=stepList[0].(map[string]interface{})
					if ok {
						stepVersion,_=stepRow["version"].(int64)
					}
				}
			}
		}
	}
	return taskVersion,stepVersion
}

func (nodeExecutor *nodeExecutorTask)confToTaskInfo(
	nodeConf *nodeExecutorTaskConf,
	taskVersion int64,
	stepInfo *map[string]interface{},
	taskID string)(*map[string]interface{}){
		
		taskInfo:=map[string]interface{}{
			"id":taskID,
			"execution_progress":nodeConf.Task.ExecutionProgress,
			"steps":map[string]interface{}{
				"fieldType":"one2many",
				"modelID":"core_task_step",
				"relatedField":"task_id",
				"list":[]map[string]interface{}{
					*stepInfo,
				},
			},
		}

		if nodeConf.Task.Name!=nil && len(*nodeConf.Task.Name)>0 {
			taskInfo["name"]=*nodeConf.Task.Name
		}

		if nodeConf.Task.ExecutionStatus!=nil && len(*nodeConf.Task.ExecutionStatus)>0 {
			taskInfo["execution_status"]=*nodeConf.Task.ExecutionStatus
		}
		
		if nodeConf.Task.ResultStatus!=nil && len(*nodeConf.Task.ResultStatus)>0 {
			taskInfo["result_status"]=*nodeConf.Task.ResultStatus
		}

		if nodeConf.Task.ErrorCode!=nil && len(*nodeConf.Task.ErrorCode)>0 {
			taskInfo["error_code"]=*nodeConf.Task.ErrorCode
		}

		if nodeConf.Task.Message!=nil && len(*nodeConf.Task.Message)>0 {
			taskInfo["message"]=*nodeConf.Task.Message
		}

		if taskVersion>=0 {
			taskInfo[data.SAVE_TYPE_COLUMN]=data.SAVE_UPDATE
			taskInfo[data.CC_VERSION]=taskVersion
		} else {
			taskInfo[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE
		}

		return &taskInfo
}

func (nodeExecutor *nodeExecutorTask)confToTaskStepInfo(
	nodeConf *nodeExecutorTaskConf,
	stepVersion int64,
	taskID string,
	taskStep int)(*map[string]interface{}){
		
		stepInfo:=map[string]interface{}{
			"id":fmt.Sprintf("%s_%d",taskID,taskStep),
			"execution_progress":nodeConf.Step.ExecutionProgress,
		}

		if nodeConf.Step.Name!=nil && len(*nodeConf.Step.Name)>0 {
			stepInfo["name"]=*nodeConf.Step.Name
		}

		if nodeConf.Step.ExecutionStatus!=nil && len(*nodeConf.Step.ExecutionStatus)>0 {
			stepInfo["execution_status"]=*nodeConf.Step.ExecutionStatus
		}
		
		if nodeConf.Step.ResultStatus!=nil && len(*nodeConf.Step.ResultStatus)>0 {
			stepInfo["result_status"]=*nodeConf.Step.ResultStatus
		}

		if nodeConf.Step.ErrorCode!=nil && len(*nodeConf.Step.ErrorCode)>0 {
			stepInfo["error_code"]=*nodeConf.Step.ErrorCode
		}

		if nodeConf.Step.Message!=nil && len(*nodeConf.Step.Message)>0 {
			stepInfo["message"]=*nodeConf.Step.Message
		}

		if stepVersion>=0 {
			stepInfo[data.SAVE_TYPE_COLUMN]=data.SAVE_UPDATE
			stepInfo[data.CC_VERSION]=stepVersion
		} else {
			stepInfo[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE
			stepInfo["task_id"]=taskID
		}

		return &stepInfo
}

func (nodeExecutor *nodeExecutorTask)getSaveObject(
	instance *flowInstance,
	nodeConf *nodeExecutorTaskConf,
	taskVersion,stepVersion int64)(*data.Save){
		taskStepInfo:=nodeExecutor.confToTaskStepInfo(nodeConf,stepVersion,*instance.TaskID,instance.TaskStep)
		taskInfo:=nodeExecutor.confToTaskInfo(nodeConf,taskVersion,taskStepInfo,*instance.TaskID)

		saver:=&data.Save{
			List:&[]map[string]interface{}{
				*taskInfo,
			},
			AppDB:instance.AppDB,
			UserID:instance.UserID,
			ModelID:"core_task",
		}
		return saver
}

func (nodeExecutor *nodeExecutorTask)saveTask(
	saver *data.Save)(int){
	
	log.Println("start nodeExecutorTask saveTask")
	
	//开启事务
	tx,err:= nodeExecutor.DataRepository.Begin()
	if err != nil {
		log.Println(err)
		return common.ResultSQLError
	}

	_,errorCode:=saver.SaveList(nodeExecutor.DataRepository,tx)
	if errorCode!=common.ResultSuccess {
		tx.Rollback()
		log.Println("end nodeExecutorTask saveTask with error")
		return errorCode
	}
	
	//提交事务
	if err := tx.Commit(); err != nil {
		log.Println(err)
		log.Println("end nodeExecutorTask saveTask with error")
		return common.ResultSQLError
	}
	log.Println("end nodeExecutorTask saveTask")
	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorTask)sendNotification(
	instance *flowInstance,
	saver *data.Save,
	userToken string){

	jsonStr, err := json.Marshal(saver)
	if err != nil {
		log.Println(err)
		return
	}

	client:=nodeExecutor.getMqttClient(instance.InstanceID)
	if client == nil {
		return
	}
	defer (*client).Disconnect(250)

	topic:=fmt.Sprintf("%s/%s",nodeExecutor.Mqtt.TaskNotificationTopic,userToken)
	token:=(*client).Publish(topic,0,false,string(jsonStr))
	token.Wait()
}

func (nodeExecutor *nodeExecutorTask)getUserToken(userID string)(string){
	//获取用户token
	duration, _ := time.ParseDuration(nodeExecutor.Redis.TokenExpired)
  loginCache:=&user.DefatultLoginCache{}
  loginCache.Init(nodeExecutor.Redis.Server,nodeExecutor.Redis.TokenDB,duration,nodeExecutor.Redis.Password)
	userToken, err := loginCache.GetUserToken(userID)
	if err != nil {
		log.Println(err)
		return ""
	}
	return userToken
}

func (nodeExecutor *nodeExecutorTask)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_TASK,
	}

	//加载节点配置
	nodeConf:=nodeExecutor.getNodeConf()
	if nodeConf==nil {
		log.Printf("nodeExecutorTask run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeFilterConfigError,params)
	}

	//获取任务信息
	taskInfo,errorCode:=nodeExecutor.getTask(instance)
	if errorCode != common.ResultSuccess {
		return node.Input,common.CreateError(errorCode,params)
	}
	//从任务信息中读取任务和步骤的版本
	taskVersion,stepVersion:=nodeExecutor.getTaskStepVersion(taskInfo)

	//根据配置信息创建或更新任务表信息
	//这里注意任务ID，和步骤ID都是保存在instance参数中的
	saver:=nodeExecutor.getSaveObject(instance,nodeConf,taskVersion,stepVersion)
	errorCode=nodeExecutor.saveTask(saver)
	if errorCode != common.ResultSuccess {
		return node.Input,common.CreateError(errorCode,params)
	}
	//通过MQTT发送通知消息给当前用户，后续需要扩展到发送给任务的相关用户，消息topic中携带用户的token
	//首先获取当前用户的token信息
	userToken:=nodeExecutor.getUserToken(instance.UserID)
	if len(userToken)>0 {
		nodeExecutor.sendNotification(instance,saver,userToken)
	}

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=node.Input

	return node.Output,nil
}