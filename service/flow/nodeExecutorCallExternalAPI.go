package flow

import (
  "time"
	"dataflow/common"
	"dataflow/data"
	"log"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"bytes"
)

const (
	//CallExternalAPIReqForEachRowYes 每行数据都调用一次API
	CallExternalAPIReqForEachRowYes="yes"
	//CallExternalAPIReqForEachRowNo 只调用一次API，传入所有数据
	CallExternalAPIReqForEachRowNo="no"
)

type APIResult struct {
	Error bool `json:"error"`
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Result modelDataItem `json:"result"`
}

//API目前仅支持POST调用，仅支持一个模型的数据，采用统一的CRV平台接口
//返回的模型数据可以和请求模型一样，或者是一个新的模型，在返回数据中给出模型ID
//如果请求的模型和返回模型一样，那么需要API接口逻辑中自己处理原有数据的更新
type CallExternalAPIConfig struct {
	URL string `json:"url"`
	ModelID string `json:"modelID"`
	ReqForEachRow string `json:"reqForEachRow"`
}

type nodeExecutorCallExternalAPI struct {
	NodeConf node
	DataRepository data.DataRepository
}

func (nodeExecutor *nodeExecutorCallExternalAPI)getNodeConf()(*CallExternalAPIConfig){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok {
		return nil
	}
	jsonStr, err := json.Marshal(mapData)
  if err != nil {
    log.Println(err)
		return nil
  }
	conf:=CallExternalAPIConfig{}
  if err := json.Unmarshal(jsonStr, &conf); err != nil {
    log.Println(err)
		return nil
  }

	return &conf
}

func (nodeExecutor *nodeExecutorCallExternalAPI)getRequestDataList(modelData *modelDataItem,conf *CallExternalAPIConfig)(*[]modelDataItem){
	requestDataList:=[]modelDataItem{}
	//如果是需要基于行调用API，那么先将数据行进行拆分
	if conf.ReqForEachRow == CallExternalAPIReqForEachRowYes {
		for _,row:=range *modelData.List {
			requestData:=modelDataItem{
				ModelID:modelData.ModelID,
				List:&[]map[string]interface{}{row},
				Total:1,
				ViewID:modelData.ViewID,
				Filter:modelData.Filter,
				Fields:modelData.Fields,
				Sorter:modelData.Sorter,
			}
			requestDataList=append(requestDataList,requestData)
		}
	} else {
		requestDataList=append(requestDataList,*modelData)
	}

	return &requestDataList
}

func (nodeExecutor *nodeExecutorCallExternalAPI)mergeResult(resultList *[]modelDataItem)(*modelDataItem){
	//合并结果
	result:=modelDataItem{}
	for index,item:=range *resultList {
		if index==0 {
			result.ModelID=item.ModelID
		}
		(*result.List)=append(*result.List,(*item.List)...)
		result.Total+=item.Total
	}
	return &result
}



func (nodeExecutor *nodeExecutorCallExternalAPI)callAPI(
	requestData modelDataItem,
	conf *CallExternalAPIConfig,
	instance *flowInstance)(*modelDataItem,*common.CommonError){
	//通过http post 请求调用外部api
	//请求数据
	requestDataJson,err:=json.Marshal(requestData)
	if err!=nil {
		return nil,common.CreateError(common.ResultJsonMarshalError,nil)
	}
	//创建HTTP请求
	req, err := http.NewRequest("POST", conf.URL, bytes.NewBuffer(requestDataJson))
	if err != nil {
		log.Println(err)
		return nil,common.CreateError(common.ResultCallExternalAPIError,nil)
	}
	//设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("UserID", instance.UserID)
	req.Header.Set("AppDB", instance.AppDB)
	//发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil,common.CreateError(common.ResultCallExternalAPIError,nil)
	}
	defer resp.Body.Close()
	//读取返回数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil,common.CreateError(common.ResultCallExternalAPIError,nil)
	}
	//解析返回数据
	apiResult:=APIResult{}
	err = json.Unmarshal(body, &apiResult)
	if err != nil {
		log.Println(err)
		return nil,common.CreateError(common.ResultCallExternalAPIError,nil)
	}
	if apiResult.Error {
		log.Println(apiResult.Message)
		params:=map[string]interface{}{
			"errorCode":apiResult.ErrorCode,
			"message":apiResult.Message,
		}
		return nil,common.CreateError(common.ResultCallExternalAPIError,params)
	}

	return &apiResult.Result,nil
}

func (nodeExecutor *nodeExecutorCallExternalAPI)dealModelData(
	modelData *modelDataItem,
	conf *CallExternalAPIConfig,
	instance *flowInstance)(*modelDataItem,*common.CommonError){
	requestDataList:=nodeExecutor.getRequestDataList(modelData,conf)
	resultList := []modelDataItem{}
	//调用API
	for _,requestData:=range *requestDataList {
		//调用API
		result,err:=nodeExecutor.callAPI(requestData,conf,instance)
		if err!=nil {
			return nil,err
		}
		resultList=append(resultList,*result)
	}
	//合并结果
	result:=nodeExecutor.mergeResult(&resultList)
	return result,nil
}

func (nodeExecutor *nodeExecutorCallExternalAPI)dealItem(
	item *flowDataItem,
	conf *CallExternalAPIConfig,
	instance *flowInstance)(*common.CommonError){
	for index,modelData:=range item.Models {
		if *modelData.ModelID==conf.ModelID {
			//调用API
			result,err:=nodeExecutor.dealModelData(&modelData,conf,instance)
			if err!=nil {
				return err
			}
			//如果返回的模型和请求模型一样，那么需要更新原有数据
			if result.ModelID==modelData.ModelID {
				//更新原有数据
				item.Models[index]=*result
			} else {
				//否则，将返回的模型数据添加到数据项中
				item.Models=append(item.Models,*result)
			}
		}
	}
	return nil
}

func (nodeExecutor *nodeExecutorCallExternalAPI)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	if node.Input.Data==nil || len(*node.Input.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=node.Input
		return node.Input,nil
	}

	req:=node.Input
	flowResult:=&flowReqRsp{
		FlowID:req.FlowID,
		FlowInstanceID:req.FlowInstanceID,
		Stage:req.Stage,
		DebugID:req.DebugID,
		UserRoles:req.UserRoles,
		UserID:req.UserID,
		AppDB:req.AppDB,
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
		GoOn:true,
	}

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_CALL_EXTERNAL_API,
	}

	//获取节点配置
	conf:=nodeExecutor.getNodeConf()
	if conf==nil {
		log.Printf("nodeExecutorCallExternalAPI run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeConfigError,params)
	}

	resultData:=make([]flowDataItem,len(*node.Input.Data))

	//对每个数据项进行调用
	for index,item:= range (*req.Data) {
		//处理每个数据项
		err:=nodeExecutor.dealItem(&item,conf,instance)
		if err!=nil {
			if err.Params==nil {
				err.Params=params
			} else {
				err.Params["nodeID"]=node.ID
				err.Params["nodeType"]=NODE_CALL_EXTERNAL_API
			}
			return node.Input,err
		}
		resultData[index]=item
	}

	flowResult.Data=&resultData
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	return flowResult,nil
}