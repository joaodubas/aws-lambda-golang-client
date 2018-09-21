

clean:
	rm -rfv terraform/.terraform
	rm -fv terraform_$(TERRAFORM_VERSION)_linux_amd64.zip
	rm -fv main
	rm -fv terraform/lambda_julius.zip

fmt:
	go fmt

build: fmt
	GOOS=linux go build lambda_client.go

run: build
	echo -e "{\"foo\": \"bar\"}" | ./lambda_client
