package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"
	"fmt"
	"os"
	"path/filepath"
	"encoding/json"
)

type Status struct {
	Building bool
	PreviousJobs int
	LastStep string
	Log string
	Error string
}

var addr = flag.String("addr", ":80", "WebServer Service");
var path = "";
var state = Status{false, 0, "", "", ""};

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
	beginProcess("all", c, req);
}

func cleanHandler(c http.ResponseWriter, req *http.Request) {
	beginProcess("clean", c, req);
}

func testHandler(c http.ResponseWriter, req *http.Request) {
	beginProcess("test", c, req);
}

func pullHandler(c http.ResponseWriter, req *http.Request) {
	beginProcess("resetgitandpull", c, req);
}

func beginProcess(rule string, c http.ResponseWriter, req *http.Request) {

	if state.Building {
		c.WriteHeader(400);
		c.Write([]byte("{reason:\"job already running\"}"));
		return;
	}

	state.PreviousJobs++;
	state.LastStep = rule;

	state.Log = "";
	state.Error = "";

	fmt.Printf("Handle %s\n", req.FormValue("source"));
	c.Header().Set("Access-Control-Allow-Origin", "*");

	path, err := filepath.Abs(path);

	if err != nil {
		c.WriteHeader(400);
		c.Write([]byte("{reason:\"" + err.Error() + "\"}"));
		return;
	}

	state.Building = true;
	
	go func() {
		
		cmd := exec.Command("make");
		
		if (rule != "") {
			cmd = exec.Command("make", rule);
		}
		
		cmd.Dir = path;

		fmt.Printf("Done setting up\n");

		out, err := cmd.CombinedOutput();

		state.Log = string(out);

		if err != nil {
			state.Error = string(err.Error());
		}

		state.Building = false;
		fmt.Printf("%s %s\n", out, err);
	}();

	c.Write([]byte(string("started")));	
}

func statusHandler(c http.ResponseWriter, req *http.Request) {
	data, err := json.Marshal(state);

	if err != nil {
		c.WriteHeader(400);
		c.Write([]byte("{reason:\"" + err.Error() + "\"}"));
		return;
	}

	c.Write(data);
}

func main() {

	directoryFlag := flag.String("path", "./Scribble/", "The path to the scribble interpretor")
	flag.Parse()
	path = *directoryFlag

	http.Handle("/", http.FileServer(http.Dir("./websrc/")));
	http.HandleFunc("/clean", cleanHandler);
	http.HandleFunc("/build", buildHandler);
	http.HandleFunc("/test",  testHandler);
	http.HandleFunc("/pull",  pullHandler);
	http.HandleFunc("/status", statusHandler);

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err);
	}
}
