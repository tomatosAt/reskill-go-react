# Include External Env and export system Env
export
# VERSION=$(shell git describe --always --long --dirty)
VERSION=v0.0.0
OUTPUT_DIR=./out
DOCKER_IMG=github.com/tomatosAt/reskill-go-react

# CI automate version set
ifeq (${CI}, true)
	VERSION=${CI_COMMIT_REF_NAME}-build${CI_PIPELINE_ID}
	DOCKER_IMG=${CI_REGISTRY_IMAGE}:${CI_COMMIT_REF_NAME}
endif

.PHONY: version clean local compile docker push dev

version:
	echo ${VERSION}

clean:
	rm -rf ./out

local:
	docker build --target=server --build-arg="APP_VERSION=beta" -t ${DOCKER_IMG} .
	
build: compile
	DOCKER_BUILDKIT=0 docker build --rm -t $(DOCKER_IMG):$(VERSION) .

compile:
	go build -ldflags "-X 'main.Version=${VERSION}'" -o ./out/server ./cmd/server/main.go
	./out/server version

docker:
	docker build --target=server --build-arg="APP_VERSION=${VERSION}" -t ${DOCKER_IMG} .

push:
	docker push ${DOCKER_IMG}:$(VERSION)

dev:
	go run -ldflags "-X 'main.Version=${VERSION}'" cmd/server/main.go start

build-mac: compile-for-build-mac
	docker buildx build --platform linux/amd64 -t $(DOCKER_IMG):$(VERSION) .

compile-for-build-mac:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.VERSION=$(VERSION)'" -o $(OUTPUT_DIR)/server ./cmd/server/main.go

db-local:
	docker-compose -f docker-compose-mock-db.yml up