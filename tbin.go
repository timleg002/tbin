package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//Constants
const (
	LISTEN_ADDR_INIT   = "/init"
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

func httphandle(url string) {
	http.HandleFunc(url, nil) //for now
}

func httplisten() {

}

func siteindexhandler(w http.ResponseWriter, req *http.Request) {
	content, _ := ioutil.ReadFile("init.html")
	_, _ = fmt.Fprintf(w, "%v", string(content))
}

func submitpastehandler(w http.ResponseWriter, req *http.Request) { //POST
	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(req.PostForm.Get("paste")) //type name not label
	pasteread(0)
}

type pastefmt struct { //no id cause count of array is the ID
	Text   string `json:"text"` //this was the error um...
	Author string `json:"author"`
	Link   string `json:"link"`
}

func pasteread(id int) {
	bytecontent, err := ioutil.ReadFile("pastes.json")
	if err != nil {
		fmt.Println(err)
	}
	var pastes []pastefmt
	err2 := json.Unmarshal(bytecontent, &pastes)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(pastes[id].Author)
}

func pastewrite(id int, paste pastefmt) {
	bytecontent, err := ioutil.ReadFile("pastes.json")
	if err != nil {
		fmt.Println(err)
	}

	npaste, err2 := json.Marshal(paste)
	if err2 != nil {
		fmt.Println(err2)
	}

	bytecontent = append(bytecontent, npaste...) //unpacking slice

	//TODO pastewrite() method,
}
