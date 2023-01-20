import { useParams } from "react-router-dom";

import Header from "./Header";
import FlowTab from "./FlowTab";
import Flow from "./Flow";
import useFrame from "../../hook/useFrame";

import './index.css';

export default function Flows(){
    const sendMessageToParent=useFrame();
    const {appID}=useParams();

    return (
    <div>
        <Header sendMessageToParent={sendMessageToParent} appID={appID} />
        <FlowTab />
        <Flow appID={appID}/>
    </div>);
}