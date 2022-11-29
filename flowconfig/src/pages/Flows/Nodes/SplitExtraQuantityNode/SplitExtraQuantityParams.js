import {Row,Col,Button,Input,Select} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined,MinusOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

const { Option } = Select;

export default function SplitExtraQuantityParams({node,labelWidth}){
    const dispatch=useDispatch();

    const setShowModels=()=>{
        const showModels=node.data?.__showModels===false?true:false;
        dispatch(updateNodeData({...node.data,__showModels:showModels}));
      }

    const onAddMatchModel=()=>{
      const models=node.data?.models?[...node.data?.models]:[];
      models.push({side:"left",modelID:"",field:"",__showModel:true});
      dispatch(updateNodeData({...node.data,models:models}));
    }

    const onDelModel=(index)=>{
      const models=[...node.data?.models];
      delete models[index];
      dispatch(updateNodeData({...node.data,models:models.filter(item=>item)}));
    }
   
    const onModelChange=(index,value)=>{
        const models=[...node.data?.models];
        models[index]=value;
        dispatch(updateNodeData({...node.data,models:models}))
    }

    const showModels=node.data?.__showModels===false?false:true;

    const models=node.data.models?.map((item,index)=>{
      return (<>
        <Row className="param-panel-row" style={{display:showModels?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>
                <div className='button' onClick={()=>onModelChange(index,{...item,__showModel:!item.__showModel})}>
                    {item.__showModel?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                Model {index}
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={()=>onDelModel(index)} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&item.__showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>ModelID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={item.modelID} onChange={(e)=>onModelChange(index,{...item,modelID:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&item.__showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Field</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={item.field} onChange={(e)=>onModelChange(index,{...item,field:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&item.__showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Side</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={item.side} size='small' onChange={(value)=>onModelChange(index,{...item,side:value})}>
                    <Option key='left'>Left</Option>
                    <Option key='right'>Right</Option>
                </Select>
            </Col>
        </Row>
      </>);
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
            <Button onClick={onAddMatchModel} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {models}
      </>
    );
}

