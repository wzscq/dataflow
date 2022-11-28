import {Row,Col,Button,Input} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined,MinusOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

export default function VerifyMatchNodeParams({node}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    const onAddRule=()=>{
        const rules=node.data.rules?[...(node.data.rules)]:[];
        rules.push({modelID:"",message:"",__showRule:true});
        dispatch(updateNodeData({...node.data,rules:rules}));
    }

    const onNodeRuleChange=(index,rule)=>{
        const rules=[...(node.data.rules)];
        rules[index]=rule;
        dispatch(updateNodeData({...node.data,rules:rules}));
    }

    const onDelRule=(index)=>{
        const rules=[...(node.data.rules)];
        delete rules[index];
        dispatch(updateNodeData({...node.data,rules:rules.filter(item=>item)}));
    }

    const showRules=node.data.__showRules===false?false:true;

    const rules=node.data.rules?.map((item,index)=>{
        return (
            <>
                <Row className="param-panel-row" style={{display:showRules?"flex":"none"}} gutter={24}>
                    <Col className="param-panel-row-label level-1" span={10}>
                        <div className='button' onClick={()=>onNodeRuleChange(index,{...item,__showRule:!item.__showRule})}>
                            {item.__showRule?<MinusSquareOutlined />:<PlusSquareOutlined />}
                        </div>
                        <span>Rule {index}</span>
                    </Col>
                    <Col className="param-panel-row-input" span={14}>
                        <Button className="button"  onClick={()=>onDelRule(index)} size='small' icon={<MinusOutlined />} />
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showRules&&item.__showRule?"flex":"none"}}  gutter={24}>
                    <Col className="param-panel-row-label level-2" span={10}>Model ID</Col>
                    <Col className="param-panel-row-input" span={14}>
                        <Input value={item.modelID} onChange={(e)=>onNodeRuleChange(index,{...item,modelID:e.target.value})}/>
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showRules&&item.__showRule?"flex":"none"}}  gutter={24}>
                    <Col className="param-panel-row-label level-2" span={10}>Message</Col>
                    <Col className="param-panel-row-input" span={14}>
                        <Input value={item.message} onChange={(e)=>onNodeRuleChange(index,{...item,message:e.target.value})}/>
                    </Col>
                </Row>
            </>
        )
    })

    return (
      <>
        <Row className="param-panel-row"  gutter={24}>
            <Col className="param-panel-row-label" span={10}>verifyID</Col>
            <Col className="param-panel-row-input" span={14}>
                <Input value={node.data.verifyID} onChange={(e)=>onNodeDataChange({...node.data,verifyID:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24}>
            <Col className="param-panel-row-label" span={10}>Failure Result</Col>
            <Col className="param-panel-row-input" span={14}>
                <Input value={node.data.failureResult} onChange={(e)=>onNodeDataChange({...node.data,failureResult:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24}>
            <Col className="param-panel-row-label" span={10}>Successful Result</Col>
            <Col className="param-panel-row-input" span={14}>
                <Input value={node.data.successfulResult} onChange={(e)=>onNodeDataChange({...node.data,successfulResult:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" span={10}>
            <div className='button' onClick={(e)=>onNodeDataChange({...node.data,__showRules:!showRules})}>
              {showRules?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Rules</span>
          </Col>
          <Col className="param-panel-row-input" span={14}>
            <Button onClick={onAddRule} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {rules}
      </>
    );
}

