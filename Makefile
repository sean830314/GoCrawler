IMAGE_REGISTRY ?= ekko771
GO_CRAWLER_SERVICE ?= go-crawler-service
CONSUMER_SERVICE ?= crawler-consumer
IMAGE := $(IMAGE_REGISTRY)/$(SERVICE)
BUILD_DATE = $(shell date +%Y-%m-%dT%H:%M:%S)
VERSION ?= latest
TAG := $(shell git describe --tags --always --dirty)
BLUE='\033[0;34m'
NC='\033[0m'

build-gocrawler-service-image: $(GOCRAWLER_IMAGE_TAGS)
$(GOCRAWLER_IMAGE_TAGS):
	@echo "start build gocrawler"
	@echo "\n${BLUE}Building Service image with labels:${NC}\n"
	@sed                                                      \
	    -e 's|{VERSION}|$(VERSION)|g'                         \
	    -e 's|{IMAGE}|$(IMAGE_REGISTRY)/$(GO_CRAWLER_SERVICE)|g'                             \
	    -e 's|{VCS-REF}|$(TAG)|g'                             \
		-e 's|{BUILD-DATE}|$(BUILD_DATE)|g'                   \
	    Dockerfile/goCrawler.Dockerfile | docker build -t $(IMAGE_REGISTRY)/$(GO_CRAWLER_SERVICE):$@ -f- . -f ./Dockerfile/goCrawler.Dockerfile
.PHONY: build-gocrawler-service-image

build-consumer-image: $(CONSUMER_IMAGE_TAGS)
$(CONSUMER_IMAGE_TAGS):
	@echo "start build consumer"
	@echo "\n${BLUE}Building Service image with labels:${NC}\n"
	@sed                                                      \
	    -e 's|{VERSION}|$(VERSION)|g'                         \
	    -e 's|{IMAGE}|$(IMAGE_REGISTRY)/$(CONSUMER_SERVICE)|g'                             \
	    -e 's|{VCS-REF}|$(TAG)|g'                             \
		-e 's|{BUILD-DATE}|$(BUILD_DATE)|g'                   \
	    Dockerfile/Consumer.Dockerfile | docker build -t $(IMAGE_REGISTRY)/$(CONSUMER_SERVICE):$@ -f- . -f ./Dockerfile/Consumer.Dockerfile
.PHONY: build-consumer-image
