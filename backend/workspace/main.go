package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    "log"
    "net/http"
    "fmt"
    "encoding/base64"
    "os"
    "database/sql"
	_ "github.com/go-sql-driver/mysql"
    "strconv"
    "bytes"
    "encoding/json"
    "io"
)

var db *sql.DB

func main() {
    var err error
    db, err=sql.Open("mysql", "root:root@tcp(db:3306)/images")
    if err != nil {
        fmt.Println("setup error");
    }
    err=db.Ping()
    if err != nil {
        fmt.Println("db error");
    }

    api := rest.NewApi()
    api.Use(&rest.CorsMiddleware{
        RejectNonCorsRequests: false,
        OriginValidator: func(origin string, request *rest.Request) bool {
            return true
        },
        AllowedMethods: []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders: []string{"Accept","Authorization","content-type"},
        AccessControlAllowCredentials: true,
        AccessControlMaxAge:           3600,
    })
    router,_:=rest.MakeRouter(
        rest.Post("/upload",Upload),
        rest.Get("/",Get),
    )
    api.SetApp(router)
    log.Fatal(http.ListenAndServe(":3000", api.MakeHandler()))
}

type PostedImage struct{
    Title string
    Data string
    Extension string
}

type PostImage struct{
    Data string
    Extension string
}


func Upload(w rest.ResponseWriter,r *rest.Request){
    fmt.Println("upload")
    image:=PostedImage{}
    err:=r.DecodeJsonPayload(&image)
    if err!=nil{
        fmt.Println("error")
    }
    fmt.Println(image.Extension)
    data,_:=base64.StdEncoding.DecodeString(image.Data)

    res,err:=db.Exec(
        "INSERT INTO images (title) VALUES (?)",
        image.Title,
    )
    id,err:=res.LastInsertId()
    id_str:=strconv.FormatInt(id,10)
    path:="./images/"+id_str+"."+image.Extension
    file,_:=os.Create(path)
    file.Write(data)

    fmt.Println("path",path)
    fmt.Println("id",id)
    fmt.Println("Ex",image.Extension)

    post_image:=&PostImage{Data:image.Data,Extension:image.Extension}
    jsonString,err:=json.Marshal(post_image)
    if err!=nil{
        fmt.Println("jsonString error")
    }
    post_res,err:=http.Post("http://torch:80","application/json",bytes.NewBuffer(jsonString))
    defer post_res.Body.Close()

    if err!=nil{
        fmt.Println("post error")
    }
    body,_:=io.ReadAll(post_res.Body)
    fmt.Println(string(body))

    res,err=db.Exec(
        "UPDATE images SET image_path = ?, category = ? WHERE id = ?",
        path,
        string(body),
        id,
    )
}

type Image struct{
    Id int
    Title string
    Data string
    Image_path string
    Category string
    Created_at string
}

func Get(w rest.ResponseWriter,r *rest.Request){
    fmt.Println("get")
    rows,err:=db.Query("SELECT * FROM images")
    if err!=nil{
        fmt.Println("error")
    }
    defer rows.Close()
    images:=[]Image{}
    for rows.Next(){
        image:=Image{}
        rows.Scan(&image.Id,&image.Title,&image.Image_path,&image.Category,&image.Created_at)

        if image.Title==""{
            continue
        }
        if image.Image_path==""{
            continue
        }
        if image.Category==""{
            continue
        }

        file,err:=os.Open(image.Image_path)
        if err!=nil{
            continue
        }
        fi,_:=file.Stat()
        size:=fi.Size()
        if size==0{
            continue
        }
        buf:=make([]byte,size)
        file.Read(buf)
        data:=base64.StdEncoding.EncodeToString(buf)
        image.Data=data

        images=append(images,image)
    }
    w.WriteJson(&images)
}
