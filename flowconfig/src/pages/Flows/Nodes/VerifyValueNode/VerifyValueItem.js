import { useDispatch } from "react-redux";
import {Row,Col,Button,Input,Select} from 'antd';
import { PlusSquareOutlined,MinusOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { updateNodeData } from '../../../../redux/flowSlice';

const { Option } = Select;

export default function VerifyValueItem({node,itemIndex,labelWidth}){
    const dispatch=useDispatch();

    const onChangeVerifyItem=(value)=>{
        const items=[...node.data?.items];
        items[itemIndex]=value
        dispatch(updateNodeData({...node.data,items:items}));
    }

    const onDelVerify=()=>{
        const items=[...node.data?.items];
        delete items[itemIndex];
        dispatch(updateNodeData({...node.data,items:items.filter(item=>item)}));
    }

    const verifyItem=node.data.items[itemIndex];
    const showItems=node.data?.__showItems===false?false:true;
    const showItem=verifyItem.__showItem;

    const itemControl=(
    <>
        <>
        <Row className="param-panel-row" style={{display:showItems&&showItem?"flex":"none"}}   gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>VerifyID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={verifyItem.verifyID} onChange={(e)=>onChangeVerifyItem({...verifyItem,verifyID:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showItems&&showItem?"flex":"none"}}  gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Tolerance</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={verifyItem?.tolerance} onChange={(e)=>onChangeVerifyItem({...verifyItem,tolerance:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showItems&&showItem?"flex":"none"}}  gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>ModelID</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={verifyItem?.modelID} onChange={(e)=>onChangeVerifyItem({...verifyItem,modelID:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showItems&&showItem?"flex":"none"}}  gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Field</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={verifyItem?.field} onChange={(e)=>onChangeVerifyItem({...verifyItem,field:e.target.value})}/>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showItems&&showItem?"flex":"none"}}  gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Aggregation</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Select value={verifyItem?.aggregation} size='small' onChange={(value)=>onChangeVerifyItem({...verifyItem,aggregation:value})}>
                  <Option key='sum'>SUM</Option>
                  <Option key='avg'>{'AVG'}</Option>
                  <Option key='count'>{'COUNT'}</Option>
                  <Option key='max'>{'MAX'}</Option>
                  <Option key='min'>{'MIN'}</Option>
                </Select>
            </Col>
        </Row>
        <Row className="param-panel-row" style={{display:showItems&&showItem?"flex":"none"}}  gutter={24}>
            <Col className="param-panel-row-label level-2" style={{width:labelWidth}}>Value</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={verifyItem?.value} onChange={(e)=>onChangeVerifyItem({...verifyItem,value:e.target.value})}/>
            </Col>
        </Row>
      </>
    </>);
    
    return (
      <>  
        <Row className="param-panel-row" style={{display:showItems?"flex":"none"}} gutter={24}>
            <Col className="param-panel-row-label level-1" style={{width:labelWidth}}>
                <div className='button' onClick={()=>onChangeVerifyItem({...verifyItem,__showItem:!showItem})}>
                {showItem?<MinusSquareOutlined />:<PlusSquareOutlined />}
                </div>
                <span>Item {itemIndex}</span>
            </Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Button className="button"  onClick={onDelVerify} size='small' icon={<MinusOutlined />} />
            </Col>
        </Row>
        {itemControl}
      </>);
}