import {Row,Col,Button,Input} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { useCallback } from 'react';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';
import NumericGroupModelItem from './NumericGroupModelItem';

export default function NumericGroupNodeParams({node}){
    const dispatch=useDispatch();

    const setShowModels=useCallback(
      ()=>{
        const showModels=node.data?.__showModels===false?true:false;
        dispatch(updateNodeData({...node.data,__showModels:showModels}));
      },
      [dispatch,node]
    );

    const onAddModel=useCallback(
      ()=>{
        const models=node.data?.models?[...node.data?.models]:[];
        models.push({modelID:"",field:"",__showModel:true});
        dispatch(updateNodeData({...node.data,models:models}));
      },
      [dispatch,node]
    );

    const onToleranceChange=(value)=>{
      dispatch(updateNodeData({...node.data,tolerance:value}));
    }

    const showModels=node.data?.__showModels===false?false:true;

    const models=node.data?.models?.map((item,index)=>{
      console.log("models index :",index,item);
      return (<NumericGroupModelItem key={index} node={node} modelIndex={index}/>)
    });

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label" span={10}>Tolerance</Col>
            <Col className="param-panel-row-input" span={14}>
                <Input value={node.data?.tolerance} onChange={(e)=>onToleranceChange(e.target.value)}/>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" span={10}>
            <div className='button' onClick={setShowModels}>
              {showModels?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Models</span>
          </Col>
          <Col className="param-panel-row-input" span={14}>
            <Button onClick={onAddModel} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {models}
      </>
    );
}

