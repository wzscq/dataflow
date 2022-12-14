package flow

import (
  "time"
	"dataflow/common"
	"dataflow/esi"
	"log"
	"encoding/json"
	"dataflow/data"
)

type nodeExecutorESI struct {
	NodeConf node
	DataRepository data.DataRepository 
}

func (nodeExecutor *nodeExecutorESI)getNodeConf()(*esi.ESIModel){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
	if err != nil {
		log.Println(err)
		return nil
	}
	log.Println(string(jsonStr))
	esiMOdel:=&esi.ESIModel{}
  if err := json.Unmarshal(jsonStr, esiMOdel); err != nil {
    log.Println(err)
		return nil
  }

	return esiMOdel
}

func (nodeExecutor *nodeExecutorESI)loadTestData()(*[]map[string]interface{}){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	testData,ok:=mapData["testData"]
	if !ok {
		return nil
	}

	testDataMap,ok:=testData.(map[string]interface{})
	if !ok {
		return nil
	}

	list:=[]map[string]interface{}{
		testDataMap,
	}
	return &list
}

func (nodeExecutor *nodeExecutorESI)getImportFile(inputRowData *map[string]interface{})(string,string,int){
	fileField,ok:=(*inputRowData)["esiFile"]
	if !ok {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("the field esiFile is not found.")
		return "","",common.ResultWrongRequest
	}

	fileValueMap,ok:=fileField.(map[string]interface{})
	if !ok {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("can not onvert esiFile value to map[stirng]interface{}.")
		return "","",common.ResultWrongRequest
	}

	listField,ok:=fileValueMap["list"]
	if !ok {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("the field esiFile.list is not found.")
		return "","",common.ResultWrongRequest
	}

	esiFileList,ok:=listField.([]interface{})
	if !ok {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("can not onvert esiFile.list to []interface{}.")
		return "","",common.ResultWrongRequest
	}

	if len(esiFileList)==0 {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("esiFile.list is empty.")
		return "","",common.ResultWrongRequest
	}

	esiFileRow,ok:=esiFileList[0].(map[string]interface{})
	if !ok {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("can not onvert esiFile.list[0] to map[stirng]interface{}.")
		return "","",common.ResultWrongRequest
	}

	//拿到文件名和文件内容
	fileNameIntreface,ok:=esiFileRow["name"]
	if !ok {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("the field esiFile.list[0].name is not found.")
		return "","",common.ResultWrongRequest
	}
	fileName,ok:=fileNameIntreface.(string)
	if !ok {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("can not onvert esiFile.list[0].name to string.")
		return "","",common.ResultWrongRequest
	}

	fileContentInterface,ok:=esiFileRow["contentBase64"]
	if !ok {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("the field esiFile.list[0].contentBase64 is not found.")
		return "","",common.ResultWrongRequest
	}
	fileContent,ok:=fileContentInterface.(string)
	if !ok {
		log.Println("nodeExecutorESI getImportFile end with error:")
		log.Println("can not onvert esiFile.list[0].contentBase64 to string.")
		return "","",common.ResultWrongRequest
	}

	return fileName,fileContent,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorESI)getRowData(
	modelID string,
	data *[]flowDataItem)(*map[string]interface{}){
		for _,dataItem:=range(*data){
			 for _,model:=range(dataItem.Models){
				if *model.ModelID==modelID {
					if model.List!=nil && len(*model.List)>0 {
						return &((*model.List)[0])
					} else {
						return nil
					}
				}
			 }
		}
		
		return nil
}

func (nodeExecutor *nodeExecutorESI)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){
	
		log.Println("nodeExecutorESI run start")
	
		req:=node.Input
		flowResult:=&flowReqRsp{
			FlowID:req.FlowID, 
			UserID:req.UserID,
			AppDB:req.AppDB,
			GoOn:true,
		}
	
		params:=map[string]interface{}{
			"nodeID":node.ID,
			"nodeType":NODE_ESI,
		}

		//如果是调试模式，择用测试数据填充List
		if instance.DebugID!=nil && len(*instance.DebugID)>0 {
			req.List=nodeExecutor.loadTestData()
		}

		//需要页面传入文件信息，ESI不再从输入参数中获取数据，而是从data中获取
		if req.Data==nil || len(*req.Data)==0 {
			log.Printf("nodeExecutorESI no input data\n")
			return flowResult,common.CreateError(common.ResultWrongRequest,params)
		}

		//加载节点配置
		esiModel:=nodeExecutor.getNodeConf()
		if esiModel==nil {
			log.Printf("nodeExecutorESI run get node config error\n")
			return flowResult,common.CreateError(common.ResultNodeConfigError,params)
		}

		inputRowData:=nodeExecutor.getRowData(esiModel.ModelID,req.Data)
		if inputRowData==nil {
			log.Printf("nodeExecutorESI no input data\n")
			return flowResult,common.CreateError(common.ResultWrongRequest,params)
		} 
		//这里考虑到导入操作
		//inputRowData:=(*req.List)[0]
		fileName,fileContent,errorCode:=nodeExecutor.getImportFile(inputRowData)
		if errorCode!=common.ResultSuccess {
			errorCode=common.ResultWrongRequest
			return flowResult,common.CreateError(errorCode,params)
		}

		//检查对应的文件名称如果已经导入过则不允许导入
		if esiModel.Options.PreventSameFile==true {
			errorCode=nodeExecutor.checkImportFile(req.AppDB,esiModel.ModelID,fileName)
			if errorCode!=common.ResultSuccess {
				return flowResult,common.CreateError(errorCode,params)
			}
		}

		log.Println(fileName)
		esiImport:=&esi.ESImport{
			AppDB:req.AppDB,
			ModelID:esiModel.ModelID,
			UserID:req.UserID,
			UserRoles:req.UserRoles,
			FileName:fileName,
			FileContent:fileContent,
			InputRowData:inputRowData,
		}
	
		result,commonErr:=esiImport.DoImport(esiModel)
		if commonErr!=nil {
			return flowResult,commonErr
		}

		list:=result["list"].([]map[string]interface{})
		total:=result["count"].(int)
		modelData:=modelDataItem{
			ModelID:&esiModel.ModelID,
			List:&list,
			Total:total,
		}

		flowData:=&[]flowDataItem{
			flowDataItem{
				Models:[]modelDataItem{
					modelData,
				},
			},
		}

		flowResult.Data=flowData
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=flowResult
		log.Println("nodeExecutorESI run end")
		return flowResult,nil
}

func (nodeExecutor *nodeExecutorESI)checkImportFile(appDB,modelID,fileName string)(int){
	query:=&data.Query{
		ModelID:modelID,
		Pagination: &data.Pagination{
			Current:1,
			PageSize:1,
		},
		Filter:&map[string]interface{}{
			esi.CC_IMPORT_FILE:fileName,
		},
		Fields:&[]data.Field{
			data.Field{
				Field:data.CC_ID,
			},
		},
		AppDB:appDB,
	}
	result,err:=query.Execute(nodeExecutor.DataRepository,false)
	if err!=common.ResultSuccess {
		return err
	}
	if result.Total>0 {
		return common.ResultESIFileAlreadyImported
	}
	return common.ResultSuccess
}