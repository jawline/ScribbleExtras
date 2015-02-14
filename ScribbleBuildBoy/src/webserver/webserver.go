package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"fmt"
	"os"
	"strconv"
	"path/filepath"
)

var addr = flag.String("addr", ":80", "WebServer Service");
var path = "";
var job = false;
var output = "";
var error = "";

func writeSource(source string) string {

	fo, err := os.Create("temp.sc");

	if err != nil {
		return "";
	}

	fo.Write([]byte(source));
	fo.Close();

	return "./temp.sc";
}

func buildHandler(c http.ResponseWriter, req *http.Request) {
	beginProcess("", c, req);
}

func cleanHandler(c http.ResponseWriter, req *http.Request) {
	beginProcess("clean", c, req);
}

func beginProcess(rule string, c http.ResponseWriter, req *http.Request) {

	if job {
		c.WriteHeader(400);
		c.Write([]byte("{reason:\"job already running\"}"));
		return;
	}

	output = "";
	error = "";

	fmt.Printf("Handle %s\n", req.FormValue("source"));
	c.Header().Set("Access-Control-Allow-Origin", "*");

	path, err := filepath.Abs(path);

	if err != nil {
		c.WriteHeader(400);
		c.Write([]byte("{reason:\"" + err.Error() + "\"}"));
		return;
	}

	job = true;
	
	go func() {
		
		cmd := exec.Command("make");
		
		if (rule != "") {
			cmd = exec.Command("make", rule);
		}
		
		cmd.Dir = path;

		fmt.Printf("Done setting up\n", output, error);

		out, err := cmd.CombinedOutput();

		output = string(out);

		if err != nil {
			error = string(err.Error());
		}

		job = false;

		fmt.Printf("%s %s\n", out, err);
	}();

	c.Write([]byte(string("started")));	
}

func statusHandler(c http.ResponseWriter, req *http.Request) {
	c.Write([]byte("{\n\"status\": " + strconv.FormatBool(job) + ",\n\"log\": \"" + output + "\",\n\"error\": \"" + error + "\"\n}"));
}

func main() {

	directoryFlag := flag.String("path", "/home/blake/Dropbox/Current/Scribble/", "The path to the scribble interpretor")
	flag.Parse()
	path = *directoryFlag

	http.Handle("/", http.FileServer(http.Dir("./websrc/")));
	http.HandleFunc("/clean", cleanHandler);
	http.HandleFunc("/build", buildHandler);
	http.HandleFunc("/status", statusHandler);

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err);
	}
}
