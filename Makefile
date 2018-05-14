VERSION=$(shell cat ./VERSION)
COMMIT=$(shell git rev-parse --short HEAD)
LATEST_TAG=$(shell git tag -l | head -n 1)

export VERSION COMMIT LATEST_TAG
.PHONY: test

test:
	@echo "=> Running tests"
	./hack/run-tests.sh

build:
	./hack/cross-platform-build.sh

verify:
	./hack/verify-version.sh

container: build
	docker build -t quay.io/vastness/vcs-webhook:${COMMIT} .

push: container
	docker push quay.io/vastness/vcs-webhook:${COMMIT}
	docker tag quay.io/vastness/vcs-webhook:${COMMIT} quay.io/vastness/vcs-webhook:${VERSION}
	docker push quay.io/vastness/vcs-webhook:${VERSION}
	docker tag quay.io/vastness/vcs-webhook:${COMMIT} quay.io/vastness/vcs-webhook:latest
	docker push quay.io/vastness/vcs-webhook:latest
