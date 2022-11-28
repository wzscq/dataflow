import {HashRouter,Routes,Route} from "react-router-dom";
import Flow from './pages/Flows';
import Dialog from './dialogs';

function App() {
  return (
    <>
      <HashRouter>
        <Routes>
          <Route path="/:appID" exact={true} element={<Flow/>} />
        </Routes>
      </HashRouter>
      <Dialog/>
    </>
  );
}

export default App;
