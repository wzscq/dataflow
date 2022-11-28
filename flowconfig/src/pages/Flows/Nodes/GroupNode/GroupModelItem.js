import { useDispatch } from "react-redux";
import {Row,Col,Button,Input} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';

export default function GroupModelItem({node,modelIndex}){
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
        modelFields.push({field:""});
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
        modelFields[index]={...modelFields[index],field:value};
        models[modelIndex]={...models[modelIndex],fields:modelFields};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const modelItem=node.data.models[modelIndex];
    const showModels=node.data?.__showModels===false?false:true;
    const showModel=modelItem.__showModel;
    const showFields=modelItem.__showFields;

    const fields=modelItem.fields?.map((item,index)=>{
        return (
        <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-3" span={10}>Field {index}</Col>
            <Col className="param-panel-row-inputwithbutton" span={14}>
                <Input value={item.field} onChange={(e)=>onFieldChange(index,e.target.value)}/>
                <Button className="button"  onClick={()=>onDelField(index)} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        );
    });

    const modelItemControl=(
    <>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" span={10}>Model ID</Col>
            <Col className="param-panel-row-input" span={14}>
                <Input value={modelItem.modelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" span={10}>
                <div className='button' onClick={setShowFields}>
                    {showFields?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Fields</span>
            </Col>
            <Col className="param-panel-row-input" span={14}>
                <Button className="button"  onClick={onAddField} size='small' icon={<PlusOutlined />} />
            </Col>
        </Row>
        {fields}
    </>);
    
    return (
      <>  
        <Row className="param-panel-row" style={{display:showModels?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-1" span={10}>
                <div className='button' onClick={setShowModel}>
                {showModel?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Model {modelIndex}</span>
            </Col>
            <Col className="param-panel-row-input" span={14}>
                <Button className="button"  onClick={onDelModel} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        {modelItemControl}
      </>);
}