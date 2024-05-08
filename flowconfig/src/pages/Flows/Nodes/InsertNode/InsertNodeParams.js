import {Row,Col,Input} from 'antd';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

export default function InsertNodeParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>SQL MaxLen</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={node.data?.sqlMaxLen} onChange={(e)=>onNodeDataChange({...node.data,sqlMaxLen:e.target.value})}/>
          </Col>
        </Row>
      </>
    );
}

