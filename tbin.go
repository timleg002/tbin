package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
)

//Constants
const (
	LISTEN_ADDR_INIT   = "/init"
	LISTEN_ADDR_SUBMIT = "/submit_paste"
)

var npastes []pastefmt

func main() {
	//httphandle("/rand")
	mux := httpinit()
	err := http.ListenAndServe("", mux) //mux is a handler
	if err != nil {
		fmt.Println(err.Error())
	}
	pasteinit() //init paste file db system
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
	pastewrite_id_impl(pastefmt{
		req.PostForm.Get("text"),
		req.PostForm.Get("author"),
		linkgen(),
	})
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

//TODO backup function for writing pastes from memory to file

/*func pasteswrite(){
	bc, err := json.Marshal(&npastes)

	f, err3 := os.Open("pastes.json")

	if err3 != nil {
		fmt.Println(err3)
	}

	defer f.Close()

	bc = append(bc, byte('\n'))

	fi, _ := f.Stat()
	_, err4 := f.WriteAt(bc, fi.Size()-2)
}*/

func pasteread(id int) pastefmt { //no need for this cause we save to npastes []pastefmt
	return npastes[id]
}

func pastewrite(id int, paste pastefmt) { //this just appends the paste to the file
	bc, err2 := json.Marshal(&paste) //&pastefmt or pastefmt !!!!!!!
	if err2 != nil {
		fmt.Println(err2)
	}

	f, err3 := os.Open("pastes.json")

	if err3 != nil {
		fmt.Println(err3)
	}

	defer f.Close()

	fi, _ := f.Stat()
	_, err4 := f.WriteAt(bc, fi.Size()-2)

	if err4 != nil {
		fmt.Println(err4)
	}

	npastes = append(npastes, paste) //to include (match!) in backup
}

func pastewrite_id_impl(paste pastefmt) {
	pastewrite(len(npastes)+1, paste) //should work idk
}
