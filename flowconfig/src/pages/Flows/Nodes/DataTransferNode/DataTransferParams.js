import {Row,Col,Button,Input} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { useCallback } from 'react';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';
import TransferSourceModelItem from './TransferSourceModelItem';

export default function DataTransferParams({node,labelWidth}){
    const dispatch=useDispatch();

    const setShowModels=useCallback(
      ()=>{
        const showModels=node.data?.__showSourceModels===false?true:false;
        dispatch(updateNodeData({...node.data,__showSourceModels:showModels}));
      },
      [dispatch,node]
    );

    const onAddModel=useCallback(
      ()=>{
        const models=node.data?.sourceModels?[...node.data?.sourceModels]:[];
        models.push({modelID:"",fields:[],__showSourceModel:true,__showFileds:true,__showUpdateFields:true});
        dispatch(updateNodeData({...node.data,sourceModels:models}));
      },
      [dispatch,node]
    );

    const onModelIDChange=(e)=>{
        dispatch(updateNodeData({...node.data,targetModelID:e.target.value}));
    }

    const onBatchNumberFieldChange=(e)=>{
      dispatch(updateNodeData({...node.data,batchNumberField:e.target.value}));
    }

    const showSourceModels=node.data?.__showSourceModels===false?false:true;

    const models=node.data?.sourceModels?.map((item,index)=>{
      console.log("models index :",index,item);
      return (<TransferSourceModelItem key={index} node={node} labelWidth={labelWidth} modelIndex={index}/>)
    });

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label level" style={{width:labelWidth}}>Target Model</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data?.targetModelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label level" style={{width:labelWidth}}>Batch Number Field</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data?.batchNumberField} onChange={onBatchNumberFieldChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <div className='button' onClick={setShowModels}>
              {showSourceModels?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Source Models</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Button onClick={onAddModel} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {models}
      </>
    );
}

