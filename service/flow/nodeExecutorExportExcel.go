package flow

import (
  "time"
	"dataflow/common"
	"encoding/json"
	"log"
	"github.com/xuri/excelize/v2"
	"net/url"
)

type ExportExcelField struct {
	Field string `json:"field"`
	Label string `json:"label"`
}

type ExportExcelSheet struct {
	SheetName string `json:"sheetName"`
	ModelID string `json:"modelID"`
	Fields []ExportExcelField `json:"fields"`
}

type nodeExportExcelConf struct {
	FileName string `json:"fileName"`
	Sheets []ExportExcelSheet `json:"sheets"`
}

type nodeExecutorExportExcel struct {
	NodeConf node
}

func (nodeExecutor *nodeExecutorExportExcel)getNodeConf()(*nodeExportExcelConf){
	mapData,ok:=nodeExecutor.NodeConf.Data.(map[string]interface{})
	if !ok {
		return nil
	}
	jsonStr, err := json.Marshal(mapData)
  if err != nil {
    log.Println(err)
		return nil
  }
	conf:=nodeExportExcelConf{}
  if err := json.Unmarshal(jsonStr, &conf); err != nil {
    log.Println(err)
		return nil
  }

	return &conf
}

func (nodeExecutor *nodeExecutorExportExcel)writeHeader(
	file *excelize.File,
	sheet *ExportExcelSheet){
		row:=1
		for col,field:=range sheet.Fields {
			cellStart,_:=excelize.CoordinatesToCellName(col+1, row)
			file.SetCellStr(sheet.SheetName,cellStart,field.Label)
		}
}

func (nodeExecutor *nodeExecutorExportExcel)writeDdataRow(
	file *excelize.File,
	sheet *ExportExcelSheet,
	rowNo int,
	rowData *map[string]interface{}){
	for col,field:=range sheet.Fields {
		value,ok:=(*rowData)[field.Field]
		if ok && value != nil {
			cellStart,_:=excelize.CoordinatesToCellName(col+1, rowNo)
			file.SetCellValue(sheet.SheetName,cellStart,value)
		}
	}
}

func (nodeExecutor *nodeExecutorExportExcel)createExcelSheet(
	file *excelize.File,
	sheet *ExportExcelSheet,
	dataItem *modelDataItem)(int){
	//创建sheet
	sheetIndex := file.NewSheet(sheet.SheetName)
    file.SetActiveSheet(sheetIndex)
	//生成header行
	nodeExecutor.writeHeader(file,sheet)
	//写入数据行
	for row,rowData:=range *dataItem.List {
		nodeExecutor.writeDdataRow(file,sheet,row+2,&rowData)
	}

	return common.ResultSuccess
}

func (nodeExecutor *nodeExecutorExportExcel)createExcelFile(
	dataItem *flowDataItem,
	nodeConf *nodeExportExcelConf)(*excelize.File,int){
		f := excelize.NewFile()
		for _,sheetConf:=range nodeConf.Sheets {
			for _,modelDataItem:=range dataItem.Models {
				if sheetConf.ModelID == *modelDataItem.ModelID {
					errorCode:=nodeExecutor.createExcelSheet(f,&sheetConf,&modelDataItem)
					if errorCode!= common.ResultSuccess {
						f=nil
						return nil,errorCode
					}
				}
			}
		}
		f.DeleteSheet("Sheet1")
		return f,common.ResultSuccess
}

func (nodeExecutor *nodeExecutorExportExcel)run(
	instance *flowInstance,
	node,preNode *instanceNode)(*flowReqRsp,*common.CommonError){

	params:=map[string]interface{}{
		"nodeID":node.ID,
		"nodeType":NODE_EXPORT_EXCEL,
	}

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
		Over:true,
	}
	
	//加载节点配置
	nodeConf:=nodeExecutor.getNodeConf()
	if nodeConf==nil {
		log.Printf("nodeExecutorExportExcel run get node config error\n")
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if len(nodeConf.FileName)==0 {
		log.Printf("nodeExecutorExportExcel node config without filename\n")
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if len(nodeConf.Sheets)==0 {
		log.Printf("nodeExecutorExportExcel node config without sheets\n")
		return flowResult,common.CreateError(common.ResultNodeConfigError,params)
	}

	if req.Data==nil || len(*req.Data)<=0 {
		log.Printf("nodeExecutorExportExcel no data to export\n")
		return flowResult,common.CreateError(common.ResultNoDataForExport,params)
	}

	if instance.GinContext == nil {
		log.Printf("nodeExecutorExportExcel node can not use in aysnc flow and sub flow.\n")
		return flowResult,common.CreateError(common.ResultNotSupportedNode,params)
	}

	//根据配置生成Excel文件
	file,errorCode:=nodeExecutor.createExcelFile(&(*req.Data)[0],nodeConf)
	if errorCode !=  common.ResultSuccess {
		return flowResult,common.CreateError(errorCode,params)
	}

	//文件流写入应答对象
	instance.GinContext.Header("Content-Type", "application/octet-stream")
	filename:=url.QueryEscape(nodeConf.FileName+"_"+instance.InstanceID+".xlsx")
  instance.GinContext.Header("Content-Disposition", "attachment;filename="+filename)
  instance.GinContext.Header("Content-Transfer-Encoding", "binary")
	
	file.Write(instance.GinContext.Writer)

	//通知外层controller不需要回复结果
	flowResult.AlreadyResponsed=true

	endTime:=time.Now().Format("2006-01-02 15:04:05")
	node.Completed=true
	node.EndTime=&endTime
	node.Output=flowResult

	return flowResult,nil
}