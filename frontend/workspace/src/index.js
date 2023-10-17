import React, {useEffect} from 'react';
import ReactDOM from 'react-dom/client';
import DropZone from "./DropZone";
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import axios from "axios";
import {Buffer} from "buffer";

class App extends React.Component{
    state={
        images:[],
    };

    constructor(props){
        super(props);
        axios.get(process.env.REACT_APP_GET_URL)
            .then(res=>{
                var images=[]
                res.data.slice(0).reverse().map(data=>{
                    images.push({Category:data.Category,Data:"data:image/png;base64,"+data.Data})
                })
                console.log(images)
                this.setState({images:images})
            })
            .catch(err=>{
                console.log("get error")
            })
    };

    get=()=>{
        axios.get(process.env.REACT_APP_GET_URL)
            .then(res=>{
                var images=[]
                res.data.slice(0).reverse().map(data=>{
                    images.push({Category:data.Category,Data:"data:image/png;base64,"+data.Data})
                })
                this.setState({images:images})
            })
    }

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
                <DropZone get={this.get} />
                {this.state.images.map(image => (
                <div key={image.objectID} style={{padding:"30px"}}>
                    <figure style={{
                        display:"flex",
                        flexFlow:"column",
                        padding:"5px",
                        maxWidth:"220px",
                        margin:"auto"
                    }}>
                    <div style={{textAlign:"center"}}>
                        <img src={image.Data} width="500px"/>
                    </div>
                        <figcaption style={{backgroundColor:"#222", color:"#fff", font:"italic smaller sans-serif", padding:"3px", textAlign:"center",}}>{image.Category}</figcaption>
                    </figure>
                </div>
                ))}
            </div>
        );
    }
}

const root=ReactDOM.createRoot(document.getElementById("root"));
root.render(<App/>);
