package flow

import (
    "time"
	"dataflow/common"
	"dataflow/data"
	"log"
	"encoding/json"
	"database/sql"
	"strconv"
	"fmt"
	"github.com/rs/xid"
)

type saveMatchedConfig struct {
	MatchResult string  `json:"matchResult"`  
	GroupModel groupModelConfig `json:"groupModel"`
}

type nodeExecutorSaveMatched struct {
	NodeConf node
	DataRepository data.DataRepository
}

func (nodeExecutor *nodeExecutorSaveMatched)getNodeConfig()(*saveMatchedConfig){
	mapData,_:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	jsonStr, err := json.Marshal(mapData)
    if err != nil {
        log.Println(err)
		return nil
    }
	conf:=&saveMatchedConfig{}
    if err := json.Unmarshal(jsonStr, conf); err != nil {
        log.Println(err)
		return nil
    }

	return conf
}

func (nodeExecutor *nodeExecutorSaveMatched)aggregationFirst(
	groupRow,modelRow map[string]interface{},
	groupField *groupModelField,
	index int){

	if index == 0 {
		value,ok:=modelRow[groupField.SourceField]
		if ok {
			groupRow[groupField.Field]=value
		}
	}
}

func (nodeExecutor *nodeExecutorSaveMatched)aggregationSum(
	groupRow,modelRow map[string]interface{},
	groupField *groupModelField,
	index int){
	
	//获取之前的数据
	preVal,ok:=groupRow[groupField.Field]
	if !ok {
		log.Println("nodeExecutorSaveMatched aggragationSum source field "+groupField.Field+" not found")
		preVal="0"
	}
	spreVal,_:=preVal.(string)
	preValFloat64,err:=strconv.ParseFloat(spreVal, 64)
	if err !=nil {
		log.Println("nodeExecutorSaveMatched aggragationSum can not convert aggregeted value to float64")
		return
	}
	//获取当前数据
	curVal,ok:=modelRow[groupField.SourceField]
	if !ok {
		log.Println("nodeExecutorSaveMatched aggragationSum no source field:"+groupField.SourceField)
		return
	}
	scurVal,_:=curVal.(string)
	curValFloat64,err:=strconv.ParseFloat(scurVal, 64)
	if err !=nil {
		log.Println("nodeExecutorSaveMatched aggragationSum can not convert source value to float64")
		return
	}

	newVal:=preValFloat64+curValFloat64
	groupRow[groupField.Field]=fmt.Sprint(newVal)
}

func (nodeExecutor *nodeExecutorSaveMatched)aggragation(
	groupRow,modelRow map[string]interface{},
	groupField *groupModelField,
	index int){
	
	if groupField.Aggregation == AGGREGATION_FIRST {
		nodeExecutor.aggregationFirst(groupRow,modelRow,groupField,index)
	} else if groupField.Aggregation == AGGREGATION_SUM {
		nodeExecutor.aggregationSum(groupRow,modelRow,groupField,index)
	}
}

