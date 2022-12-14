import { useDispatch } from "react-redux";
import {Row,Col,Button,Input,Select} from 'antd';
import { PlusSquareOutlined,MinusOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';

const { Option } = Select;

export default function VerifyValueModelItem({node,modelIndex,labelWidth}){
    const dispatch=useDispatch();

    const setShowModel=()=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],__showModel:!(models[modelIndex].__showModel)}
        dispatch(updateNodeData({...node.data,models:models}));
    }
    
    const onModelIDChange=(e)=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],modelID:e.target.value}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onFieldChange=(e)=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],field:e.target.value}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onSideChange=(value)=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],side:value}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onDelModel=()=>{
        const models=[...node.data?.models];
        delete models[modelIndex];
        dispatch(updateNodeData({...node.data,models:models.filter(item=>item)}));
    }

    const modelItem=node.data.models[modelIndex];
    const showModels=node.data?.__showModels===false?false:true;
    const showModel=modelItem.__showModel;

    const modelItemControl=(
    <>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Model ID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={modelItem.modelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Field</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={modelItem.field} onChange={onFieldChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Side</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={modelItem.side} size='small' onChange={onSideChange}>
                    <Option key='left'>Left</Option>
                    <Option key='right'>Right</Option>
                </Select>
            </Col>
        </Row>
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