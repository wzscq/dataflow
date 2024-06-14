import {Row,Col,Button,Select} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { useCallback } from 'react';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';
import TransfromModelItem from './TransformModelItem';

const {Option}=Select;

export default function DataTransformParams({node,labelWidth}){
    const dispatch=useDispatch();

    const setShowModels=useCallback(
      ()=>{
        const showModels=node.data?.__showModels===false?true:false;
        dispatch(updateNodeData({...node.data,__showModels:showModels}));
      },
      [dispatch,node]
    );

    const setOutputType=(value)=>{
        dispatch(updateNodeData({...node.data,outputType:value}));
    }

    const onAddModel=useCallback(
      ()=>{
        const models=node.data?.models?[...node.data?.models]:[];
        models.push({modelID:"",fields:[],__showModel:true,__showFileds:true});
        dispatch(updateNodeData({...node.data,models:models}));
      },
      [dispatch,node]
    );

    const showModels=node.data?.__showModels===false?false:true;

    const models=node.data?.models?.map((item,index)=>{
      console.log("models index :",index,item);
      return (<TransfromModelItem key={index} node={node} modelIndex={index} labelWidth={labelWidth}/>)
    });

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>Output</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={node.data.outputType} size='small' onChange={setOutputType}>
                    <Option key='all'>All</Option>
                    <Option key='modified'>Modified</Option>
                    <Option key='original'>Original</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
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
      </>
    );
}

