package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

func main() {
	if info, err := os.Stdin.Stat(); err != nil {
		log.Fatalf("stat: couldn't fetch stdin stat: %v", err)
	} else if !stdinFromPipe(info) || stdinIsEmpty(info) {
		help()
	}

	output, err := read(os.Stdin)
	if err != nil {
		log.Fatalf("read: couldn't read from stdin: %v", err)
	}
	log.Printf("read: get content: %s", string(output))

	port := os.Getenv("_LAMBDA_SERVER_PORT")
	if port == "" {
		log.Fatal("You need export the variable _LAMBDA_SERVER_PORT")
	}

	log.Println("Connecting to: tcp://localhost:" + port)
	client, err := rpc.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to: tcp://localhost:" + port)

	req := messages.InvokeRequest{}
	req.Payload = []byte(string(output))
	rep := messages.InvokeResponse{}
	err = client.Call("Function.Invoke", req, &rep)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Stdout:\n", string(rep.Payload))
}

// stdinFromPipe check if stdin was set from a pipe.
func stdinFromPipe(info os.FileInfo) bool {
	return info.Mode()&os.ModeNamedPipe != os.ModeNamedPipe
}

// stdinIsEmpyt check if stdin is empty.
func stdinIsEmpty(info os.FileInfo) bool {
	return info.Size() <= 0
}

// help show help message and exit.
func help() {
	_, cmd := filepath.Split(os.Args[0])
	msg := `%s is intended to work with pipes.

Usage: echo '{"foo": "bar"}' | lambda_client
`
	fmt.Printf(msg, cmd)
	os.Exit(1)
}

// read all content from a given reader into an array of runes.
func read(r io.Reader) ([]rune, error) {
	buff := bufio.NewReader(r)
	var output []rune

	for {
		if input, _, err := buff.ReadRune(); err != nil {
			if err == io.EOF {
				break
			} else {
				return output, err
			}
		} else {
			output = append(output, input)
		}
	}
	return output, nil
}
