import { useDispatch } from "react-redux";
import {Row,Col,Button,Input,Select} from 'antd';
import { PlusSquareOutlined,MinusOutlined,PlusOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';
import EBProcessingStepModelItem from './EBProcessingStepModelItem';

const { Option } = Select;

export default function EBProcessingStepItem({node,stepIndex,labelWidth}){
    const dispatch=useDispatch();

    const setShowStep=()=>{
        const bpSteps=[...node.data?.bpSteps];
        bpSteps[stepIndex]={...bpSteps[stepIndex],__showStep:!(bpSteps[stepIndex].__showStep)}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const setShowModels=()=>{
        const bpSteps=[...node.data?.bpSteps];
        bpSteps[stepIndex]={...bpSteps[stepIndex],__showModels:!(bpSteps[stepIndex].__showModels)}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }
    
    const onWriteoffTypeChange=(value)=>{
        const bpSteps=[...node.data?.bpSteps];
        bpSteps[stepIndex]={...bpSteps[stepIndex],writeoffType:value}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const onMonoNegativeWriteoffMethodChange=(value)=>{
        const bpSteps=[...node.data?.bpSteps];
        bpSteps[stepIndex]={...bpSteps[stepIndex],monoNegativeWriteoffMethod:value}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const onDelStep=()=>{
        const bpSteps=[...node.data?.bpSteps];
        delete bpSteps[stepIndex];
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps.filter(item=>item)}));
    }

    const onAddModel=()=>{
        const bpSteps=[...node.data?.bpSteps];
        const models=[...bpSteps[stepIndex].models];
        models.push({modelID:"",side:"",ajust:"",zeroStop:"1",__showModel:true});
        bpSteps[stepIndex]={...bpSteps[stepIndex],models:models}
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const stepItem=node.data.bpSteps[stepIndex];
    const showSteps=node.data?.__showSteps===false?false:true;
    const showStep=stepItem.__showStep;
    const showModels=stepItem.__showModels;

    const models=stepItem.models.map((item,index)=>{
        console.log("step index :",index,item);
        return (<EBProcessingStepModelItem labelWidth={labelWidth} modelIndex={index} key={index} node={node} stepIndex={stepIndex}/>)
    });

    const stepItemControl=(
    <>
        <Row className="param-panel-row" style={{display:showSteps&&showStep?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Writeoff Type</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={stepItem.writeoffType} onChange={(e)=>onWriteoffTypeChange(e.target.value)}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSteps&&showStep?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Mono Negative</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={stepItem.monoNegativeWriteoffMethod} size='small' onChange={onMonoNegativeWriteoffMethodChange}>
                    <Option key='0'>Writeoff Negative Value</Option>
                    <Option key='1'>Writeoff Positive Value</Option>
                    <Option key='2'>Writeoff Left Value</Option>
                    <Option key='3'>Writeoff Right Value</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSteps&&showStep?"flex":"none"}} gutter={24}>
          <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
            <div className='button' onClick={setShowModels}>
              {showModels?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Models</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Button onClick={onAddModel} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {models}
    </>);
    
    return (
      <>  
        <Row className="param-panel-row" style={{display:showSteps?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>
                <div className='button' onClick={setShowStep}>
                {showStep?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Step {stepIndex}</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onDelStep} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        {stepItemControl}
      </>);
}