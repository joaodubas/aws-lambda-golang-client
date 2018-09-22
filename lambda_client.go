package main

import (
	"bufio"
	"errors"
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

	client, err := newClient()
	if err != nil {
		log.Fatalf("client: error: %v", err)
	}

	resp, err := rpcCall(client, []byte(string(output)))
	if err != nil {
		log.Fatalf("rpc: communication error: %v", err)
	}
	log.Printf("rpc: get: %s", string(resp.Payload))
}

// stdinFromPipe check if stdin was set from a pipe.
func stdinFromPipe(info os.FileInfo) bool {
	return info.Mode()&os.ModeNamedPipe == os.ModeNamedPipe
}

// stdinIsEmpyt check if stdin is empty.
func stdinIsEmpty(info os.FileInfo) bool {
	return info.Size() < 0
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

// newClient create a rpc client for a given lambda server.
func newClient() (*rpc.Client, error) {
	port := os.Getenv("_LAMBDA_SERVER_PORT")
	if port == "" {
		return nil, errors.New("missing _LAMBDA_SERVER_PORT env variable")
	}

	addr := fmt.Sprintf("localhost:%s", port)
	log.Printf("client: connecting to: tcp://%s", addr)

	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	log.Printf("client: connected to: tcp://:", addr)

	return client, nil
}

// rpcCall invoke lambda function, passing the given content.
func rpcCall(client *rpc.Client, content []byte) (messages.InvokeResponse, error) {
	resp := messages.InvokeResponse{}
	err := client.Call("Function.Invoke", newRequest(content), &resp)
	return resp, err
}

// newRequest create a lambda request with the given content.
func newRequest(content []byte) messages.InvokeRequest {
	return messages.InvokeRequest{
		Payload: content,
	}
}
