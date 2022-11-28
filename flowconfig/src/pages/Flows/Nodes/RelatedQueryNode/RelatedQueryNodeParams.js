import {Row,Col,Button} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { useCallback } from 'react';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';
import RelatedQueryModelItem from './RelatedQueryModelItem';

export default function RelatedQueryNodeParams({node}){
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
        models.push({modelID:"",fields:[],filter:{},__showModel:true,__showFileds:true,sorter:[],__showSorter:true});
        dispatch(updateNodeData({...node.data,models:models}));
      },
      [dispatch,node]
    );

    const showModels=node.data?.__showModels===false?false:true;

    const models=node.data?.models?.map((item,index)=>{
      console.log("models index :",index,item);
      return (<RelatedQueryModelItem key={index} node={node} modelIndex={index}/>)
    });

    return (
      <>
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

