package flow

import (
    "time"
	"dataflow/common"
	"dataflow/data"
	"encoding/json"
	"log"
	"github.com/rs/xid"
)

const (
	KEY_FIELD_YES = "1"
)

type dataTransferField struct {
	KeyField string `json:"keyField"`
	SourceField string `json:"sourceField"`
	TargetField string `json:"targetField"`
}

type updateField struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type sourceModel struct {
	ModelID string `json:"modelID"`
	BatchNumberField string `json:"batchNumberField"`
	Fields []dataTransferField `json:"fields"`
	UpdateFields []updateField `json:"updateFields"`
}

type dataTransferConf struct {
	TargetModelID string `json:"targetModelID"`
	BatchNumberField string `json:"batchNumberField"`
	SourceModels []sourceModel  `json:"sourceModels"`
}

type nodeExecutorDataTransfer struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorDataTransfer)getNodeConf()(*dataTransferConf){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	log.Println(string(jsonStr))
	conf:=&dataTransferConf{}
    if err := json.Unmarshal(jsonStr, conf); err != nil {
        log.Println(err)
		return nil
    }

	return conf
}

func (nodeExecutor *nodeExecutorDataTransfer)getTargetRow(
	sourceModelCfg *sourceModel,
	sourceDataRow *map[string]interface{},
	batchNumber string,
	targetBatchNumberField string)(*map[string]interface{}){
	
	targetRow:=map[string]interface{}{}
	//循环每个字段,按照原字段名取值，并复制到目标字段
	for _,fieldCfg:=range(sourceModelCfg.Fields) {
		srcVal,ok:=(*sourceDataRow)[fieldCfg.SourceField]
		if ok {
			targetRow[fieldCfg.TargetField]=srcVal
		}
	}

	if len(targetBatchNumberField)>0 {
		targetRow[targetBatchNumberField]=batchNumber
	}
	targetRow[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE

	return &targetRow
}

func (nodeExecutor *nodeExecutorDataTransfer)getUpdateRow(
	sourceModelCfg *sourceModel,
	sourceDataRow *map[string]interface{},
	batchNumber string)(*map[string]interface{}){
	
	updateRow:=map[string]interface{}{}
	//循环每个字段,按照原字段名取值，并复制到目标字段
	for _,updateFieldCfg:=range(sourceModelCfg.UpdateFields) {
		updateRow[updateFieldCfg.Field]=updateFieldCfg.Value
	}

	if len(sourceModelCfg.BatchNumberField)>0 {
		updateRow[sourceModelCfg.BatchNumberField]=batchNumber
	}

	updateRow[data.CC_VERSION]=(*sourceDataRow)[data.CC_VERSION]
	updateRow[data.CC_ID]=(*sourceDataRow)[data.CC_ID]
	updateRow[data.SAVE_TYPE_COLUMN]=data.SAVE_UPDATE

	return &updateRow
}

func (nodeExecutor *nodeExecutorDataTransfer)getTargetData(
	sourceModelCfg *sourceModel,
	sourceModelData *modelDataItem,
	targetData *[]map[string]interface{},
	updateData *[]map[string]interface{},
	batchNumber string,
	targetBatchNumberField string){
	
	//循环处理每一行数据
	for _,dataRow:=range(*sourceModelData.List) {
		targetRow:=nodeExecutor.getTargetRow(sourceModelCfg,&dataRow,batchNumber,targetBatchNumberField)
		(*targetData)=append((*targetData),*targetRow)
		if len(sourceModelCfg.UpdateFields)>0 {
			updateRow:=nodeExecutor.getUpdateRow(sourceModelCfg,&dataRow,batchNumber)
			(*updateData)=append((*updateData),*updateRow)
		}
	}
}

func (nodeExecutor *nodeExecutorDataTransfer)getBatchNumber()(string){
	return xid.New().String()
}

func (nodeExecutor *nodeExecutorDataTransfer)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorDataTransfer run start")
	
	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_DATA_TRANSFER,
	}
	//加载节点配置
	nodeCfg:=nodeExecutor.getNodeConf()
	if nodeCfg==nil {
		log.Printf("nodeExecutorDataTransfer run get node config error\n")
		return node.Input,common.CreateError(common.ResultNodeConfigError,params)
	}

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
		UserID:req.UserID,
		AppDB:req.AppDB,
	}

	//生成处理批次号
	batchNumber:=nodeExecutor.getBatchNumber()
	//将数据从源表字段中读取，并写入到目标表对应字段上，
	//如果有多个原表的数据，则写入时按照关键字段索引，将多个原表的数据写入到同一个目标表的数据行上
	//这里先遍历所有配置中的源表
	//数据先按照关键字段索引放入map中
	targetData:=[]map[string]interface{}{}
	updateData:=[]map[string]interface{}{}
	sourceModelID:=""
	for _,sourceModel:= range (nodeCfg.SourceModels) {
		//处理每一个flowDataItem
		for _,flowDataItem:= range (*req.Data) {
			//处理每个模型的数据
			for _,modelItem:= range (flowDataItem.Models) {
				//判断如果当前数据模型是一个数据源，则处理这个模型的数据
				if sourceModel.ModelID == *modelItem.ModelID {
					sourceModelID=sourceModel.ModelID
					nodeExecutor.getTargetData(
						&sourceModel,
						&modelItem,
						&targetData,
						&updateData,
						batchNumber,
						nodeCfg.BatchNumberField)
				}
			}
		}
	}

	models:=[]modelDataItem{}
	if len(targetData)>0 {
		targetModelItem:=modelDataItem{
			ModelID:&nodeCfg.TargetModelID,
			List:&targetData,
		}
		models=append(models,targetModelItem)
	}

	if len(updateData)>0 {
		updateModelItem:=modelDataItem{
			ModelID:&sourceModelID,
			List:&updateData,
		}
		models=append(models,updateModelItem)
	}

	flowData:=[]flowDataItem{
		flowDataItem{
			Models:models,
		},
	}

	flowResult.Data=&flowData

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorGroup run end")
	return flowResult,nil
}