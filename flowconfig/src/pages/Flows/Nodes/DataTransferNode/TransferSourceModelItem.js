import { useDispatch } from "react-redux";
import {Row,Col,Button,Input,Select} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';

const {Option}=Select;

export default function TransferSourceModelItem({node,modelIndex,labelWidth}){
    const dispatch=useDispatch();

    const setShowModel=()=>{
        const models=[...node.data?.sourceModels];
        models[modelIndex]={...models[modelIndex],__showModel:!(models[modelIndex].__showModel)}
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const setShowFields=()=>{
        const models=[...node.data?.sourceModels];
        models[modelIndex]={...models[modelIndex],__showFields:!(models[modelIndex].__showFields)}
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }
    
    const setShowUpdateFields=()=>{
        const models=[...node.data?.sourceModels];
        models[modelIndex]={...models[modelIndex],__showUpdateFields:!(models[modelIndex].__showUpdateFields)}
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const onModelIDChange=(e)=>{
        const models=[...node.data?.sourceModels];
        models[modelIndex]={...models[modelIndex],modelID:e.target.value}
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const onBatchNumberFieldChange=(e)=>{
        const models=[...node.data?.sourceModels];
        models[modelIndex]={...models[modelIndex],batchNumberField:e.target.value}
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const onDelModel=()=>{
        const models=[...node.data?.sourceModels];
        delete models[modelIndex];
        dispatch(updateNodeData({...node.data,sourceModels:models.filter(item=>item)}));
    }

    const onAddField=()=>{
        const models=[...node.data?.sourceModels];
        const modelFields=models[modelIndex].fields?[...(models[modelIndex].fields)]:[];
        modelFields.push({sourceField:"",targetField:"",keyField:'0',__showField:true});
        models[modelIndex]={...models[modelIndex],fields:modelFields};
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const onAddUpdateField=()=>{
        const models=[...node.data?.sourceModels];
        const modelUpdateFields=models[modelIndex].updateFields?[...(models[modelIndex].updateFields)]:[];
        modelUpdateFields.push({field:"",value:"",__showUpdateField:true});
        models[modelIndex]={...models[modelIndex],updateFields:modelUpdateFields};
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const onDelField=(index)=>{
        const models=[...node.data?.sourceModels];
        const modelFields=models[modelIndex].fields?[...(models[modelIndex].fields)]:[];
        delete modelFields[index];
        models[modelIndex]={...models[modelIndex],fields:modelFields.filter(item=>item)};
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const onDelUpdateField=(index)=>{
        const models=[...node.data?.sourceModels];
        const modelFields=models[modelIndex].updateFields?[...(models[modelIndex].updateFields)]:[];
        delete modelFields[index];
        models[modelIndex]={...models[modelIndex],updateFields:modelFields.filter(item=>item)};
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }
    
    const onFieldChange=(index,value)=>{
        const models=[...node.data?.sourceModels];
        const modelFields=[...(models[modelIndex].fields)];
        modelFields[index]=value;
        models[modelIndex]={...models[modelIndex],fields:modelFields};
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const onUpdateFieldChange=(index,value)=>{
        const models=[...node.data?.sourceModels];
        const modelFields=[...(models[modelIndex].updateFields)];
        modelFields[index]=value;
        models[modelIndex]={...models[modelIndex],updateFields:modelFields};
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const setShowField=(index,showField)=>{
        const models=[...node.data?.sourceModels];
        const modelFields=[...(models[modelIndex].fields)];
        modelFields[index]={...modelFields[index],__showField:showField};
        models[modelIndex]={...models[modelIndex],fields:modelFields};
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const setShowUpdateField=(index,showUpdateField)=>{
        const models=[...node.data?.sourceModels];
        const modelFields=[...(models[modelIndex].updateFields)];
        modelFields[index]={...modelFields[index],__showUpdateFeild:showUpdateField};
        models[modelIndex]={...models[modelIndex],updateFields:modelFields};
        dispatch(updateNodeData({...node.data,sourceModels:models}));
    }

    const modelItem=node.data.sourceModels[modelIndex];
    const showSourceModels=node.data?.__showSourceModels===false?false:true;
    const showModel=modelItem.__showModel;
    const showFields=modelItem.__showFields;
    const showUpdateFields=modelItem.__showUpdateFields;

    const fields=modelItem.fields?.map((item,index)=>{
        const showField=item.__showField
        return (
            <>
            <Row className="param-panel-row" style={{display:showSourceModels&&showModel&&showFields?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>
                    <div className='button' onClick={()=>setShowField(index,!showField)}>
                        {showField?<MinusSquareOutlined />:<PlusSquareOutlined />}
                    </div>
                    <span>Field {index}</span>
                </Col>
                <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <span>{item.targetField}</span>
                    <Button className="button"  onClick={()=>onDelField(index)} size='small' icon={<MinusOutlined />} />
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showSourceModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Source Field</Col>
                <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Input value={item.sourceField} onChange={(e)=>onFieldChange(index,{...item,sourceField:e.target.value})}/>
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showSourceModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Target Field</Col>
                <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Input value={item.targetField} onChange={(e)=>onFieldChange(index,{...item,targetField:e.target.value})}/>
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showSourceModels&&showModel&&showFields&&showField?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Key Field</Col>
                <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Select value={item.keyField} size='small' onChange={(value)=>onFieldChange(index,{...item,keyField:value})}>
                        <Option key='0'>No</Option>
                        <Option key='1'>Yes</Option>
                    </Select>
                </Col>
            </Row>
            </>
        );
    });

    const updateFields=modelItem.updateFields?.map((item,index)=>{
        const showUPdateField=item.__showUpdateField;
        return (
            <>
                <Row className="param-panel-row" style={{display:showSourceModels&&showModel&&showUpdateFields?"flex":"none"}} gutter={24}>
                    <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>
                        <div className='button' onClick={()=>setShowUpdateField(index,!showUPdateField)}>
                            {showUPdateField?<MinusSquareOutlined />:<PlusSquareOutlined />}
                        </div>
                        <span>Field {index}</span>
                    </Col>
                    <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <span>{item.targetField}</span>
                        <Button className="button"  onClick={()=>onDelUpdateField(index)} size='small' icon={<MinusOutlined />} />
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showSourceModels&&showModel&&showUpdateFields&&showUPdateField?"flex":"none"}} gutter={24}>
                    <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Field</Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Input value={item.field} onChange={(e)=>onUpdateFieldChange(index,{...item,field:e.target.value})}/>
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showSourceModels&&showModel&&showUpdateFields&&showUPdateField?"flex":"none"}} gutter={24}>
                    <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Value</Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Input value={item.value} onChange={(e)=>onUpdateFieldChange(index,{...item,value:e.target.value})}/>
                    </Col>
                </Row>
            </>
        )
    })

    const modelItemControl=(
    <>
        <Row className="param-panel-row" style={{display:showSourceModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Model ID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={modelItem.modelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSourceModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Batch Number Field</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={modelItem.batchNumberField} onChange={onBatchNumberFieldChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showSourceModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
                <div className='button' onClick={setShowUpdateFields}>
                    {showUpdateFields?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Update Fields</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onAddUpdateField} size='small' icon={<PlusOutlined />} />
            </Col>
        </Row>
        {updateFields}
        <Row className="param-panel-row" style={{display:showSourceModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
                <div className='button' onClick={setShowFields}>
                    {showFields?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Fields</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onAddField} size='small' icon={<PlusOutlined />} />
            </Col>
        </Row>
        {fields}
    </>);
    
    return (
      <>  
        <Row className="param-panel-row" style={{display:showSourceModels?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>
                <div className='button' onClick={setShowModel}>
                {showModel?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Model {modelIndex}</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onDelModel} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        {modelItemControl}
      </>);
}