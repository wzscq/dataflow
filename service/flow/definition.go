package flow

import (
	"log"
	"os"
	"dataflow/common"
	"encoding/json"
	"github.com/rs/xid"
    "time"
	"io/ioutil"
	"strings"
)

type node struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Data interface{} `json:"data"`
	Position interface{} `json:"position"`
}

type edge struct {
	ID string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type flowConf struct {
	 Label interface{} `json:"label"`
	 Description interface{} `json:"description"`
	 Nodes []node `json:"nodes"`
	 Edges []edge `json:"edges"`
}

func deleteFlow(appDB,flowID string)(int){
	flowCfgFile:="apps/"+appDB+"/flows/"+flowID+".json"
	err := os.Remove(flowCfgFile)
    if err != nil {
        log.Println(err)
		return common.ResultDeleteFlowFileError
    }
	return common.ResultSuccess
}

func addFlowConf(appDB,flowID string,flowConf *flowConf)(int){
	flowCfgFile:="apps/"+appDB+"/flows/"+flowID+".json"
	if _, err := os.Stat(flowCfgFile); err == nil {
		return common.ResultFlowIDAlreadyExist
	}
	
	return saveFlowConf(appDB,flowID,flowConf)
}

func saveFlowConf(appDB,flowID string,flowConf *flowConf)(int){
	jsonStr, err := json.MarshalIndent(flowConf, "", "    ")
    if err != nil {
        log.Println(err)
    } else {
		flowCfgFile:="apps/"+appDB+"/flows/"+flowID+".json"
		ioutil.WriteFile(flowCfgFile, jsonStr, 0644)
	}
	return common.ResultSuccess
}

func loadFlowConf(appDB,flowID string)(*flowConf,int){
	//load flow config from file
	flowCfgFile:="apps/"+appDB+"/flows/"+flowID+".json"
	filePtr, err := os.Open(flowCfgFile)
	if err != nil {
		log.Printf("Open flow configuration file failed [Err:%s] \n", err.Error())
		return nil,common.ResultOpenFileError
	}
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	flowConf:=&flowConf{}
	err = decoder.Decode(flowConf)
	if err != nil {
		log.Printf("json file decode failed [Err:%s] \n", err.Error())
		return nil,common.ResultJsonDecodeError
	}
	return flowConf,common.ResultSuccess
}

func getInstanceID(appDB,flowID string)(string){
	guid := xid.New().String()
	nowStr:= time.Now().Format("20060102150405")
	return appDB+"_"+flowID+"_"+nowStr+"_"+guid
}

func logInstance(instance *flowInstance){
	//打印一下流的实例内容
	jsonStr, err := json.MarshalIndent(instance, "", "    ")
    if err != nil {
        log.Println(err)
    } else {
		fileName:="apps/"+instance.AppDB+"/instances/"+instance.InstanceID+".json"
		ioutil.WriteFile(fileName, jsonStr, 0644)
	}
}

func logInstanceNode(instance *flowInstance,node *instanceNode){
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
}

func createInstance(appDB,flowID,userID string,debugID *string,flowCfg *flowConf)(*flowInstance,int){
	//允许前端直接将配置传递到后端运行，如果没有给则根据flowID从文件加载
	if flowCfg==nil {
		var errorCode int
		flowCfg,errorCode=loadFlowConf(appDB,flowID)	
		if(errorCode!=common.ResultSuccess){
			return nil,errorCode
		}
	}

	instanceID:=getInstanceID(appDB,flowID)

	instance:=&flowInstance{
		AppDB:appDB,
	 	FlowID:flowID,
	 	InstanceID:instanceID,
	 	UserID:userID,
	 	FlowConf:flowCfg,
		Completed:false,
		DebugID:debugID,
		StartTime:time.Now().Format("2006-01-02 15:04:05"),
	}
	
	return instance,common.ResultSuccess
}

func getAppFlows(appDB string)(*[]string,int){
	//flow path
	flowConfPath:="apps/"+appDB+"/flows/"
	flowConfDir, err := os.Open(flowConfPath)
	if err != nil {
		log.Printf("Open flow configuration folder failed [Err:%s] \n", err.Error())
		return nil,common.ResultOpenFileError
	}
	defer flowConfDir.Close()

	flowConfFiles, err := flowConfDir.Readdir(0)
	if err != nil {
		log.Printf("Read flow configuration folder failed [Err:%s] \n", err.Error())
		return nil,common.ResultOpenFileError
	}
    // Loop over files.
	if len(flowConfFiles) <=0 {
		return nil,common.ResultSuccess
	} 

	appFlows:=make([]string,len(flowConfFiles))
    for index := range(flowConfFiles) {
        flowConfFile := flowConfFiles[index]
		flowID:=flowConfFile.Name()
		if strings.Contains(flowID, ".json") {
			flowID=flowID[:len(flowID)-5]
		}
		appFlows[index]=flowID
    }
	return &appFlows,common.ResultSuccess
}