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
    "time"
	"math/rand"
	"sync"
)

var db *sql.DB
var name_map map[string]string

func main() {
	name_map = make(map[string]string)
    var err error
    db, err=sql.Open("mysql", "root:root@tcp(db:3306)/images")
    if err != nil {
        fmt.Println("setup error");
    }
    for{
        err=db.Ping()
        if err != nil {
            fmt.Println("wait db");
        } else {
			fmt.Println("db is ok");
			break
		}
		time.Sleep(1*time.Second)
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
		rest.Get("/get",Get),
		rest.Get("/getnum",Getnum),
		rest.Get("/getname",Getname),
		rest.Get("/long-polling",LongPolling),
    )
    api.SetApp(router)
    fmt.Println("start api")
    log.Fatal(http.ListenAndServe(":3000", api.MakeHandler()))
}

type PostedImage struct{
    Title string
    Data string
	Post_user string
    Extension string
}

type PostImage struct{
    Data string
    Extension string
}

func Genid(n int) string {
	var s = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = s[rand.Intn(len(s))]
	}
	return string(b)
}

func Getname(w rest.ResponseWriter,r *rest.Request){
	token, err := r.Cookie("session_id")
	name:="guest"
	if err == nil{
		id:=token.Value
		name=name_map[id]
		if name==""{
			name="guest"
		}
		delete(name_map,id)
		new_id := Genid(20)
		name_map[new_id]=name
		cookie := &http.Cookie{Name: "session_id", Value: new_id}
		writer := w.(http.ResponseWriter)
		http.SetCookie(writer, cookie)
	}
	w.WriteJson(&name)
}

var msgCh = make(chan int,100)
var wg sync.WaitGroup

func Upload(w rest.ResponseWriter,r *rest.Request){
    fmt.Println("upload")
    image:=PostedImage{}
    err:=r.DecodeJsonPayload(&image)
	fmt.Println(image.Title)
    if err!=nil{
        fmt.Println("error")
    }
    fmt.Println(image.Extension)
    data,_:=base64.StdEncoding.DecodeString(image.Data)

	if image.Post_user==""{
		image.Post_user="guest"
	}
    res,err:=db.Exec(
        "INSERT INTO images (title,post_user) VALUES (?,?)",
        image.Title,
		image.Post_user,
    )
	if image.Post_user!="guest"{
		fmt.Println("not guest")
		new_id := Genid(20)
		name_map[new_id]=image.Post_user
		cookie := &http.Cookie{Name: "session_id", Value: new_id}
		writer := w.(http.ResponseWriter)
		http.SetCookie(writer, cookie)
	}
    id,err:=res.LastInsertId()
    id_str:=strconv.FormatInt(id,10)
    path:="./images/"+id_str+"."+image.Extension
    file,_:=os.Create(path)
    file.Write(data)

    fmt.Println("path",path)
    fmt.Println("id",id)
    fmt.Println("Ex",image.Extension)
	fmt.Println("title",image.Title)
	fmt.Println("post_user",image.Post_user)

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
	close(msgCh)
	wg.Wait()
	msgCh = make(chan int,100)
	fmt.Println("finish polling send")
}

type Image struct{
    Id int
    Title string
	Post_user string
    Data string
    Image_path string
    Category string
    Created_at string
}

func Getnum(w rest.ResponseWriter,r *rest.Request){
	fmt.Println("get num")
	rows,err:=db.Query("SELECT COUNT(*) FROM images")

	if err!=nil{
		fmt.Println("error")
	}else{
		var count int
		for rows.Next(){
			rows.Scan(&count)
		}
		fmt.Println(count)
		w.WriteJson(&count)
	}
}

func LongPolling(w rest.ResponseWriter,r *rest.Request){
	fmt.Println("long polling")
	wg.Add(1)
	<-msgCh
	fmt.Println("long polling update")
	Get(w,r)
	wg.Done()
}

func Get(w rest.ResponseWriter,r *rest.Request){
    fmt.Println("test test get")
    rows,err:=db.Query("SELECT * FROM images")
    if err!=nil{
        fmt.Println("error")
    }
    defer rows.Close()
    images:=[]Image{}
    for rows.Next(){
        image:=Image{}
        rows.Scan(&image.Id,&image.Title,&image.Post_user,&image.Image_path,&image.Category,&image.Created_at)
		fmt.Println(image)

        if image.Title==""{
	    fmt.Println("Title is None")
            continue
        }
        if image.Image_path==""{
	    fmt.Println("Image_path is None")
            continue
        }
        if image.Category==""{
	    fmt.Println("category is none")
            continue
        }

        file,err:=os.Open(image.Image_path)
        if err!=nil{
	    fmt.Println("file can not open")
            continue
        }
        fi,_:=file.Stat()
        size:=fi.Size()
        if size==0{
	    fmt.Println("file size is 0")
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