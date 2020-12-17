package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
)

//Constants
const (
	ListenAddrInit         = "/"
	ListenAddrSubmit       = "/submit_paste"
	ListenAddrLinkWildcard = "/paste/"
	ListenAddrStatic       = "/static/"
	UriLink                = "/paste/" //URI LINK IS RELATIVE TO WEBPAGE, NOT RELATIVE TO THE INTRANET/INTERNET
)

func main() {
	InitPasteSlice() //init paste file db system
	mux := HttpInit()
	err := http.ListenAndServe("", mux) //mux is a handler
	ErrHandle(err, "mux")
	http.Handle(ListenAddrStatic, http.FileServer(http.Dir("static/")))
}

func HttpInit() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(ListenAddrInit, SiteIndexHandler)
	mux.HandleFunc(ListenAddrSubmit, SubmitPasteHandler)
	mux.HandleFunc(ListenAddrLinkWildcard, ShowPasteHandler)
	return mux
}

func SiteIndexHandler(w http.ResponseWriter, _ *http.Request) {
	content, _ := ioutil.ReadFile("init.html")
	_, _ = fmt.Fprintf(w, "%v", string(content))
}

func SubmitPasteHandler(w http.ResponseWriter, req *http.Request) { //POST
	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	strLink := RandomLink()
	pasteType{
		req.PostForm.Get("paste"), //paste in html, text in json and code... fucked up but 😈🚫🤫🤘
		req.PostForm.Get("author"),
		strLink,
	}.WritePaste()
	_, _ = fmt.Fprintf(w, `<html>
<head></head><body></body>
`+"Your paste is available at "+`<a href="`+"%s"+`">here</a>`, UriLink+strLink)
	//_, _ = w.Write([]byte("Your paste is available at " + `<a href="` + URI_LINK+str_link + `">here</a>`))
}

func RandomLink() string { //maybe some checking if the value already exists?🤔
	ret := make([]byte, 8) //should be 8 bytes == 8 chars in 1 word
	for i := 0; i < 8; i++ {
		ret[i] = RandByte()
	}
	return string(ret)
}

var acceptableChars = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

func RandByte() byte {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(acceptableChars))))
	return acceptableChars[num.Int64()]
}

func ErrHandle(e error, errorDesc string) {
	if e != nil {
		_ = fmt.Errorf("errorin: %s, error!: %s", errorDesc, e.Error())
	}
}
