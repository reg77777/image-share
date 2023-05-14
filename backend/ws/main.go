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
            return origin == "http://localhost:8080"
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

    post_image:=&PostImage{Data:image.Data,Extension:image.Extension}
    jsonString,err:=json.Marshal(post_image)
    if err!=nil{
        fmt.Println("jsonString error")
    }
    post_res,err:=http.Post("http://172.21.0.3:80","application/json",bytes.NewBuffer(jsonString))
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
    id int
    title string
    image_path string
    category string
    created_at string
}

func Get(w rest.ResponseWriter,r *rest.Request){
    rows,err:=db.Query("SELECT * FROM images")
    if err!=nil{
        fmt.Println("error")
    }
    defer rows.Close()
    for rows.Next(){
        image:=Image{}
        rows.Scan(&image.id,&image.title,&image.image_path,&image.category,&image.created_at)
        fmt.Println(image)
    }
    w.WriteJson("hello")
}
