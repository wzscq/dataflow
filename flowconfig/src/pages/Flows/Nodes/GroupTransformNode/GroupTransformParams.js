import {Row,Col,Button,Input} from 'antd';
import { AlignCenterOutlined } from '@ant-design/icons';
import { useDispatch } from 'react-redux';
import {openDialog} from '../../../../redux/dialogSlice';

const initContent=
`/*
分组数据处理逻辑用于对分组数进行变更修改，这里的输入和输出都是一个数据分组，格式请参照输入参数的说明
这里的代码中仅包含函数体内的部分，函数名称、参数和花括号程序会在运行时自动加上，
输入参数：
groupItem:分组数据，格式如下
    {
      models:[  //这个部分时数据分组中的业务数据，一个模型数据一项
        {
          modelID:"",  //数据模型ID
          list:[       //模型数据数组
            {
              fieldname1:value1,
              fieldname2:value2,
              ...
            },
            ...
          ]   
        },
        ...
      ],
      verifyResult:[  //这个部分表示流程中前序数据校验节点的校验结果，一般不需要处理，只需要将输入数据组中的对应值复制到输出分组中即可
        {
          verfiyID:"",
          verfiyType:"",
          message:"",
	        result:""
        },
        ...
      ]
    }

函数中必须包含return语句用于返回数据，返回语句的一般格式如下：
    return varObject; 
    其中varObject保存返回数据，格式和输入数据一致。

处理函数中可使用以下全局参数
g_BatchNumber:只读参数，处理节点生成的UUID，表示处理批次号，可以用这个变量对模型中的字段赋值来标识同一批次处理的数据。
g_Index:读写参数，处理批次内序号，初始值为1，在同一批次处理的不同分组数据中这个值是共享的，可以实现跨分组数据的编号。

*/`

export default function DataTransformParams({node}){
    const dispatch=useDispatch();

    const onEditScript=()=>{
      dispatch(openDialog({type:'funcScript',title:'Function Script',param:{node:node,initContent}}));
    }

    return (
      <>
        <Row className="param-panel-row" gutter={24}>
            <Col className="param-panel-row-label" span={10}>Function Script</Col>
            <Col className="param-panel-row-inputwithbutton" span={14}>
                <Input disabled={true} value={node.data.funcScript?.name}/>
                <Button className="button"  onClick={()=>{onEditScript()}} size='small' icon={<AlignCenterOutlined />} />
            </Col>
        </Row>
      </>
    );
}