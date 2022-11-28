import {Row,Col,Input} from 'antd';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

export default function DelayNodeParams({node}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" span={10}>
            <span>Seconds</span>
          </Col>
          <Col className="param-panel-row-input" span={14}>
            <Input value={node.data?.seconds} onChange={(e)=>onNodeDataChange({...node.data,seconds:e.target.value})}/>
          </Col>
        </Row>
      </>
    );
}

