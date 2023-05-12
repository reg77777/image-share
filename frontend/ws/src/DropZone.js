import React from "react";
import Dropzone from "react-dropzone";
import {Button,Container,makeStyles} from "@mui/material";
import axios from "axios";

axios.defaults.baseURL='http://localhost:8081'

class DropZone extends React.Component{
    constructor(props){
        super(props);
    };

    state={
        filenames:[],
        dropzone_color: "blue"
    };

    dropzoneStyleCreator = (color) => {
        return {
            background: color,
            width: "100px",
            height: "100px",
            transitionProperty: "all",
            transitionDuration: "0.3s",
            transitionDelay: "0s",
            transitionTimingFunction: "ease-in",
        }
    };

    onMouseEnter=()=>{
        this.setState(state=>({
            dropzone_color: "red"
        }))
    };

    onMouseLeave=()=>{
        this.setState(state=>({
            dropzone_color: "blue"
        }))
    };

    onDrop=(files)=>{
        axios.get("http://localhost:8081")
            .then(res=>{
                console.log(res)
            })
            .catch(err=>{
                console.log(err)
            });
        for(let i=0;i<files.length;i++){
            console.log(files[i])
            this.setState(state=>({
                filenames: [...state.filenames,URL.createObjectURL(files[i])]
            }))
        }
    };

    render(){
        return(
            <div class='container' style={this.dropzoneStyleCreator(this.state.dropzone_color)} onMouseEnter={this.onMouseEnter} onMouseLeave={this.onMouseLeave}>
                <Dropzone onDrop={this.onDrop} noClick={true}>
                    {({getRootProps,getInputProps})=> (
                        <section className="container">
                            <div {...getRootProps({className:'dropzone'})}>
                                <input {...getInputProps()} />
                                <p> Drag image here </p>
                            </div>
                            <aside>
                                <h4>
                                    files
                                    {this.state.dropzone_color}
                                </h4>
                                <ul>
                                    {this.state.filenames.map(file=>
                                        <li>
                                            <p>{file}</p>
                                            <img src={file}/>
                                        </li>)
                                    }
                                </ul>
                            </aside>
                        </section>
                    )}
                </Dropzone>
            </div>
        )
    };
}

export default DropZone;
