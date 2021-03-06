.PHONY: test build run

GOLANG_VERSION ?= 1.12
USERNAME ?= patricksanders
REPO_NAME ?= statsdebug
TRAVIS_REPO_SLUG ?= ${USERNAME}/${REPO_NAME}
VERSION ?= 0.2.0
DOCKER_TAG := ${TRAVIS_REPO_SLUG}:${VERSION}
LATEST_TAG := ${TRAVIS_REPO_SLUG}:latest

test:
	@docker run --rm \
		-v $(shell pwd):/opt/${REPO_NAME} \
		-w /opt/${REPO_NAME} \
		golang:${GOLANG_VERSION} \
		go test

build:
	@docker build -t ${TRAVIS_REPO_SLUG} \
		--build-arg GOLANG_VERSION=${GOLANG_VERSION} \
		.
	docker tag ${TRAVIS_REPO_SLUG} ${DOCKER_TAG}

publish:
	@docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}
	@docker push ${DOCKER_TAG}
	@docker push ${LATEST_TAG}

run:
	@docker run --rm -it -p 8080:8080 -p 8125:8125/udp --rm ${DOCKER_TAG}
