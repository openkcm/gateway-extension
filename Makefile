NAME := gateway-extension

IMGTAG ?= latest
IMG ?= localhost/$(NAME):${IMGTAG}

tools.dir = tools

.PHONY: install-controller-gen
install-controller-gen:
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@latest

.PHONY: generate-crds
generate-crds:
	controller-gen crd:allowDangerousTypes=true,generateEmbeddedObjectMeta=true \
			object:headerFile="$(tools.dir)/boilerplate.generatego.txt",year=2025 paths="{./...}" \
            output:crd:artifacts:config=charts/gateway-extension/crds; \
    goimports -w .

.PHONY: build
build: clean
	go build -o $(NAME) ./cmd/$(NAME)
	sha256sum $(NAME)

.PHONY: docker-build
docker-build:
	GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o $(NAME) ./cmd/$(NAME)
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push:
	docker push ${IMG}

.PHONY: docker-build-push
docker-build-push: docker-build docker-push

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run -v --fix ./...

.PHONY: clean
clean:
	rm -f cover.out cover.html $(NAME)
	rm -rf cover/

.PHONY: test
test: clean
	mkdir -p cover/integration cover/unit
	go clean -testcache

	# unit tests
	go test -count=1 -race -cover ./... -args -test.gocoverdir="${PWD}/cover/unit"

	# integration tests
	GOCOVERDIR="${PWD}/cover/integration" go test -count=1 -race --tags=integration ./integration

	# merge coverage
	go tool covdata textfmt -i=./cover/unit,./cover/integration -o cover.out

	# On a Mac, you can use the following command to open the coverage report in the browser
	# go tool cover -html=cover.out -o cover.html && open cover.html
