import {Row,Col,Button,Input,Select} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';
import EBProcessingStepItem from './EBProcessingStepItem';

const { Option } = Select;

export default function EBProcessingNodeParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
      dispatch(updateNodeData(data));
    }

    const onAddStep=()=>{
        const bpSteps=node.data?.bpSteps?[...node.data?.bpSteps]:[];
        bpSteps.push({models:[],writeoffType:"",__showStep:true,__showModels:true});
        dispatch(updateNodeData({...node.data,bpSteps:bpSteps}));
    }

    const showSteps=node.data?.__showSteps===false?false:true;

    const bpSteps=node.data?.bpSteps?.map((item,index)=>{
      console.log("step index :",index,item);
      return (<EBProcessingStepItem key={index} node={node} stepIndex={index} labelWidth={labelWidth}/>)
    });

    return (
      <>
        <Row className="param-panel-row"  gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>Writeoff ModelID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.writeoffDetailModelID} onChange={(e)=>onNodeDataChange({...node.data,writeoffDetailModelID:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>Group ModelID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.matchGroupModelID} onChange={(e)=>onNodeDataChange({...node.data,matchGroupModelID:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>Deal Amount</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={node.data.dealAmount} size='small' onChange={value=>onNodeDataChange({...node.data,dealAmount:value})}>
                    <Option key='0'>no</Option>
                    <Option key='1'>yes</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>Deal Quantity</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={node.data.dealQuantity} size='small' onChange={value=>onNodeDataChange({...node.data,dealQuantity:value})}>
                    <Option key='0'>no</Option>
                    <Option key='1'>yes</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <div className='button' onClick={()=>onNodeDataChange({...node.data,__showSteps:!showSteps})}>
              {showSteps?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Steps</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Button onClick={onAddStep} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {bpSteps}
      </>
    );
}

