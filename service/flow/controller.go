package flow

import (
	"log"
	"github.com/gin-gonic/gin"
	"dataflow/common"
	"dataflow/data"
	"net/http"
	"encoding/json"
)

type modelDataItem struct {
	ModelID *string `json:"modelID,omitempty"`
	ViewID *string `json:"viewID,omitempty"`
	Filter *map[string]interface{} `json:"filter,omitempty"`
	List *[]map[string]interface{} `json:"list,omitempty"`
	Fields *interface{} `json:"fields,omitempty"`
	Sorter *interface{} `json:"sorter,omitempty"`
	Total int `json:"total"`
}

type verifyResultItem struct {
	VerfiyID string `json:"verfiyID,omitempty"`
	VerfiyType string `json:"verfiyType,omitempty"`
	Message string `json:"message,omitempty"`
	Result string `json:"result,omitempty"`
}

type flowDataItem struct {
	VerifyResult []verifyResultItem `json:"verifyResult,omitempty"`
	Models []modelDataItem `json:"models,omitempty"`
}

type CommonHeader struct {
	Token     string  `json:"token"`
	UserID    string  `json:"userID"`
	AppDB     string  `json:"appDB"`
	UserRoles string  `json:"userRoles"`
}

type flowReqRsp struct {
	FlowID string `json:"flowID"`
	FlowInstanceID *string `json:"flowInstanceID,omitempty"`
	//增加任务调度需要的相关属性
	TaskID *string `json:"taskID,omitempty"`  //taskID属性标识任务，如果请求中没有提供taskID则默认和flowInstanceID一致
	TaskStep int `json:"taskStep"`            //taskStep标识任务的步骤，默认0
	Stage *int `json:"stage,omitempty"`
	DebugID *string `json:"debugID,omitempty"`
	UserRoles string  `json:"userRoles"`
	UserID    string  `json:"userID"`
	AppDB     string  `json:"appDB"`
	FlowConf *flowConf `json:"flowConf,omitempty"`
	ModelID *string `json:"modelID"`
	ViewID *string `json:"viewID"`
	FilterData *[]data.FilterDataItem `json:"filterData"`
	Filter *map[string]interface{} `json:"filter"`
	List *[]map[string]interface{} `json:"list"`
	Total int `json:"total"`
	GoOn bool   //返回当前节点是否需要继续执行后续节点，默认true，继续执行
	Over bool   //返回是否终止流的执行，默认false
	//Fields *[]field `json:"fields"`
	//Sorter *[]sorter `json:"sorter"`
	SelectedRowKeys *[]string `json:"selectedRowKeys"`
	Pagination *data.Pagination `json:"pagination"`
	Operation *map[string]interface{} `json:"operation,omitempty"`
	/**
	data结构示例：
	[
		{
			"verifyResult":[
				{
					"verifyID":"",
					"modelID":"",
					"verifyType":"",
					"message":""
				}
			],
			"models":[{
				ModelID:"modelid1",
				ViewID:"view1",
				Filter:{...},
				List:[{...},{...}],
				Fields:{},
				Sorter:{},
			},
			{
				ModelID:"modelid2",
				ViewID:"view2",
				Filter:{...},
				List:[{...},{...}],
				Fields:{},
				Sorter:{},
			}]
		}
	]
	**/
	Data *[]flowDataItem `json:"data,omitempty"`	
}

type FlowController struct {
	DataRepository data.DataRepository
	InstanceRepository FlowInstanceRepository
	Mqtt common.MqttConf
}

//start flow by mqtt
func (controller *FlowController)StartFlow(reqPayload []byte){
	log.Println("FlowController start StartFlow")
	var req flowReqRsp
	if err := json.Unmarshal(reqPayload, &req); err != nil {
		log.Println(err)
		log.Println("FlowController end StartFlow with error")
		return
	}

	req.GoOn=true  //这个值设置节点是否继续运行默认为true
	req.Over=false 

	//创建流
	flowInstance,errorCode:=createInstance(
		req.AppDB,
		req.FlowID,
		req.UserID,
		req.FlowInstanceID,
		req.DebugID,
		req.TaskID,
		req.TaskStep,
		req.FlowConf,
		controller.InstanceRepository)
	if errorCode!=common.ResultSuccess {
		log.Printf("create flow instance error with no %d",errorCode)
		log.Println("FlowController end StartFlow with error")
		return
	}
	//执行流
	flowInstance.push(controller.DataRepository,&req,&controller.Mqtt)
	log.Println("FlowController end StartFlow")
}

