test:
	@go test ./... -v
cover:
	@go test ./... -v -coverprofile cover.out.tmp && cat cover.out.tmp | grep -v "mock_" > cover.out && go tool cover -html=cover.out -o cover.html && firefox cover.html
