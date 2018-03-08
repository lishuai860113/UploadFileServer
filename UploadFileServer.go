package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "io"
    "os"
    "flag"
)

var indextpl string = `<html>
<head>
       <title>Upload file</title>
</head>
<body>
<form enctype="multipart/form-data" action="/upload" method="post">
    <input type="file" name="uploadfile" />
    <input type="submit" value="upload" />
</form>

</body>
</html>`

func upload(w http.ResponseWriter, r *http.Request) {
       fmt.Println("method:", r.Method)
       if r.Method == "GET" {
           t, _ := template.New("foo").Parse(indextpl)
           t.Execute(w, nil)
       } else {
           r.ParseMultipartForm(32 << 20)
           file, handler, err := r.FormFile("uploadfile")
           if err != nil {
               fmt.Println(err)
               return
           }
           defer file.Close()
           t, _ := template.New("foo").Parse(indextpl)
           t.Execute(w, nil)
           f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
           if err != nil {
               fmt.Println(err)
               return
           }
           defer f.Close()
           io.Copy(f, file)
       }
}



func main() {
    portPtr := flag.String("port", "9090", "listen port")
    flag.Parse()
    http.HandleFunc("/upload", upload)
    http.Handle("/", http.FileServer(http.Dir(".")))
    err := http.ListenAndServe(":" + *portPtr, nil) 
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
