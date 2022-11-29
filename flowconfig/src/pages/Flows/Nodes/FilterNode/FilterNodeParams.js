import {Row,Col,Button,Input} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined,MinusOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';

export default function FilterNodeParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onNodeDataChange=(data)=>{
        dispatch(updateNodeData(data));
    }

    const onAddFilterItem=()=>{
        const filter=node.data.filter?[...(node.data.filter)]:[];
        filter.push({modelID:"",message:"",__showFilter:true});
        dispatch(updateNodeData({...node.data,filter:filter}));
    }

    const onFilterItemChange=(index,filterItem)=>{
        const filter=[...(node.data.filter)];
        filter[index]=filterItem;
        dispatch(updateNodeData({...node.data,filter:filter}));
    }

    const onDelFilterItem=(index)=>{
        const filter=[...(node.data.filter)];
        delete filter[index];
        dispatch(updateNodeData({...node.data,filter:filter.filter(item=>item)}));
    }

    const showFilter=node.data.__showFilter===false?false:true;

    const filterItems=node.data.filter?.map((item,index)=>{
        return (
            <>
                <Row className="param-panel-row" style={{display:showFilter?"flex":"none"}} gutter={24}>
                    <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>
                        <div className='button' onClick={()=>onFilterItemChange(index,{...item,__showFilter:!item.__showFilter})}>
                            {item.__showFilter?<MinusSquareOutlined />:<PlusSquareOutlined />}
                        </div>
                        <span>Item {index}</span>
                    </Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Button className="button"  onClick={()=>onDelFilterItem(index)} size='small' icon={<MinusOutlined />} />
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showFilter&&item.__showFilter?"flex":"none"}}  gutter={24}>
                    <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Verify ID</Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Input value={item.verifyID} onChange={(e)=>onFilterItemChange(index,{...item,verifyID:e.target.value})}/>
                    </Col>
                </Row>
                <Row className="param-panel-row" style={{display:showFilter&&item.__showFilter?"flex":"none"}}  gutter={24}>
                    <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Result</Col>
                    <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                        <Input value={item.result} onChange={(e)=>onFilterItemChange(index,{...item,result:e.target.value})}/>
                    </Col>
                </Row>
            </>
        )
    })

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <div className='button' onClick={(e)=>onNodeDataChange({...node.data,__showFilter:!showFilter})}>
              {showFilter?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Filter</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Button onClick={onAddFilterItem} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {filterItems}
      </>
    );
}

