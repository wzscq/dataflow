import { useDispatch } from "react-redux";
import {Row,Col,Button,Input} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusOutlined,MinusSquareOutlined,AlignCenterOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';
import {openDialog} from '../../../../redux/dialogSlice';

const initContent=
`/*
数据行字段处理函数，可用于生成字段的值，或根据其它字段值计算某个字段的值
入参如下：
field:当前处理的字段名称
row:当前处理的数据行
返回值：
直接返回需要处理的字段的值,比如：  return fieldValue;
*/`

export default function TransferModelItem({node,modelIndex,labelWidth}){
    const dispatch=useDispatch();

    const setShowModel=()=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],__showModel:!(models[modelIndex].__showModel)}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const setShowFields=()=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],__showFields:!(models[modelIndex].__showFields)}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onModelIDChange=(e)=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],modelID:e.target.value}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onDelModel=()=>{
        const models=[...node.data?.models];
        delete models[modelIndex];
        dispatch(updateNodeData({...node.data,models:models.filter(item=>item)}));
    }

    const onAddField=()=>{
        const models=[...node.data?.models];
        const modelFields=models[modelIndex].fields?[...(models[modelIndex].fields)]:[];
        modelFields.push({sourceField:"",targetField:"",keyField:'0',__showField:true});
        models[modelIndex]={...models[modelIndex],fields:modelFields};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onDelField=(index)=>{
        const models=[...node.data?.models];
        const modelFields=models[modelIndex].fields?[...(models[modelIndex].fields)]:[];
        delete modelFields[index];
        models[modelIndex]={...models[modelIndex],fields:modelFields.filter(item=>item)};
        dispatch(updateNodeData({...node.data,models:models}));
    }
    
    const onFieldChange=(index,value)=>{
        const models=[...node.data?.models];
        const modelFields=[...(models[modelIndex].fields)];
        modelFields[index]=value;
        models[modelIndex]={...models[modelIndex],fields:modelFields};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const setShowField=(index,showField)=>{
        const models=[...node.data?.models];
        const modelFields=[...(models[modelIndex].fields)];
        modelFields[index]={...modelFields[index],__showField:showField};
        models[modelIndex]={...models[modelIndex],fields:modelFields};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onEditScript=(fieldIndex)=>{
        dispatch(openDialog({type:'funcScript',title:'Function Script',param:{node:node,modelIndex:modelIndex,fieldIndex:fieldIndex,initContent}}));
    }

    const modelItem=node.data.models[modelIndex];
    const showModels=node.data?.__showModels===false?false:true;
    const showModel=modelItem.__showModel;
    const showFields=modelItem.__showFields;

    const fields=modelItem.fields?.map((item,index)=>{
        const showField=item.__showField
        return (
            <>
            <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields?"flex":"none"}} gutter={24}>
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
            <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Field</Col>
                <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Input value={item.field} onChange={(e)=>onFieldChange(index,{...item,field:e.target.value})}/>
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Function Script</Col>
                <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Input disabled={true} value={item.funcScript?.name}/>
                    <Button className="button"  onClick={()=>{onEditScript(index)}} size='small' icon={<AlignCenterOutlined />} />
                </Col>
            </Row>
            </>
        );
    });

    const modelItemControl=(
    <>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Model ID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={modelItem.modelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
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
        <Row className="param-panel-row" style={{display:showModels?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>
                <div className='button' onClick={setShowModel}>
                {showModel?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Model {modelIndex}</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onDelModel} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        {modelItemControl}
      </>);
}