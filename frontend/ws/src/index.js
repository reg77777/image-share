import React from 'react';
import ReactDOM from 'react-dom/client';
import DropZone from "./DropZone";
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';

class App extends React.Component{
    render(){
      return(
        <div className="App">
          <header className="App-header">
            <AppBar position="static">
                <Toolbar>
                    <Typography>
                        APP
                    </Typography>
                </Toolbar>
            </AppBar>
          </header>
          <DropZone/>
        </div>
      );
    }
}

const root=ReactDOM.createRoot(document.getElementById("root"));
root.render(<App/>);
