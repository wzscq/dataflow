import Header from "./Header";
import FlowTab from "./FlowTab";
import Flow from "./Flow";
import './index.css';
import { useParams } from "react-router-dom";

export default function Flows(){
    const {appID}=useParams();

    return (
    <div>
        <Header appID={appID} />
        <FlowTab />
        <Flow appID={appID}/>
    </div>);
}