func (nodeExecutor *nodeExecutorSaveMatched)createGroup(
	dataItem *flowDataItem,
	saveMatchedConf *saveMatchedConfig,
	tx *sql.Tx,
	instance *flowInstance,
	groupIndex int,
	batchNo string)(string,int){

	groupID:=fmt.Sprintf("%s%05d",batchNo,groupIndex)
	
	groupRow:=map[string]interface{}{}
	//根据配置生成分组记录数据
	//迭代每个model
	for _,modelData:=range dataItem.Models {
		if modelData.List !=nil && len(*modelData.List)>0 {
			for index,modelRow:=range  (*modelData.List) {
				for _,groupField:=range saveMatchedConf.GroupModel.Fields {
					if groupField.SourceModel == *modelData.ModelID {
						nodeExecutor.aggragation(groupRow,modelRow,&groupField,index)
					}
				}
			}
		}
	}

	groupRow[CC_MATCH_RESULT]=saveMatchedConf.MatchResult
	groupRow[data.SAVE_TYPE_COLUMN]=data.SAVE_CREATE
	groupRow[data.CC_ID]=groupID

	saver:=data.Save{
		List:&[]map[string]interface{}{
			groupRow,
		},
		AppDB:instance.AppDB,
		UserID:instance.UserID,
		ModelID:saveMatchedConf.GroupModel.ModelID,
	}

	result,errorCode:=saver.SaveList(nodeExecutor.DataRepository,tx)
	if errorCode != common.ResultSuccess {
		return "",errorCode
	}
	log.Println("nodeExecutorSaveMatched createGroup ",result.List)
	/*resultID,_:=result.List[0][data.CC_ID]
	var groupID string
	switch resultID.(type){
	case string:
		groupID,_=resultID.(string)
	case int64:
		intVal,_:=resultID.(int64)
		groupID=strconv.FormatInt(intVal,10)
	} 
	log.Printf("nodeExecutorSaveMatched resultID type %T",resultID)
	log.Println("nodeExecutorSaveMatched groupID ",groupID)*/
	return groupID,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorSaveMatched)saveModels(
	dataItem *flowDataItem,
	groupID string,
	tx *sql.Tx,
	instance *flowInstance)(int){

	for _,modelData:=range (dataItem.Models) {
		if modelData.List!=nil && len(*modelData.List)>0 {
			rows:=make([]map[string]interface{},len(*modelData.List))
			for index,row:=range (*modelData.List) {
				rows[index]=map[string]interface{}{}
				rows[index][data.CC_ID]=row[data.CC_ID]
				rows[index][data.CC_VERSION]=row[data.CC_VERSION]
				rows[index][data.SAVE_TYPE_COLUMN]=data.SAVE_UPDATE
				rows[index][CC_MATCH_GROUP]=groupID
				rows[index][CC_MATCH_STATUS]=MATCH_STATUS_SUCCESS
				rows[index][CC_MATCH_MESSAGE]=""
			}

			saver:=data.Save{
				List:&rows,
				AppDB:instance.AppDB,
				UserID:instance.UserID,
				ModelID:*modelData.ModelID,
			}
			_,errorCode:=saver.SaveList(nodeExecutor.DataRepository,tx)
			if errorCode != common.ResultSuccess {
				return errorCode
			}
		}
	}

	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorSaveMatched)save(
	dataItem *flowDataItem,
	saveMatchedConf *saveMatchedConfig,
	instance *flowInstance,
	groupIndex int,
	batchNo string){
	
	log.Println("start nodeExecutorSaveMatched save")
	
	//开启事务
	tx,err:= nodeExecutor.DataRepository.Begin()
	if err != nil {
		log.Println(err)
		return
	}

	//创建分组记录数据
	groupID,errorCode:=nodeExecutor.createGroup(dataItem,saveMatchedConf,tx,instance,groupIndex,batchNo)
	if errorCode!=common.ResultSuccess {
		tx.Rollback()
		log.Println("end nodeExecutorSaveMatched save")
		return
	}
	//将分组号更新到左右表，同时更新左右表数据的匹配状态
	errorCode=nodeExecutor.saveModels(dataItem,groupID,tx,instance)
	if errorCode!=common.ResultSuccess {
		tx.Rollback()
		log.Println("end nodeExecutorSaveMatched save")
		return
	}
	
	//提交事务
	if err := tx.Commit(); err != nil {
		log.Println(err)
	}
	log.Println("end nodeExecutorSaveMatched save")
}

func (nodeExecutor *nodeExecutorSaveMatched)getBatchNumber()(string){
	return xid.New().String()
}

func (nodeExecutor *nodeExecutorSaveMatched)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	log.Println("nodeExecutorSaveMatched run start")

	jsonStr, _ := json.MarshalIndent(node, "", "    ")
	log.Println(string(jsonStr))

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
		"nodeType":NODE_SAVE_MATCHED,
	}
	//加载节点配置
	saveMatchedConf:=nodeExecutor.getNodeConfig()
	if saveMatchedConf==nil {
		log.Printf("nodeExecutorSaveMatched run get node config error\n")
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if req.Data==nil || len(*req.Data)==0 {
		endTime:=time.Now().Format("2006-01-02 15:04:05")
		node.Completed=true
		node.EndTime=&endTime
		node.Output=req
		log.Printf("nodeExecutorSaveMatched run end without data\n")
		return req,nil
	}

	batchNo:=nodeExecutor.getBatchNumber()

	//一个分组作为一个独立的事务进行保存
	for index,item:= range (*req.Data) {
		nodeExecutor.save(&item,saveMatchedConf,instance,index,batchNo)
	}
	
	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult
	log.Println("nodeExecutorSaveMatched run end")
	return flowResult,nil
}