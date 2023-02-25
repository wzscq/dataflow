import {Row,Col,Button,Input} from 'antd';
import { PlusOutlined,PlusSquareOutlined,MinusSquareOutlined } from '@ant-design/icons';
import { useCallback } from 'react';
import { useDispatch } from 'react-redux';

import { updateNodeData } from '../../../../redux/flowSlice';
import ExportExcelModelItem from './ExportExcelModelItem';

export default function RequestQueryNodeParams({node,labelWidth}){
    const dispatch=useDispatch();

    const setShowSheets=useCallback(
      ()=>{
        const showSheets=node.data?.__showSheets===false?true:false;
        dispatch(updateNodeData({...node.data,__showSheets:showSheets}));
      },
      [dispatch,node]
    );

    const onFileNameChange=(e)=>{
      dispatch(updateNodeData({...node.data,fileName:e.target.value}));
    }

    const onAddSheet=useCallback(
      ()=>{
        const sheets=node.data?.sheets?[...node.data?.sheets]:[];
        sheets.push({modelID:"",fields:[],__showSheet:true,__showFileds:true});
        dispatch(updateNodeData({...node.data,sheets:sheets}));
      },
      [dispatch,node]
    );

    const showSheets=node.data?.__showSheets===false?false:true;

    const sheets=node.data?.sheets?.map((item,index)=>{
      return (<ExportExcelModelItem labelWidth={labelWidth} key={index} node={node} sheetIndex={index}/>)
    });

    return (
      <>
        <Row className="param-panel-row" style={{display:"flex"}} gutter={24}>
            <Col className="param-panel-row-label " style={{width:labelWidth}}>File Name</Col>
            <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
                <Input value={node.data.fileName} onChange={onFileNameChange}/>
            </Col>
        </Row>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <div className='button' onClick={setShowSheets}>
              {showSheets?<MinusSquareOutlined />:<PlusSquareOutlined />}
            </div>
            <span>Sheets</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Button onClick={onAddSheet} className='button' size='small' icon={<PlusOutlined />} />
          </Col>
        </Row>
        {sheets}
      </>
    );
}

