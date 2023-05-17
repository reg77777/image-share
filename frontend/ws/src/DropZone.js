import React from "react";
import Dropzone from "react-dropzone";
import {Button,Container,makeStyles} from "@mui/material";
import axios from "axios";
import {Buffer} from 'buffer';

//post_url='http://133.167.112.49:8081/'
const post_url='http://localhost:8081'

class DropZone extends React.Component{
    constructor(props){
        super(props);
    };

    state={
        extension:"",
        image:"",
        title:"",
        mouseenter:false,
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

    onDrop=(files)=>{
        for(let i=0;i<files.length;i++){
            var file=files[i]
            var reader=new FileReader
            reader.readAsDataURL(file)
            var s=file.name
            var extension=s.split('.').pop()
            var base_name=file.name.split('.').shift()
            reader.onload=()=>{
                var val = reader.result.replace(/data:.*\/.*;base64,/, '');

                axios.post(post_url+'/upload',{data:val,extension:extension,title:s})
            }
        }
    };

    onDragEnter=()=>{
        this.setState({mouseenter:true})
    }

    onMouseLeave=()=>{
        this.setState({mouseenter:false})
    }

    make_style=(mouseenter)=>{
        var color="";
        if(mouseenter){
            color="red"
        }
        else{
            color="black"
        }
        return {
            padding:"100px",
            textAlign:"center",
            border:"2px solid",
            color:color,
        }
    }

    render(){
        return(
          <Dropzone onDrop={this.onDrop} noClick={true}>
            {({getRootProps, getInputProps}) => (
              <section className="container">
                <div {...getRootProps({className: 'dropzone'})}>
                  <input {...getInputProps()} height="500px"/>
                  <p onDragEnter={this.onDragEnter} onMouseLeave={this.onMouseLeave} onDragLeave={this.onMouseLeave}  style={this.make_style(this.state.mouseenter)}>
                ここにファイルをドラッグアンドドロップしてください
                </p>
                </div>
              </section>
            )}
          </Dropzone>
        )
    };
}

export default DropZone;
