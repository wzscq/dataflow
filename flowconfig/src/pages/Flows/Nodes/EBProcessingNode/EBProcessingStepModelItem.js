import { useDispatch } from "react-redux";
import {Row,Col,Button,Input,Select} from 'antd';
import { PlusSquareOutlined,MinusOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';

const { Option } = Select;

export default function EBProcessingStepModelItem({node,stepIndex,modelIndex,labelWidth}){
    const dispatch=useDispatch();

    const setShowModel=()=>{
        const bpSteps=[...node.data?.bpSteps];
        const models=[...bpSteps[stepIndex].models];
        models[modelIndex]={...models[modelIndex],__showModel:!models[modelIndex].__showModel};
        bpSteps[stepIndex]={...bpSteps[stepIndex],models:models}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const onModelIDChange=(e)=>{
        const bpSteps=[...node.data?.bpSteps];
        const models=[...bpSteps[stepIndex].models];
        models[modelIndex]={...models[modelIndex],modelID:e.target.value};
        bpSteps[stepIndex]={...bpSteps[stepIndex],models:models}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const onSideChange=(value)=>{
        const bpSteps=[...node.data?.bpSteps];
        const models=[...bpSteps[stepIndex].models];
        models[modelIndex]={...models[modelIndex],side:value};
        bpSteps[stepIndex]={...bpSteps[stepIndex],models:models}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const onAjustChange=(value)=>{
        const bpSteps=[...node.data?.bpSteps];
        const models=[...bpSteps[stepIndex].models];
        models[modelIndex]={...models[modelIndex],ajust:value};
        bpSteps[stepIndex]={...bpSteps[stepIndex],models:models}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const onZeroStopChange=(value)=>{
        const bpSteps=[...node.data?.bpSteps];
        const models=[...bpSteps[stepIndex].models];
        models[modelIndex]={...models[modelIndex],zeroStop:value};
        bpSteps[stepIndex]={...bpSteps[stepIndex],models:models}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const onDelModel=()=>{
        const bpSteps=[...node.data?.bpSteps];
        const models=[...bpSteps[stepIndex].models];
        delete models[modelIndex]
        bpSteps[stepIndex]={...bpSteps[stepIndex],models:models.filter(item=>item)}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const stepItem=node.data.bpSteps[stepIndex];
    const modelItem=stepItem.models[modelIndex];
    const showSteps=node.data?.__showSteps===false?false:true;
    const showStep=stepItem.__showStep;
    const showModels=stepItem.__showModels;
    const showModel=modelItem.__showModel;

    return (
      <>  
        <Row className="param-panel-row" style={{display:showSteps&&showStep&&showModels?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>
                <div className='button' onClick={setShowModel}>
                {showModel?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Model {modelIndex}</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onDelModel} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSteps&&showStep&&showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Model ID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={modelItem.modelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSteps&&showStep&&showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Side</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={modelItem.side} size='small' onChange={onSideChange}>
                    <Option key='left'>Left</Option>
                    <Option key='right'>Right</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSteps&&showStep&&showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Is Ajust</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={modelItem.ajust} size='small' onChange={onAjustChange}>
                    <Option key='0'>no</Option>
                    <Option key='1'>yes</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSteps&&showStep&&showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Zero Stop</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={modelItem.zeroStop} size='small' onChange={onZeroStopChange}>
                    <Option key='0'>no</Option>
                    <Option key='1'>yes</Option>
                </Select>
            </Col>
        </Row>
      </>);
}