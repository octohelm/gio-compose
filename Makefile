APP_ID = tech.innoai.example

dev:
	go run ./cmd/example

BUILDER_IMAGE = ghcr.io/innoai-tech/gogio:builder

DOCKER_MAKE = docker run -it \
	-v=${PWD}:/go/src \
	-w=/go/src \
	--platform=linux/${GOARCH} \
	--entrypoint=/usr/bin/make $(BUILDER_IMAGE)

build.windows:
	gogio -target=windows -o ./build/example.exe ./cmd/example

build.linux:
	CGO_ENABLED=1 go build -o=./build/example_linux_${GOARCH} ./cmd/example

build.macos:
	gogio -target=macos -o=./build/example.app ./cmd/example

build.ios-no-sign:
	gogio -target ios -appid=$(APP_ID) -o=./build/example.ipa ./cmd/example

build.android:
	ANDROID_SDK_ROOT = ${ANDROID_HOME} \
		gogio -target android -appid=$(APP_ID) -o ./build/example.apk ./cmd/example

pkg.%:
	$(DOCKER_MAKE) build.$*

install.apk:
	adb install build/example.apk

debug.ios:
	xcrun simctl install booted ./build/example.app

ship.builder:
	docker buildx build --push --platform=linux/amd64,linux/arm64 -t $(BUILDER_IMAGE) -f ./hack/Dockerfile .

test:
	go test -v ./pkg/...

test.race:
	go test -v -race ./pkg/...

doc:
	godoc -http=:6060

fmt:
	goimports -w -l .

dep:
	go get -u ./...

GIO_COMPOSE = go run ./cmd/giocompose

preview: preview.install
	cd ./pkg && giocompose preview ./component/m3 ButtonPreview

preview.install:
	 go build -o $(shell go env GOPATH)/bin/giocompose ./cmd/giocompose


