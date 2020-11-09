package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//Constants
const (
	LISTEN_ADDR_INIT = "/init"
	LISTEN_ADDR_SUBMIT = "/submit_paste"
)

func main() {
	//httphandle("/rand")
	mux := httpinit()
	err := http.ListenAndServe("", mux) //mux is a handler
	if err != nil {
		fmt.Println(err.Error())
	}
}

func httpinit() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(LISTEN_ADDR_INIT, siteindexhandler)
	mux.HandleFunc(LISTEN_ADDR_SUBMIT, submitpastehandler)
	return mux
}

func httphandle(url string){
	http.HandleFunc(url, nil)//for now
}

func httplisten(){

}

func siteindexhandler(w http.ResponseWriter, req *http.Request) {
	content, _ := ioutil.ReadFile("init.html")
	_, _ = fmt.Fprintf(w, "%v", string(content))
}

func submitpastehandler(w http.ResponseWriter, req *http.Request) {//POST
	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(req.PostForm.Get("paste"))//type name not label
}

