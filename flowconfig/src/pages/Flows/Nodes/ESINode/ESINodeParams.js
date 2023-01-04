import {Row,Col,Button} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined} from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';
import ESINodeModelItem from './ESINodeModelItem';

export default function ESINodeParams({node,labelWidth}){
    const dispatch=useDispatch();
    const setShowModels=()=>{
      const showModels=node.data?.__showModels===false?true:false;
      dispatch(updateNodeData({...node.data,__showModels:showModels}));
    };

    const onAddModel=()=>{
      const models=node.data?.models?[...node.data?.models]:[];
      models.push({modelID:"",fields:[],__showModel:true,__showFileds:true,sheets:[],__showSheets:true,__showOptions:true});
      dispatch(updateNodeData({...node.data,models:models}));
    };

    const showModels=node.data?.__showModels===false?false:true;

    const models=node.data?.models?.map((item,index)=>{
      return (<ESINodeModelItem  labelWidth={labelWidth} key={index} node={node} modelIndex={index}/>)
    });

    return (
      <>
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

