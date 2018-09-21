package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/rpc"
	"os"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeNamedPipe != os.ModeNamedPipe || info.Size() < 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: echo '{\"foo\": \"bar\"}' | lambda_client")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	log.Println("Stdin:\n", string(output))

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
