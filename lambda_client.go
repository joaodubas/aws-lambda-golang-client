package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

func main() {
	output, err := fetchContent()
	if err != nil {
		log.Fatalf("read: couldn't read from stdin: %v", err)
	}
	log.Printf("read: get content: %s", string(output))

	client, err := newClient()
	if err != nil {
		log.Fatalf("client: error: %v", err)
	}

	resp, err := rpcCall(client, output)
	if err != nil {
		log.Fatalf("rpc: communication error: %v", err)
	}
	log.Printf("rpc: get: %s", string(resp.Payload))
}

// fetchContent fetch content either from args or standard input.
func fetchContent() ([]byte, error) {
	useArgs := len(os.Args) == 2
	if useArgs {
		return fetchContentFromArgs()
	} else {
		return fetchContentFromStdin()
	}
}

// fetchContentFromArgs fetch content from first arg sent to command.
func fetchContentFromArgs() ([]byte, error) {
	content := os.Args[1]
	if content == "-h" {
		help()
	}
	return []byte(os.Args[1]), nil
}

// fetchContentFromStdin fetch content from standard input.
func fetchContentFromStdin() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}

// help show help message and exit.
func help() {
	_, cmd := filepath.Split(os.Args[0])
	msg := `%s invoke local lambda function, through rpc.

Usage:

    $ echo '{"foo": "bar"}' | lambda_client

    $ lambda_client '{"foo": "bar"}'

    $ lambda_client <<EOF
    > {"foo": "bar"}
    > EOF
`
	fmt.Fprintf(os.Stderr, msg, cmd)
	os.Exit(1)
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
	log.Printf("client: connected to: tcp://%s", addr)

	return client, nil
}

// rpcCall invoke lambda function, passing the given content.
func rpcCall(client *rpc.Client, content []byte) (messages.InvokeResponse, error) {
	resp := messages.InvokeResponse{}
	err := client.Call("Function.Invoke", newRequest(content), &resp)
	return resp, err
}

// newRequest create a lambda request with the given content.
func newRequest(content []byte) *messages.InvokeRequest {
	return &messages.InvokeRequest{
		Payload: content,
	}
}
