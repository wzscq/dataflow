import React, { memo } from 'react';
import { Handle } from 'react-flow-renderer';

import '../node.css';

export default memo(({ data, isConnectable,selected }) => {
  return (
    <>
      <Handle
        type="target"
        position="left"
        id="a"
        style={{left:-7, width:10,height:10, top: 14, background: '#555' }}
        isConnectable={isConnectable}
      />
      {data.label}
      <Handle
        type="source"
        position="right"
        style={{right:-7, width:10,height:10, top: 14, background: '#555' }}
        isConnectable={isConnectable}
      />
    </>
  );
});