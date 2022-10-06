VENTI_VERSION=v0.1.9
IMAGE_REPO=ghcr.io/kuoss

LDFLAGS += -X "main.ventiVersion=$(VENTI_VERSION)"
MAKEFLAGS += -j2

web-dev:
	cd web && npm run dev --clearScreen=false

go-dev:
	API_ONLY=1 VENTI_VERSION=${VENTI_VERSION} air

dev: go-dev web-dev

stage: web-build
	PORT=3000 go run -ldflags '$(LDFLAGS)' .

base2:
	go mod download -x && go build -ldflags '$(LDFLAGS)' -o /app/venti


git-push:
	git add -A; git commit -am ${VENTI_VERSION}; git push

docker-build-base1:
	docker build -t venti:base1 -f Dockerfile.base1_alpine .

docker-build-base2:
	docker build -t venti:base2 -f Dockerfile.base2_golang .

docker-build-base3:
	docker build -t venti:base3 -f Dockerfile.base3_vue .

docker-build:
	docker build -t ${IMAGE_REPO}/venti:${VENTI_VERSION} --build-arg VENTI_VERSION=${VENTI_VERSION} . && docker push ${IMAGE_REPO}/venti:${VENTI_VERSION} 
