import {Row,Col,Input} from 'antd';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

export default function FlowNodeParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>ErrorCode</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={node.data?.errorCode} onChange={(e)=>onNodeDataChange({...node.data,errorCode:e.target.value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>Message</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={node.data?.message} onChange={(e)=>onNodeDataChange({...node.data,message:e.target.value})}/>
          </Col>
        </Row>
      </>
    );
}

