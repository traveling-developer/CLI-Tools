.PHONY: build build-all defuddle clean

build: defuddle/defuddle-linux-arm64
	go build -o web2md .

build-all: defuddle/defuddle-linux-x64 defuddle/defuddle-linux-arm64 defuddle/defuddle-darwin-x64 defuddle/defuddle-darwin-arm64 defuddle/defuddle-windows-x64.exe
	GOOS=linux   GOARCH=amd64 go build -o web2md-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -o web2md-linux-arm64 .
	GOOS=darwin  GOARCH=amd64 go build -o web2md-darwin-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -o web2md-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o web2md-windows-amd64.exe .

defuddle: \
	defuddle/defuddle-linux-x64 \
	defuddle/defuddle-linux-arm64 \
	defuddle/defuddle-darwin-x64 \
	defuddle/defuddle-darwin-arm64 \
	defuddle/defuddle-windows-x64.exe

defuddle/defuddle-linux-x64 defuddle/defuddle-linux-arm64 defuddle/defuddle-darwin-x64 defuddle/defuddle-darwin-arm64 defuddle/defuddle-windows-x64.exe:
	docker build -t defuddle-build -f defuddle/dockerfile ../..
	docker create --name tmp defuddle-build
	docker cp tmp:/output/defuddle-linux-x64    defuddle/defuddle-linux-x64
	docker cp tmp:/output/defuddle-linux-arm64  defuddle/defuddle-linux-arm64
	docker cp tmp:/output/defuddle-darwin-x64   defuddle/defuddle-darwin-x64
	docker cp tmp:/output/defuddle-darwin-arm64 defuddle/defuddle-darwin-arm64
	docker cp tmp:/output/defuddle-windows-x64.exe  defuddle/defuddle-windows-x64.exe
	docker rm tmp
	docker rmi defuddle-build

clean:
	rm -f web2md web2md-* defuddle/defuddle-* defuddle/defuddle-*.exe
