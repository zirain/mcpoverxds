DOCKER_REPO ?= zirain/mcp-over-xds
#DOCKER_TAG ?= $(shell git rev-parse --abbrev-ref HEAD)-$(shell date -u +"%Y-%m-%d")-$(shell git rev-parse --short HEAD)
DOCKER_TAG ?= $(shell git rev-parse --abbrev-ref HEAD)-latest
GO_BUILD ?= go build
BIN_OUT ?= ./out

.PHONY: clean
clean:
	rm -rf $(BIN_OUT)
	#docker rmi $(DOCKER_REPO):$(DOCKER_TAG)
	docker image prune -f

.PHONY: build
build:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o=$(BIN_OUT)/mcp-over-xds ./cmd

.PHONY: docker-build
docker-build: clean
	docker build . -t $(DOCKER_REPO):$(DOCKER_TAG)

.PHONY: docker-publish
docker-publish: docker-build
	docker push $(DOCKER_REPO):$(DOCKER_TAG)


.PHONY: k8s
k8s: docker-build
	kubectl delete -f k8s/all-in-one.yaml && kubectl apply -f k8s/all-in-one.yaml
