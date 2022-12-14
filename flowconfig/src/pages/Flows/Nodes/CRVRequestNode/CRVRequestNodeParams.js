import {Row,Col,Button} from 'antd';
import { AlignCenterOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';

import {openDialog} from '../../../../redux/dialogSlice';

export default function CRVRequestNodeParams({node,labelWidth}){
    const dispatch=useDispatch();

    const onSetTestData=()=>{
      dispatch(openDialog({type:'testData',title:'Test Data',param:{node:node}}));
    }
    return (
      <>
        <Row className="param-panel-row" gutter={24}>
          <Col className="param-panel-row-label" style={{width:labelWidth}}>
            <span>Test Data</span>
          </Col>
          <Col className="param-panel-row-input" style={{width:'calc(100% - '+labelWidth+'px)'}}>
            <Button onClick={onSetTestData} className='button' size='small' icon={<AlignCenterOutlined />} />
          </Col>
        </Row>
      </>
    );
}

