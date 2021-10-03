IMAGE_REPO=docker.io/arriqaaq/kubediff
TAG=$(shell cut -d'=' -f2- .release)

.DEFAULT_GOAL := build
.PHONY: release git-tag check-git-status build image pre-build tag-image publish system-check

release: check-git-status image tag-image publish git-tag
	@echo "Successfully releeased version $(TAG)"

# Create a git tag
git-tag:
	@echo "Creating a git tag"
	@git add .release helm/kubediff deploy-all-in-one.yaml
	@git commit -m "Release $(TAG)" ;
	@git tag $(TAG) ;
	@git push --tags origin develop;
	@echo 'Git tag pushed successfully' ;

# Check git status
check-git-status:
	@echo "Checking git status"
	@if [ -n "$(shell git tag | grep $(TAG))" ] ; then echo 'ERROR: Tag already exists' && exit 1 ; fi
	@if [ -z "$(shell git remote -v)" ] ; then echo 'ERROR: No remote to push tags to' && exit 1 ; fi
	@if [ -z "$(shell git config user.email)" ] ; then echo 'ERROR: Unable to detect git credentials' && exit 1 ; fi

# Build the binary
build: pre-build
	GOOS_VAL=$(shell go env GOOS) GOARCH_VAL=$(shell go env GOARCH) go build -o $(shell go env GOPATH)/bin/kubediff
	@echo "Build completed successfully"

# Build the image
image: pre-build
	@echo "Building docker image"
	@docker build --build-arg GOOS_VAL=linux --build-arg GOARCH_VAL=$(shell go env GOARCH) -t $(IMAGE_REPO) -f Dockerfile --no-cache .
	@echo "Docker image build successfully"

# system checks
system-check:
	@echo "Checking system information"
	@if [ -z "$(shell go env GOOS)" ] || [ -z "$(shell go env GOARCH)" ] ; \
	then \
	echo 'ERROR: Could not determine the system architecture.' && exit 1 ; \
	else \
	echo 'GOOS: $(shell go env GOOS)' ; \
	echo 'GOARCH: $(shell go env GOARCH)' ; \
	echo 'System information checks passed.'; \
	fi ;

# Pre-build checks
pre-build: system-check

# Tag images
tag-image:
	@echo 'Tagging image'
	@docker tag $(IMAGE_REPO) $(IMAGE_REPO):$(TAG)

# Docker push image
publish:
	@echo "Pushing docker image to repository"
	@docker login
	@docker push $(IMAGE_REPO):$(TAG)
	@docker push $(IMAGE_REPO):latest
