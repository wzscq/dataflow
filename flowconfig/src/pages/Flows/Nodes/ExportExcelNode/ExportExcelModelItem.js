import { useDispatch } from "react-redux";
import {Row,Col,Button,Input} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';

export default function RequestQueryModelItem({node,sheetIndex,labelWidth}){
    const dispatch=useDispatch();

    const setShowSheet=()=>{
        const sheets=[...node.data?.sheets];
        sheets[sheetIndex]={...sheets[sheetIndex],__showSheet:!(sheets[sheetIndex].__showSheet)}
        dispatch(updateNodeData({...node.data,sheets:sheets}));
    }

    const setShowFields=()=>{
        const sheets=[...node.data?.sheets];
        sheets[sheetIndex]={...sheets[sheetIndex],__showFields:!(sheets[sheetIndex].__showFields)}
        dispatch(updateNodeData({...node.data,sheets:sheets}));
    }

    const setShowField=(index,showField)=>{
        const sheets=[...node.data?.sheets];
        const sheetFields=[...(sheets[sheetIndex].fields)];
        sheetFields[index]={...sheetFields[index],__showField:showField};
        sheets[sheetIndex]={...sheets[sheetIndex],fields:sheetFields};
        dispatch(updateNodeData({...node.data,sheets:sheets}));
    }
    
    const onModelIDChange=(e)=>{
        const sheets=[...node.data?.sheets];
        sheets[sheetIndex]={...sheets[sheetIndex],modelID:e.target.value}
        dispatch(updateNodeData({...node.data,sheets:sheets}));
    }

    const onSheetNameChange=(e)=>{
        const sheets=[...node.data?.sheets];
        sheets[sheetIndex]={...sheets[sheetIndex],sheetName:e.target.value}
        dispatch(updateNodeData({...node.data,sheets:sheets}));
    }

    const onDelSheet=()=>{
        const sheets=[...node.data?.sheets];
        delete sheets[sheetIndex];
        dispatch(updateNodeData({...node.data,sheets:sheets.filter(item=>item)}));
    }

    const onAddField=()=>{
        const sheets=[...node.data?.sheets];
        const sheetFields=sheets[sheetIndex].fields?[...(sheets[sheetIndex].fields)]:[];
        sheetFields.push({field:"",__showField:true});
        sheets[sheetIndex]={...sheets[sheetIndex],fields:sheetFields};
        dispatch(updateNodeData({...node.data,sheets:sheets}));
    }

    const onDelField=(index)=>{
        const sheets=[...node.data?.sheets];
        const sheetFields=sheets[sheetIndex].fields?[...(sheets[sheetIndex].fields)]:[];
        delete sheetFields[index];
        sheets[sheetIndex]={...sheets[sheetIndex],fields:sheetFields.filter(item=>item)};
        dispatch(updateNodeData({...node.data,sheets:sheets}));
    }

    const onFieldChange=(index,value)=>{
        const sheets=[...node.data?.sheets];
        const sheetFields=[...(sheets[sheetIndex].fields)];
        sheetFields[index]={...sheetFields[index],...value};
        sheets[sheetIndex]={...sheets[sheetIndex],fields:sheetFields};
        dispatch(updateNodeData({...node.data,sheets:sheets}));
    }

    const sheetItem=node.data.sheets[sheetIndex];
    const showSheets=node.data?.__showSheets===false?false:true;
    const showSheet=sheetItem.__showSheet;
    const showFields=sheetItem.__showFields;

    const fields=sheetItem.fields?.map((item,index)=>{
        const showField=item.__showField
        return (
            <>
            <Row className="param-panel-row" style={{display:showSheets&&showSheet&&showFields?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>
                    <div className='button' onClick={()=>setShowField(index,!showField)}>
                        {showField?<MinusSquareOutlined />:<PlusSquareOutlined />}
                    </div>
                    <span>Field {index}</span>
                </Col>
                <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <span>{item.targetField}</span>
                    <Button className="button"  onClick={()=>onDelField(index)} size='small' icon={<MinusOutlined />} />
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showSheets&&showSheet&&showFields&&showField?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Field</Col>
                <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Input value={item.field} onChange={(e)=>onFieldChange(index,{...item,field:e.target.value})}/>
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showSheets&&showSheet&&showFields&&showField?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Label</Col>
                <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Input value={item.label} onChange={(e)=>onFieldChange(index,{...item,label:e.target.value})}/>
                </Col>
            </Row>
            </>
        );
    });

    const sheetItemControl=(
    <>
        <Row className="param-panel-row" style={{display:showSheets&&showSheet?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Sheet Name</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={sheetItem.sheetName} onChange={onSheetNameChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSheets&&showSheet?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Model ID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={sheetItem.modelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSheets&&showSheet?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
                <div className='button' onClick={setShowFields}>
                    {showFields?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Fields</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onAddField} size='small' icon={<PlusOutlined />} />
            </Col>
        </Row>
        {fields}
    </>);
    
    return (
      <>  
        <Row className="param-panel-row" style={{display:showSheets?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>
                <div className='button' onClick={setShowSheet}>
                {showSheet?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Sheet {sheetIndex}</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onDelSheet} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        {sheetItemControl}
      </>);
}