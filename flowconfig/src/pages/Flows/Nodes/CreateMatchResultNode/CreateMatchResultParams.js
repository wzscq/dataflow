import {Row,Col,Button,Input,Select} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined,MinusOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

const { Option } = Select;

export default function CreateMatchResultParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    const onAddGroupField=()=>{
        const fields=node.data.groupModel?.fields?[...node.data.groupModel.fields]:[];
        fields.push({field:"",sourceModel:"",sourceField:"",aggregation:"",__showField:true});
        dispatch(updateNodeData({...node.data,groupModel:{...node.data.groupModel,fields:fields}}));
    }

    const onGroupFieldChange=(index,fieldItem)=>{
        const fields=[...node.data.groupModel.fields];
        fields[index]=fieldItem;
        dispatch(updateNodeData({...node.data,groupModel:{...node.data.groupModel,fields:fields}}));
    }

    const onDelGroupField=(index)=>{
        const fields=[...node.data.groupModel.fields];
        delete fields[index];
        dispatch(updateNodeData({...node.data,groupModel:{...node.data.groupModel,fields:fields.filter(item=>item)}}));
    }

    const showGroupModel=node.data.__showGroupModel===false?false:true;
    const showGroupModelFields=node.data.__showGroupModelFields===false?false:true;

    const fields=node.data.groupModel?.fields?.map((item,index)=>{
        return (
            <>
                <Row className="param-panel-row" style={{display:showGroupModelFields?"flex":"none"}} gutter={24}>
                    <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
                        <div className='button' onClick={()=>onGroupFieldChange(index,{...item,__showField:!item.__showField})}>
                            {item.__showField?<MinusSquareOutlined />:<PlusSquareOutlined />}
                        </div>
                        <span>Field {index}</span>
                    </Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Button className="button"  onClick={()=>onDelGroupField(index)} size='small' icon={<MinusOutlined />} />
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showGroupModelFields&&item.__showField?"flex":"none"}}  gutter={24}>
                    <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>Field</Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Input value={item.field} onChange={(e)=>onGroupFieldChange(index,{...item,field:e.target.value})}/>
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showGroupModelFields&&item.__showField?"flex":"none"}}  gutter={24}>
                    <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>Source Model</Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Input value={item.sourceModel} onChange={(e)=>onGroupFieldChange(index,{...item,sourceModel:e.target.value})}/>
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showGroupModelFields&&item.__showField?"flex":"none"}}  gutter={24}>
                    <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>Source Field</Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Input value={item.sourceField} onChange={(e)=>onGroupFieldChange(index,{...item,sourceField:e.target.value})}/>
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showGroupModelFields&&item.__showField?"flex":"none"}} gutter={24}>
                    <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>Aggregation</Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Select value={item.aggregation} size='small' onChange={(value)=>onGroupFieldChange(index,{...item,aggregation:value})}>
                            <Option key='first'>first</Option>
                            <Option key='sum'>sum</Option>
                        </Select>
                    </Col>
                </Row>
            </>
        )
    })

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label" style={{width:labelWidth}}>Match Result</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={node.data.matchResult} size='small' onChange={(value)=>onNodeDataChange({...node.data,matchResult:value})}>
                    <Option key='0'>Completely</Option>
                    <Option key='1'>Partially</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <div className='button' onClick={()=>onNodeDataChange({...node.data,__showGroupModel:!showGroupModel})}>
              {showGroupModel?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Group Model</span>
          </Col>
        </Row>
        <Row className="param-panel-row"  gutter={24} style={{display:showGroupModel?"flex":"none"}}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>Model ID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.groupModel?.modelID} onChange={(e)=>onNodeDataChange({...node.data,groupModel:{...node.data.groupModel,modelID:e.target.value}})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24} style={{display:showGroupModel?"flex":"none"}}>
          <Col className="param-panel-row-label  level-1" style={{width:labelWidth}}>
            <div className='button' onClick={()=>onNodeDataChange({...node.data,__showGroupModelFields:!showGroupModelFields})}>
              {showGroupModelFields?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Fields</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Button onClick={onAddGroupField} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {fields}
      </>
    );
}

