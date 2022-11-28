import React, { memo } from 'react';
import { Handle } from 'react-flow-renderer';

import './node.css';

export default memo(({ data, isConnectable,selected }) => {
  return (
    <>
      {data.label}
      <Handle
        type="source"
        position="right"
        id="a"
        style={{right:-7,width:10,height:10, top: 14, background: '#555' }}
        isConnectable={isConnectable}
      />
    </>
  );
});