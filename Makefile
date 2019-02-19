LAMBDA_FUNCTION_NAME = "lambda-go-demo"

build: build/handler

bundle: build/lambda.zip

run:
	go run main.go

deploy: bundle
	aws lambda update-function-code \
		--function-name $(LAMBDA_FUNCTION_NAME) \
		--zip-file fileb://build/lambda.zip \
		--publish 

clean:
	@rm -f build/handler build/lambda.zip

build/handler:
	mkdir -p build
	GOOS=linux GOARCH=amd64 go build -o build/handler main.go

build/lambda.zip: build/handler
	cd build && zip lambda.zip handler

.PONNY: build run bundle deploy clean
