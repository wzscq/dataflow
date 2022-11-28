import { useDispatch } from "react-redux";
import {Row,Col,Button,Input} from 'antd';
import { PlusSquareOutlined,MinusOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';

export default function GroupModelItem({node,modelIndex}){
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
            <Col className="param-panel-row-label level-2" span={10}>Model ID</Col>
            <Col className="param-panel-row-input" span={14}>
                <Input value={modelItem.modelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" span={10}>Field</Col>
            <Col className="param-panel-row-input" span={14}>
                <Input value={modelItem.field} onChange={onFieldChange}/>
            </Col>
        </Row>
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