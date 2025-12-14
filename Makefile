build-new:
	@go build -o ./build/new ./cmd/new/main.go
	@cp ./build/new "$(GO_PATH)/bin"

build-envtoggle:
	@go build -o ./build/envtoggle ./cmd/envtoggle/main.go
	@cp ./build/envtoggle "$(GO_PATH)/bin"

build-livephoto:
	@go build -o ./build/livephoto ./cmd/livephoto/main.go
	@cp ./build/livephoto "$(GO_PATH)/bin"

build-gitclean:
	@go build -o ./build/gitclean ./cmd/gitclean/main.go
	@cp ./build/gitclean "$(GO_PATH)/bin"

build-renamer:
	@go build -o ./build/renamer ./cmd/renamer/main.go
	@cp ./build/renamer "$(GO_PATH)/bin"