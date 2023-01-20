import {Row,Col,Input,Select,InputNumber} from 'antd';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

const { Option } = Select;

export default function CRVFormNodeParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>Title</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={node.data?.title} onChange={(e)=>onNodeDataChange({...node.data,title:e.target.value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>Url</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={node.data?.url} onChange={(e)=>onNodeDataChange({...node.data,url:e.target.value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>Location</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Select value={node.data?.location} size='small' onChange={(value)=>onNodeDataChange({...node.data,location:value})}>
                <Option key='modal'>Modal</Option>
                <Option key='tab'>Tab</Option>
            </Select>
          </Col>
        </Row>        
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>Key</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Input value={node.data?.key} onChange={(e)=>onNodeDataChange({...node.data,key:e.target.value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>Width</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <InputNumber size='small' value={node.data?.width} onChange={(value)=>onNodeDataChange({...node.data,width:value})}/>
          </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>Height</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <InputNumber size='small' value={node.data?.height} onChange={(value)=>onNodeDataChange({...node.data,height:value})}/>
          </Col>
        </Row>
      </>
    );
}

