package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var PasteHolder []pasteType

type pasteType struct { //no id cause count of array is the ID
	Text   string `json:"text"` //this was the error um...
	Author string `json:"author"`
	Link   string `json:"link"`
}

func ShowPasteHandler(w http.ResponseWriter, req *http.Request) {
	str := strings.TrimPrefix(req.URL.RequestURI(), UriLink)
	fmt.Println(str)
	paste := PasteFmtWithLinkName(str)
	if paste == nil {
		_, _ = fmt.Fprintf(w, "Paste doesn't exist!")
	} else {
		_, _ = fmt.Fprintf(w, "paste: %s\nauthor: %s\n", paste.Text, paste.Author)
	}
}

func PasteFmtWithLinkName(link string) *pasteType {
	for i := range PasteHolder {
		if PasteHolder[i].Link == link {
			return &PasteHolder[i]
		}
	}
	return nil
}

func InitPasteSlice() {
	byteContent, err := ioutil.ReadFile("pastes.json")
	if err != nil {
		fmt.Println(err)
	}
	err2 := json.Unmarshal(byteContent, &PasteHolder) //load content into &PasteHolder
	if err2 != nil {
		fmt.Println(err2)
	}
}

//No backup function needed since WritePasteImplyId() updates both the PasteHolder object and both the file so no need for file backup, cause memory loss will not cause anything

func (paste pasteType) WritePaste() { //this just appends the paste to the file
	bc, err := json.Marshal(&paste)
	ErrHandle(err, "json.Marshal")
	f, err := os.OpenFile("pastes.json", os.O_WRONLY, os.ModePerm)
	ErrHandle(err, "OpenFile")

	defer f.Close()

	fi, _ := f.Stat()
	_, err = f.WriteAt([]byte{'}', ','}, fi.Size()-3)
	ErrHandle(err, "WriteAt") //it overwrites the '}' character so we need to put it back hehe 😹 ;;
	_, err = f.WriteAt(bc, fi.Size()-1)
	ErrHandle(err, "WriteAt")
	_, err = f.WriteAt([]byte{'\n', ']'}, fi.Size()-1+int64(len(bc)))
	ErrHandle(err, "WriteAt") //ends the json

	PasteHolder = append(PasteHolder, paste) //to include (match!) in backup
}
