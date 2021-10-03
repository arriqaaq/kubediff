IMAGE_REPO=docker.io/arriqaaq/kubediff
TAG=$(shell cut -d'=' -f2- .release)

.DEFAULT_GOAL := build
.PHONY: release git-tag check-git-status build image tag-image publish

release: check-git-status image tag-image publish git-tag
	@echo "Successfully released version $(TAG)"

git-tag:
	@echo "Creating a git tag"
	@git add .release helm/kubediff hack/deploy.yaml CHANGELOG.md
	@git commit -m "Release $(TAG)" ;
	@git tag ${TAG} ;
	@git push --tags;
	@echo 'Git tag pushed successfully' ;

check-git-status:
	@echo "Checking git status"
	@if [ -n "$(shell git tag | grep $(TAG))" ] ; then echo 'ERROR: Tag already exists' && exit 1 ; fi

build:
	GOOS_VAL=$(shell go env GOOS) GOARCH_VAL=$(shell go env GOARCH) go build -o $(shell go env GOPATH)/bin/kubediff
	@echo "Build complete"

image:
	@echo "Building docker image"
	@docker build -t $(IMAGE_REPO) -f Dockerfile --no-cache .
	@echo "Docker image built"

tag-image:
	@echo 'Tagging image'
	@docker tag $(IMAGE_REPO) $(IMAGE_REPO):$(TAG)

publish:
	@echo "Pushing docker image to repository"
	@docker login
	@docker push $(IMAGE_REPO):$(TAG)
	@docker push $(IMAGE_REPO):latest
