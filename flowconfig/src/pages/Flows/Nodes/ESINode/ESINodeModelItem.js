import {Row,Col,Button,Input,Select,InputNumber} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined,MinusOutlined,AlignCenterOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';
import {openDialog} from '../../../../redux/dialogSlice';

const {Option}=Select;


export default function ESINodeModelItem({node,modelIndex,labelWidth}){
  const dispatch=useDispatch();

  const onModelIDChange=(e)=>{
    const models=[...node.data?.models];
    models[modelIndex]={...models[modelIndex],modelID:e.target.value}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onShowModelChange=()=>{
    const models=[...node.data?.models];
    models[modelIndex]={...models[modelIndex],__showModel:!(models[modelIndex].__showModel)}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onDelModel=(modelIndex)=>{
    const models=[...node.data?.models];
    delete models[modelIndex];
    dispatch(updateNodeData({...node.data,models:models.filter(item=>item)}));
  }

  const onFileModelChange=(e)=>{
    const models=[...node.data?.models];
    models[modelIndex]={...models[modelIndex],fileModel:e.target.value}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onFileFieldChange=(e)=>{
    const models=[...node.data?.models];
    models[modelIndex]={...models[modelIndex],fileField:e.target.value}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const setShowFields=(showFields)=>{
    const models=[...node.data?.models];
    models[modelIndex]={...models[modelIndex],__showFields:showFields}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const setShowOptions=(showOptions)=>{
    const models=[...node.data?.models];
    models[modelIndex]={...models[modelIndex],__showOptions:showOptions}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const setShowSheets=(showSheets)=>{
    const models=[...node.data?.models];
    models[modelIndex]={...models[modelIndex],__showSheets:showSheets}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onAddField=()=>{
    const models=[...node.data?.models];
    let fields=models[modelIndex].fields?models[modelIndex].fields:[];
    fields=[...fields,{field:"",__showField:true,labelRegexp:"",excelRangeType:"",endRow:"",emptyValue:""}];
    models[modelIndex]={...models[modelIndex],fields:fields}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onAddSheet=()=>{
    const models=[...node.data?.models];
    let sheets=models[modelIndex].sheets?models[modelIndex].sheets:[];
    sheets=[...sheets,{type:"name",__showSheet:true,optional:"no",value:""}];
    models[modelIndex]={...models[modelIndex],sheets:sheets}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onDelField=(index)=>{
    const models=[...node.data?.models];
    let fields=[...(models[modelIndex].fields?models[modelIndex].fields:[])];
    delete fields[index];
    models[modelIndex]={...models[modelIndex],fields:fields.filter(item=>item)}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onDelSheet=(index)=>{
    const models=[...node.data?.models];
    let sheets=[...(models[modelIndex].sheets?models[modelIndex].sheets:[])];
    delete sheets[index];
    models[modelIndex]={...models[modelIndex],sheets:sheets.filter(item=>item)}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onFieldChange=(index,value)=>{
    const models=[...node.data?.models];
    let fields=[...(models[modelIndex].fields?models[modelIndex].fields:[])];
    fields[index]=value;
    models[modelIndex]={...models[modelIndex],fields:fields}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onSheetChange=(index,value)=>{
    const models=[...node.data?.models];
    let sheets=[...(models[modelIndex].sheets?models[modelIndex].sheets:[])];
    sheets[index]=value;
    models[modelIndex]={...models[modelIndex],sheets:sheets}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onOptionChange=(value)=>{
    const models=[...node.data?.models];
    models[modelIndex]={...models[modelIndex],options:value}
    dispatch(updateNodeData({...node.data,models:models}));
  }

  const onSetTestData=()=>{
    dispatch(openDialog({type:'setESITestData',title:'ESI Test Data',param:{node:node,modelIndex:modelIndex}}));
  }

  const modelItem=node.data.models[modelIndex];
  const showModels=node.data?.__showModels===false?false:true;
  const showModel=modelItem.__showModel===false?false:true;
  const showFields=modelItem.__showFields===false?false:true;
  const showOptions=modelItem.__showOptions===false?false:true;
  const showSheets=modelItem.__showSheets===false?false:true;
  const options=modelItem.options?modelItem.options:{};
  const testData=modelItem.testData?modelItem.testData:{};
  let testFilename="";
  if(testData.esiFile?.list?.length>0){
    testFilename=testData.esiFile.list[0].name;
  }

  const sheets=modelItem.sheets?.map((item,index)=>{
    const showSheet=item.__showSheet;
    return (
      <>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showSheets?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>
              <div className='button' onClick={()=>onSheetChange(index,{...item,__showSheet:!showSheet})}>
                  {showSheet?<MinusSquareOutlined />:<PlusSquareOutlined />}
              </div>
              <span>Sheet {index}</span>
          </Col>
          <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Button className="button"  onClick={()=>onDelSheet(index)} size='small' icon={<MinusOutlined />} />
          </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showSheets&&showSheet?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Type</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
          <Select value={item?.type} size='small' onChange={(value)=>onSheetChange(index,{...item,type:value})}>
            <Option key='index'>Index</Option>
            <Option key='name'>Name</Option>
          </Select>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showSheets&&showSheet?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Value</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={item.value} onChange={(e)=>onSheetChange(index,{...item,value:e.target.value})}/>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showSheets&&showSheet?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Optional</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Select value={item?.optional} size='small' onChange={(value)=>onSheetChange(index,{...item,optional:value})}>
              <Option key='no'>No</Option>
              <Option key='yes'>Yes</Option>
            </Select>
        </Col>
      </Row>
      </>)
  });

  const fields=modelItem.fields?.map((item,index)=>{
    const showField=item.__showField;
    return (
      <>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>
              <div className='button' onClick={()=>onFieldChange(index,{...item,__showField:!showField})}>
                  {showField?<MinusSquareOutlined />:<PlusSquareOutlined />}
              </div>
              <span>Field {index}</span>
          </Col>
          <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <span>{item.field}</span>
              <Button className="button"  onClick={()=>onDelField(index)} size='small' icon={<MinusOutlined />} />
          </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Field</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={item.field} onChange={(e)=>onFieldChange(index,{...item,field:e.target.value})}/>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>LabelRegexp</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={item.labelRegexp} onChange={(e)=>onFieldChange(index,{...item,labelRegexp:e.target.value})}/>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>ExcelRangeType</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
          <Select value={item?.excelRangeType} size='small' onChange={(value)=>onFieldChange(index,{...item,excelRangeType:value})}>
            <Option key='table'>Table</Option>
            <Option key='rightVal'>Right Value</Option>
            <Option key='auto'>Auto Detect</Option>
          </Select>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>EndRow</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Select value={item?.endRow} size='small' onChange={(value)=>onFieldChange(index,{...item,endRow:value})}>
              <Option key='no'>No</Option>
              <Option key='yes'>Yes</Option>
              <Option key='auto'>Auto Detect</Option>
            </Select>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>EmptyValue</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Select value={item?.emptyValue} size='small' onChange={(value)=>onFieldChange(index,{...item,emptyValue:value})}>
              <Option key='no'>No</Option>
              <Option key='yes'>Yes</Option>
            </Select>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Source</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Select value={item?.source} size='small' onChange={(value)=>onFieldChange(index,{...item,source:value})}>
              <Option key='input'>Input</Option>
              <Option key='file'>File</Option>
            </Select>
        </Col>
      </Row>
      </>
    );
  });

  return (
    <>
      <Row className="param-panel-row" style={{display:showModels?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>
              <div className='button' onClick={()=>onShowModelChange(!showModel)}>
                  {showModel?<MinusSquareOutlined />:<PlusSquareOutlined />}
              </div>
              <span>Model {modelIndex}</span>
          </Col>
          <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Button className="button"  onClick={()=>onDelModel(modelIndex)} size='small' icon={<MinusOutlined />} />
          </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Model ID</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Input value={modelItem.modelID} onChange={onModelIDChange}/>
          </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>File Model</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Input value={modelItem.fileModel} onChange={onFileModelChange}/>
          </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>File Field</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Input value={modelItem.fileField} onChange={onFileFieldChange}/>
          </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
              <div className='button' onClick={()=>setShowSheets(!showSheets)}>
                  {showSheets?<MinusSquareOutlined />:<PlusSquareOutlined />}
              </div>
              <span>Sheets</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Button className="button"  onClick={onAddSheet} size='small' icon={<PlusOutlined />} />
          </Col>
      </Row>
      {sheets}
      <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
              <div className='button' onClick={()=>setShowFields(!showFields)}>
                  {showFields?<MinusSquareOutlined />:<PlusSquareOutlined />}
              </div>
              <span>Fields</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Button className="button"  onClick={onAddField} size='small' icon={<PlusOutlined />} />
          </Col>
      </Row>
      {fields}
      <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
              <div className='button' onClick={()=>setShowOptions(!showOptions)}>
                  {showOptions?<MinusSquareOutlined />:<PlusSquareOutlined />}
              </div>
              <span>Options</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              
          </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showOptions?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>Prevent Same File</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
          <Select value={options.preventSameFile===true?"true":"false"} size='small' onChange={(value)=>onOptionChange({...options,preventSameFile:value==="true"?true:false})}>
              <Option key={"false"}>No</Option>
              <Option key={"true"}>Yes</Option>
          </Select>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showOptions?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>GenerateRowID</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Select value={options.generateRowID===true?"true":"false"} size='small' onChange={(value)=>onOptionChange({...options,generateRowID:value==="true"?true:false})}>
              <Option key={"false"}>No</Option>
              <Option key={"true"}>Yes</Option>
            </Select>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel&&showOptions?"flex":"none"}} gutter={24}>
        <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>MaxHeaderRow</Col>
        <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
          <InputNumber size='small' value={options.maxHeaderRow} onChange={(value)=>onOptionChange({...options,maxHeaderRow:value})}/>
        </Col>
      </Row>
      <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
              TestData
          </Col>
          <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Input disabled={true} value={testFilename}/>
              <Button className="button"  onClick={onSetTestData} size='small' icon={<AlignCenterOutlined />} />
          </Col>
      </Row>
    </>
  );
}