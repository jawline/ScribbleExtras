package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"fmt"
	"os"
)

var addr = flag.String("addr", ":9760", "WebServer Service")
var path = ""

func writeSource(source string) string {

	fo, err := os.Create("temp.sc")

	if err != nil {
		return ""
	}

	fo.Write([]byte(source))

	fo.Close()

	return "./temp.sc"
}

func buildHandler(c http.ResponseWriter, req *http.Request) {
	fmt.Printf("Handle %s\n", req.FormValue("source"))

	c.Header().Set("Access-Control-Allow-Origin", "*")

	out, err := exec.Command(path, writeSource(req.FormValue("source"))).Output()

	if err != nil {
		c.Write([]byte("Error:"+err.Error()+"\n"+ string(out)))
	} else {
		c.Write([]byte(out))
	}
}

func main() {
	directoryFlag := flag.String("path", "scribble", "The path to the scribble directory")
	flag.Parse()
	path = *directoryFlag

	http.HandleFunc("/startbuild", buildHandler)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
