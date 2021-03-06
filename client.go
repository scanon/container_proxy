package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type ProcReq struct {
	Cmd []string `json:"cmd"`
}

type Message struct {
	Msg  string `json:"msg"`
	Line string `json:"line"`
	Err  bool   `json:"error"`
	Exit int    `json:"exit"`
}

type Resp struct {
	Recieved bool   `json:"received"`
	JobId    string `json:"jid"`
}
type Messages struct {
	Messages []Message `json:"msgs"`
}

func main() {
	flag.Parse()
	var sock = "/tmp/api.sock"
	var err error

	cmd := os.Args
	m := ProcReq{Cmd: cmd}
	b, err := json.Marshal(&m)
	if err != nil {
		panic(err)
	}

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", sock)
			},
		},
	}

	var response *http.Response
	response, err = httpc.Post("http://unix/submit", "application/octet-stream", bytes.NewBuffer(b))
	if err != nil {
		log.Fatal(err)
	}

	if response.Body != nil {
		defer response.Body.Close()
	}

	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	var message Resp
	jsonErr := json.Unmarshal(body, &message)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	var messages Messages
	var cont bool = true
	for cont {
		response, err = httpc.Get("http://unix/output/" + message.JobId)
		if err != nil {
			log.Fatal(err)
		}

		if response.Body != nil {
			defer response.Body.Close()
		}

		body, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		jsonErr := json.Unmarshal(body, &messages)
		if jsonErr != nil {
			log.Print(jsonErr)
		}
		for _, msg := range messages.Messages {
			if msg.Msg == "finished" {
				cont = false
				os.Exit(msg.Exit)
			} else if msg.Msg == "output" {
				if msg.Err {
					os.Stderr.WriteString(msg.Line)
				} else {
					os.Stdout.WriteString(msg.Line)
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
