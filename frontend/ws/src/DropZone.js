import React from "react";
import Dropzone from "react-dropzone";
import {Button,Container,makeStyles} from "@mui/material";
import axios from "axios";
import {Buffer} from 'buffer';

axios.defaults.baseURL='http://localhost:8081'

class DropZone extends React.Component{
    constructor(props){
        super(props);
    };

    state={
        extension:"",
        image:"",
        title:"",
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

    onDrop=(file)=>{
        axios.post("http://localhost:8081/upload",{data:this.state.image,extension:this.state.extension,title:this.state.title})
            .then(res=>{
                console.log(res)
            })
            .catch(err=>{
                console.log(err)
            });
        /*
        this.setState(state=>({
            filenames: [...state.filenames,URL.createObjectURL(file)]
        }))
        */
    };

    onChange=(e)=>{
        var reader=new FileReader;
        reader.readAsDataURL(e.target.files[0])
        var extension=e.target.files[0].name.split('.').pop()
        var base_name=e.target.files[0].name.split('.').shift()
        reader.onload=()=>{
            var val = reader.result.replace(/data:.*\/.*;base64,/, '');
            this.setState({image:val,extension:extension,title:base_name})
            console.log(val)
        }
    }

    render(){
        return(
            <div>
                <form onSubmit={this.onDrop}>
                    <input type="file" name="file" onChange={this.onChange}/>
                    <input type="submit"/>
                </form>
            {/*
            <div class='container' style={this.dropzoneStyleCreator(this.state.dropzone_color)} onMouseEnter={this.onMouseEnter} onMouseLeave={this.onMouseLeave}>
                <input type="submit"/>
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
            */}
            </div>

        )
    };
}

export default DropZone;
