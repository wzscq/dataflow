import {Row,Col,Button,Input,Select,InputNumber} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined,MinusOutlined,AlignCenterOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';
import {openDialog} from '../../../../redux/dialogSlice';

const {Option}=Select;

export default function ESINodeParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onModelIDChange=(e)=>{
      dispatch(updateNodeData({...node.data,modelID:e.target.value}));
    }

    const setShowFields=(showFields)=>{
      dispatch(updateNodeData({...node.data,__showFields:showFields}));
    }

    const setShowOptions=(showOptions)=>{
      dispatch(updateNodeData({...node.data,__showOptions:showOptions}));
    }

    const onAddField=()=>{
      let fields=node.data.fields?node.data.fields:[];
      fields=[...fields,{field:"",__showField:true,labelRegexp:"",excelRangeType:"",endRow:"",emptyValue:""}];
      dispatch(updateNodeData({...node.data,fields:fields}));
    }

    const onDelField=(index)=>{
      let fields=[...node.data.fields];
      delete fields[index];
      dispatch(updateNodeData({...node.data,fields:fields.filter(item=>item)}));
    }

    const onFieldChange=(index,value)=>{
      let fields=[...node.data.fields];
      fields[index]=value;
      dispatch(updateNodeData({...node.data,fields:fields}));
    }

    const onOptionChange=(value)=>{
      dispatch(updateNodeData({...node.data,options:value}));
    }

    const onSetTestData=()=>{
      dispatch(openDialog({type:'setESITestData',title:'ESI Test Data',param:{node:node}}));
    }

    const showFields=node.data?.__showFields===false?false:true;
    const showOptions=node.data?.__showOptions===false?false:true;
    const options=node.data?.options?node.data.options:{};
    const testData=node.data?.testData?node.data.testData:{};
    let testFilename="";
    if(testData.esiFile?.list?.length>0){
      testFilename=testData.esiFile.list[0].name;
    }

    const fields=node.data?.fields?.map((item,index)=>{
      const showField=item.__showField;
      return (
        <>
        <Row className="param-panel-row" style={{display:showFields?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>
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
        <Row className="param-panel-row" style={{display:showFields&&showField?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Field</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Input value={item.field} onChange={(e)=>onFieldChange(index,{...item,field:e.target.value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showFields&&showField?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>LabelRegexp</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Input value={item.labelRegexp} onChange={(e)=>onFieldChange(index,{...item,labelRegexp:e.target.value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showFields&&showField?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>ExcelRangeType</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Select value={item?.excelRangeType} size='small' onChange={(value)=>onFieldChange(index,{...item,excelRangeType:value})}>
              <Option key='table'>Table</Option>
              <Option key='rightVal'>Right Value</Option>
              <Option key='auto'>Auto Detect</Option>
            </Select>
          </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showFields&&showField?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>EndRow</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Select value={item?.endRow} size='small' onChange={(value)=>onFieldChange(index,{...item,endRow:value})}>
                <Option key='no'>No</Option>
                <Option key='yes'>Yes</Option>
                <Option key='auto'>Auto Detect</Option>
              </Select>
          </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showFields&&showField?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>EmptyValue</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Select value={item?.emptyValue} size='small' onChange={(value)=>onFieldChange(index,{...item,emptyValue:value})}>
                <Option key='no'>No</Option>
                <Option key='yes'>Yes</Option>
              </Select>
          </Col>
        </Row>
        </>
      );
    });

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>Model ID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.modelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>
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
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>
                <div className='button' onClick={()=>setShowOptions(!showOptions)}>
                    {showOptions?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Options</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showOptions?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Prevent Same File</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Select value={options.preventSameFile===true?"true":"false"} size='small' onChange={(value)=>onOptionChange({...options,preventSameFile:value==="true"?true:false})}>
                <Option key={"false"}>No</Option>
                <Option key={"true"}>Yes</Option>
            </Select>
          </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showOptions?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>GenerateRowID</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
              <Select value={options.generateRowID===true?"true":"false"} size='small' onChange={(value)=>onOptionChange({...options,generateRowID:value==="true"?true:false})}>
                <Option key={"false"}>No</Option>
                <Option key={"true"}>Yes</Option>
              </Select>
          </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showOptions?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>MaxHeaderRow</Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <InputNumber size='small' value={options.maxHeaderRow} onChange={(value)=>onOptionChange({...options,maxHeaderRow:value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>
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

