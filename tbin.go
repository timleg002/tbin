package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
)

//Constants
const (
	LISTEN_ADDR_INIT          = "/init"
	LISTEN_ADDR_SUBMIT        = "/submit_paste"
	LISTEN_ADDR_LINK_WILDCARD = "/paste/"
	URI_LINK                  = "/paste/" //URI LINK IS RELATIVE TO WEBPAGE, NOT RELATIVE TO THE INTRANET/INTERNET
)

var npastes []pastefmt

func main() {
	//httphandle("/rand")
	pasteinit() //init paste file db system
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
	mux.HandleFunc(LISTEN_ADDR_LINK_WILDCARD, pasteshowhandler)
	return mux
}

func pasteshowhandler(w http.ResponseWriter, req *http.Request) {
	str := strings.TrimPrefix(req.URL.RequestURI(), URI_LINK)
	fmt.Println(str)
	paste := returnpastefmtwithlinkname(str)
	if paste == nil {
		_, _ = fmt.Fprintf(w, "Paste doesn't exist!")
	} else {
		_, _ = fmt.Fprintf(w, "paste: %s\nauthor: %s\n", paste.Text, paste.Author)
	}
}

func returnpastefmtwithlinkname(link string) *pastefmt {
	for i := range npastes {
		if npastes[i].Link == link {
			return &npastes[i]
		}
	}
	return nil
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
	linkstr := linkgen()
	pastewrite_id_impl(pastefmt{
		req.PostForm.Get("paste"), //paste in html, text in json and code... fucked up but 😈🚫🤫🤘
		req.PostForm.Get("author"),
		linkstr,
	})
	_, _ = fmt.Fprintf(w, `<html>
<head></head><body></body>
`+"Your paste is available at "+`<a href="`+"%s"+`">here</a>`, URI_LINK+linkstr)
	//_, _ = w.Write([]byte("Your paste is available at " + `<a href="` + URI_LINK+linkstr + `">here</a>`))
}

func linkgen() string { //maybe some checking if the value already eixsts?🤔
	ret := make([]byte, 8) //should be 8 bytes == 8 chars in 1 word
	for i := 0; i < 8; i++ {
		ret[i] = randbyte()
	}
	return string(ret)
}

var acceptablechars = []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

func randbyte() byte {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(acceptablechars))))
	return acceptablechars[num.Int64()]
}

func errhandle(e error, errordesc string) {
	if e != nil {
		_ = fmt.Errorf("errorin: %s, error!: %s", errordesc, e.Error())
	}
}

type pastefmt struct { //no id cause count of array is the ID
	Text   string `json:"text"` //this was the error um...
	Author string `json:"author"`
	Link   string `json:"link"`
}

func pasteinit() {
	bytecontent, err := ioutil.ReadFile("pastes.json")
	if err != nil {
		fmt.Println(err)
	}
	err2 := json.Unmarshal(bytecontent, &npastes) //load content into &npastes
	if err2 != nil {
		fmt.Println(err2)
	}
}

//No backup function needed since pastewrite() updates both the npastes object and both the file so no need for file backup, cause memory loss will not cause anything

func pasteread(id int) pastefmt { //no need for this cause we save to npastes []pastefmt
	return npastes[id]
}

func pastewrite(id int, paste pastefmt) { //this just appends the paste to the file
	bc, err2 := json.Marshal(&paste) //&pastefmt or pastefmt !!!!!!!
	if err2 != nil {
		fmt.Println(err2)
	}

	f, err3 := os.OpenFile("pastes.json", os.O_WRONLY, os.ModePerm)

	if err3 != nil {
		fmt.Println(err3)
	}

	defer f.Close()

	fi, _ := f.Stat()
	_, err6 := f.WriteAt([]byte{'}', ','}, fi.Size()-3) //it overwrites the '}' character so we need to put it back hehe😹
	_, err4 := f.WriteAt(bc, fi.Size()-1)
	_, err5 := f.WriteAt([]byte{'\n', ']'}, fi.Size()-1+int64(len(bc))) //ends the json

	if err4 != nil {
		fmt.Println(err4)
	}

	if err5 != nil {
		fmt.Println(err5)
	}

	if err6 != nil {
		fmt.Println(err6)
	}

	npastes = append(npastes, paste) //to include (match!) in backup
}

func pastewrite_id_impl(paste pastefmt) {
	pastewrite(len(npastes)+1, paste) //should work idk
}
