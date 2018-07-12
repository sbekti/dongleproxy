package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

type message struct {
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	binary, lookErr := exec.LookPath("asterisk")
	if lookErr != nil {
		panic(lookErr)
	}

	cmd := exec.Command(binary, "-rx", "dongle show devices")

	cmd.Env = os.Environ()

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out.String())

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(out.String()))
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var message message
	_ = json.NewDecoder(r.Body).Decode(&message)

	binary, lookErr := exec.LookPath("asterisk")
	if lookErr != nil {
		panic(lookErr)
	}

	cmd := exec.Command(binary, "-rx",
		fmt.Sprintf("dongle sms dongle0 %s %s",
			message.Recipient,
			message.Content))

	cmd.Env = os.Environ()

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", out.String())

	json.NewEncoder(w).Encode(message)
}

func main() {
	port := os.Args[1]
	router := mux.NewRouter()
	router.HandleFunc("/", getStatus).Methods("GET")
	router.HandleFunc("/", sendMessage).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
