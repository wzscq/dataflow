import {useState} from 'react';
import { useSelector,useDispatch } from "react-redux";
import {Space,Button,Upload} from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import {closeDialog} from '../../../../../redux/dialogSlice';
import { updateNodeData } from '../../../../../redux/flowSlice';
import './index.css';

export default function DialogESITestData({dialogIndex}){
  const dispatch = useDispatch();
  const dialogItem=useSelector(state=>state.dialog.dialogs[dialogIndex]);
  const {node}=dialogItem.param;
  console.log(node);
  const initFileList=[];
  if(node.data?.testData?.esiFile?.list.length>0){
    const testFile=node.data.testData.esiFile.list[0];
    initFileList.push({
      uid:testFile.id,
      name:testFile.name,
      contentBase64:testFile.contentBase64
    });
  }
  const [fileList,setFileList]=useState(initFileList);

  const onCancel=()=>{
    dispatch(closeDialog());
  }

  const onOk=()=>{
    let testData={};
    if(fileList.length>0){
      const selectedFile=fileList[0]
      testData={
        esiFile:{
          list:[
            {
              id:selectedFile.uid,
              name:selectedFile.name,
              contentBase64:selectedFile.contentBase64
            }
          ]
        }
      };
    }
    dispatch(updateNodeData({...node.data,testData:testData}));
    dispatch(closeDialog());
  }

  const props = {
    accept:"*.xlsx",
    showUploadList:{
        showDownloadIcon:true,
        showRemoveIcon:true,
    },
    onRemove: file => {
        const index = fileList.indexOf(file);
        const newFileList = fileList.slice();
        newFileList.splice(index, 1);
        setFileList(newFileList);
    },
    beforeUpload: file => {
        const reader = new FileReader();
        reader.onload=(e)=>{
            const fileTmp={uid:file.uid,name:file.name,contentBase64:e.target.result};
            setFileList([...fileList,fileTmp]);
        };
        reader.readAsDataURL(file);
        return false;
    },
    fileList,
  };

  return (
    <div className="dialog-func-script">
      <div style={{maxHeight:350,overflowY:'auto'}}>
        <Upload {...props}>
            {
                (fileList.length<1)?(
                    <Button icon={<UploadOutlined />}>
                      SelectFile
                    </Button>
                ):null
            }
        </Upload>
      </div>
      <div className="dialog-bottom-bar">
          <Space style={{float:'right'}}>
              <Button style={{minWidth:100}} size="small" onClick={onOk}>Ok</Button>
              <Button style={{minWidth:100}} size="small" onClick={onCancel}>Close</Button>
          </Space>
      </div>
    </div>
  );
}