# aws lambda golang client

## Install

```bash
go get github.com/lichti/aws-lambda-golang-client
go install github.com/lichti/aws-lambda-golang-client
```

## Example

- Run you golang lambda:

```bash
export _LAMBDA_SERVER_PORT=55555
./your_lambda
```

- Tringging your lambda function

```bash
export _LAMBDA_SERVER_PORT=55555
echo '{"foo": "bar"}' | ./aws-lambda-golang-client
```

## How golang lambda works ?

The golang aws lambda function create a RCP server and wait a remote call with aws parameters(`InvokeRequest`), this parameters contain you json payload. After function execute, the RPC server return a response containing de function return as a payload on aws format (`InvokeResponse`)
