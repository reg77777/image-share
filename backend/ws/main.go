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
    file,_:=os.Create("./images/"+strconv.FormatInt(id,10)+"."+image.Extension)
    file.Write(data)
}

func Get(w rest.ResponseWriter,r *rest.Request){
    fmt.Println("get")
    w.WriteJson("hello")
}
