import {Row,Col,Input,Button} from 'antd';
import { useDispatch } from 'react-redux';
import { PlusSquareOutlined,PlusOutlined,MinusSquareOutlined } from '@ant-design/icons';

import { updateNodeData } from '../../../../redux/flowSlice';
import VerifyValueItem from './VerifyValueItem';

export default function VerifyValueNodeParams({node}){
    const dispatch=useDispatch();

    const setShowItems=()=>{
      const showItems=node.data?.__showItems===false?true:false;
      dispatch(updateNodeData({...node.data,__showItems:showItems}));
    };

    const onAddItem=()=>{
        const items=node.data?.items?[...node.data?.items]:[];
        items.push({verifyID:"",tolerance:"0",aggregation:"",operator:"",value:"",modelID:"",field:"",__showItem:true});
        dispatch(updateNodeData({...node.data,items:items}));
    };

    const showItems=node.data?.__showItems===false?false:true;

    const items=node.data?.items?.map((item,index)=>{
      console.log("models index :",index,item);
      return (<VerifyValueItem key={index} node={node} itemIndex={index}/>)
    });

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" span={10}>
            <div className='button' onClick={setShowItems}>
              {showItems?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Models</span>
          </Col>
          <Col className="param-panel-row-input" span={14}>
            <Button onClick={onAddItem} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {items}
      </>
    );
}