func (controller *FlowController)start(c *gin.Context){
	log.Println("FlowController start start")
	var header CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	//加载一个流的配置
	var req flowReqRsp
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("FlowController start end")
		return
	}

	req.UserID=header.UserID
	req.AppDB=header.AppDB
	req.UserRoles=header.UserRoles
	req.GoOn=true  //这个值设置节点是否继续运行默认为true
	req.Over=false 

	//创建流
	flowInstance,errorCode:=createInstance(
		req.AppDB,
		req.FlowID,
		req.UserID,
		req.FlowInstanceID,
		req.DebugID,
		req.TaskID,
		req.TaskStep,
		req.FlowConf,
		controller.InstanceRepository)
	if errorCode!=common.ResultSuccess {
		rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("FlowController start end")
		return
	}
	//执行流
	result,err:=flowInstance.push(controller.DataRepository,&req,&controller.Mqtt)

	rsp:=common.CreateResponse(err,result)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("FlowController start end")
}

func (controller *FlowController)list(c *gin.Context){
	log.Println("FlowController list start")

	var header CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}
	
	var req flowReqRsp
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("FlowController list end")
		return
	} 

	//读取所有流定义文件的名字，名字就是ID
	result,errorCode:=getAppFlows(header.AppDB)
	rsp:=common.CreateResponse(common.CreateError(errorCode,nil),result)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("FlowController start end")
}

func (controller *FlowController)getConfig(c *gin.Context){
	log.Println("FlowController getConfig start")

	var header CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//加载一个流的配置
	var req flowReqRsp
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("FlowController getConfig end")
		return
	} 

	//获取流配置
	flowConf,errorCode:=loadFlowConf(header.AppDB,req.FlowID)
	//logInstance(flowInstance)
	rsp:=common.CreateResponse(common.CreateError(errorCode,nil),flowConf)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("FlowController getConfig end")
}

func (controller *FlowController)saveConfig(c *gin.Context){
	log.Println("FlowController saveConfig start")

	var header CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//加载一个流的配置
	var req flowReqRsp
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("FlowController saveConfig end")
		return
	} 

	//保存流配置
	errorCode:=saveFlowConf(header.AppDB,req.FlowID,req.FlowConf)
	//logInstance(flowInstance)
	rsp:=common.CreateResponse(common.CreateError(errorCode,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("FlowController saveConfig end")
}

func (controller *FlowController)addFlow(c *gin.Context){
	log.Println("FlowController addFlow start")

	var header CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//加载一个流的配置
	var req flowReqRsp
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("FlowController addFlow end")
		return
	} 

	//保存流配置
	errorCode:=addFlowConf(header.AppDB,req.FlowID,req.FlowConf)

	rsp:=common.CreateResponse(common.CreateError(errorCode,nil),req.FlowID)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("FlowController addFlow end")
}

func (controller *FlowController)deleteFlow(c *gin.Context){
	log.Println("FlowController deleteFlow start")

	var header CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	//加载一个流的配置
	var req flowReqRsp
	if err := c.BindJSON(&req); err != nil {
		log.Println(err)
		rsp:=common.CreateResponse(common.CreateError(common.ResultWrongRequest,nil),nil)
		c.IndentedJSON(http.StatusOK, rsp)
		log.Println("FlowController deleteFlow end")
		return
	} 

	//删除流配置
	errorCode:=deleteFlow(header.AppDB,req.FlowID)

	rsp:=common.CreateResponse(common.CreateError(errorCode,nil),req.FlowID)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("FlowController deleteFlow end")
}

func (controller *FlowController)getMqttServer(c *gin.Context){
	log.Println("FlowController getMqttServer start")
	//加载一个流的配置
	rsp:=common.CreateResponse(nil,controller.Mqtt)
	c.IndentedJSON(http.StatusOK, rsp)
	log.Println("FlowController getMqttServer end")
}

func (controller *FlowController) Bind(router *gin.Engine) {
	log.Println("Bind FlowController")
	router.POST("/flow/start", controller.start)
	router.POST("/flow/list",controller.list)
	router.POST("/flow/getConfig",controller.getConfig)
	router.POST("/flow/saveConfig",controller.saveConfig)
	router.POST("/flow/addFlow",controller.addFlow)
	router.POST("/flow/deleteFlow",controller.deleteFlow)
	router.POST("/flow/getMqttServer",controller.getMqttServer)
}