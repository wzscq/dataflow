import { useDispatch } from "react-redux";
import {Row,Col,Button,Input,Select} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusOutlined,MinusSquareOutlined,AlignCenterOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';
import {openDialog} from '../../../../redux/dialogSlice';

const { Option } = Select;

export default function QueryModelItem({node,modelIndex,labelWidth}){
    const dispatch=useDispatch();

    const setShowModel=()=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],__showModel:!(models[modelIndex].__showModel)}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const setShowFields=()=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],__showFields:!(models[modelIndex].__showFields)}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const setShowSorter=()=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],__showSorter:!(models[modelIndex].__showSorter)}
        dispatch(updateNodeData({...node.data,models:models}));
    }
    
    const onModelIDChange=(e)=>{
        const models=[...node.data?.models];
        models[modelIndex]={...models[modelIndex],modelID:e.target.value}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onDelModel=()=>{
        const models=[...node.data?.models];
        delete models[modelIndex];
        dispatch(updateNodeData({...node.data,models:models.filter(item=>item)}));
    }

    const onAddField=()=>{
        const models=[...node.data?.models];
        const modelFields=models[modelIndex].fields?[...(models[modelIndex].fields)]:[];
        modelFields.push({field:""});
        models[modelIndex]={...models[modelIndex],fields:modelFields};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onAddSorter=()=>{
        const models=[...node.data?.models];
        const modelSorter=models[modelIndex].sorter?[...(models[modelIndex].sorter)]:[];
        modelSorter.push({field:"",order:"",__showSorter:true});
        models[modelIndex]={...models[modelIndex],sorter:modelSorter};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onDelField=(index)=>{
        const models=[...node.data?.models];
        const modelFields=models[modelIndex].fields?[...(models[modelIndex].fields)]:[];
        delete modelFields[index];
        models[modelIndex]={...models[modelIndex],fields:modelFields.filter(item=>item)};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onDelSorter=(index)=>{
        const models=[...node.data?.models];
        const modelSorter=models[modelIndex].sorter?[...(models[modelIndex].sorter)]:[];
        delete modelSorter[index];
        models[modelIndex]={...models[modelIndex],sorter:modelSorter.filter(item=>item)};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onFilterChange=(e)=>{
        const models=[...node.data?.models];
        var obj = JSON.parse(e.target.value);
        models[modelIndex]={...models[modelIndex],filter:obj}
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onFieldChange=(index,value)=>{
        const models=[...node.data?.models];
        const modelFields=[...(models[modelIndex].fields)];
        modelFields[index]={...modelFields[index],field:value};
        models[modelIndex]={...models[modelIndex],fields:modelFields};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onSorterChange=(index,sorter)=>{
        const models=[...node.data?.models];
        const modelSorter=[...(models[modelIndex].sorter)];
        modelSorter[index]=sorter;
        models[modelIndex]={...models[modelIndex],sorter:modelSorter};
        dispatch(updateNodeData({...node.data,models:models}));
    }

    const onSetFilter=()=>{
        dispatch(openDialog({type:'queryFilter',title:'Query Filter',param:{node:node,modelIndex:modelIndex}}));
    }

    const modelItem=node.data.models[modelIndex];
    const showModels=node.data?.__showModels===false?false:true;
    const showModel=modelItem.__showModel;
    const showFields=modelItem.__showFields;
    const showSorter=modelItem.__showSorter;

    const fields=modelItem.fields?.map((item,index)=>{
        return (
        <Row className="param-panel-row" style={{display:showModels&&showModel&&showFields?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>Field {index}</Col>
            <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={item.field} onChange={(e)=>onFieldChange(index,e.target.value)}/>
                <Button className="button"  onClick={()=>onDelField(index)} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        );
    });

    const sorter=modelItem.sorter?.map((item,index)=>{
        const {__showSorter,field,order}=item;
        return (<>
            <Row className="param-panel-row" style={{display:showModels&&showModel&&showSorter?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-3" style={{width:labelWidth}}>
                    <div className='button' onClick={()=>onSorterChange(index,{...item,__showSorter:!__showSorter})}>
                        {__showSorter?<MinusSquareOutlined />:<PlusSquareOutlined />}
                    </div>
                    <span>Sorter {index}</span>
                </Col>
                <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Button className="button"  onClick={()=>onDelSorter(index)} size='small' icon={<MinusOutlined />} />
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showModels&&showModel&&showSorter&&__showSorter?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Field</Col>
                <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Input value={field} onChange={(e)=>onSorterChange(index,{...item,field:e.target.value})}/>
                </Col>
            </Row>
            <Row className="param-panel-row" style={{display:showModels&&showModel&&showSorter&&__showSorter?"flex":"none"}} gutter={24}>
                <Col className="param-panel-row-label level-4" style={{width:labelWidth}}>Order</Col>
                <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                    <Select value={order} size='small' onChange={(value)=>onSorterChange(index,{...item,order:value})}>
                        <Option key='asc'>ASC</Option>
                        <Option key='desc'>DESC</Option>
                    </Select>
                </Col>
            </Row>
        </>
        );
    });

    const modelItemControl=(
    <>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Model ID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={modelItem.modelID} onChange={onModelIDChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
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
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Filter</Col>
            <Col className="param-panel-row-inputwithbutton" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input disabled={true} value={JSON.stringify(modelItem.filter)} onChange={onFilterChange}/>
                <Button className="button"  onClick={onSetFilter} size='small' icon={<AlignCenterOutlined />} />
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showModels&&showModel?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>
                <div className='button' onClick={setShowSorter}>
                    {showSorter?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Sorter</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onAddSorter} size='small' icon={<PlusOutlined />} />
            </Col>
        </Row>
        {sorter}
    </>);
    
    return (
      <>  
        <Row className="param-panel-row" style={{display:showModels?"flex":"none"}} gutter={24}>
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