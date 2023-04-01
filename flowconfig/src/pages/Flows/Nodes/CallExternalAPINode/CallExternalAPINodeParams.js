import {Row,Col,Input,Select} from 'antd';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

const {Option}=Select;

export default function CallExternalAPINodeParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>URL</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={node.data?.url} onChange={(e)=>onNodeDataChange({...node.data,url:e.target.value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>ModelID</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={node.data?.modelID} onChange={(e)=>onNodeDataChange({...node.data,modelID:e.target.value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>RequestEachRow</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Select value={node.data?.reqForEachRow} size='small' onChange={(value)=>onNodeDataChange({...node.data,reqForEachRow:value})}>
              <Option key='no'>No</Option>
              <Option key='yes'>Yes</Option>
            </Select>
          </Col>
        </Row>
      </>
    );
}

