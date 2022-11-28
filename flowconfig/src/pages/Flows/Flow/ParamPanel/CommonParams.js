import {Row,Col,Input, InputNumber} from 'antd';
import { useCallback } from 'react';
import { useDispatch } from 'react-redux';

import {updateNodeData} from '../../../../redux/flowSlice';

export default function CommonParam({node}){
    const dispatch=useDispatch();

    const onLabelChange=useCallback(
        (e)=>{
            const { value } = e.target;
            dispatch(updateNodeData({...node.data,label:value}));
        },
        [dispatch,node]
    );

    const onPriorityChange=(value)=>{
        dispatch(updateNodeData({...node.data,priority:value}));
    }

    return (
        <>
            <Row className="param-panel-row" gutter={24}>
                <Col className="param-panel-row-title" span={24}>{node.type} Node: {node.id} </Col>
            </Row>
            <Row className="param-panel-row" gutter={24}>
                <Col className="param-panel-row-label" span={10}>Label</Col>
                <Col className="param-panel-row-input" span={14}>
                    <Input value={node.data.label} onChange={onLabelChange}/>
                </Col>
            </Row>
            <Row className="param-panel-row" gutter={24}>
                <Col className="param-panel-row-label" span={10}>Priority</Col>
                <Col className="param-panel-row-input" span={14}>
                    <InputNumber size='small' value={node.data.priority} onChange={onPriorityChange}/>
                </Col>
            </Row>
        </>
    )
}