VENTI_VERSION=v0.1.10
IMAGE_REPO=ghcr.io/kuoss

LDFLAGS += -X "main.ventiVersion=$(VENTI_VERSION)"
MAKEFLAGS += -j2

web-build:
	cd web && VITE_SERVER_HMR_HOST=localhost npm run build

web-dev:
	cd web && VITE_SERVER_HMR_HOST=localhost npm run dev --clearScreen=false

go-dev:
	API_ONLY=1 VENTI_VERSION=${VENTI_VERSION} air

dev: go-dev web-dev


stage:
	skaffold dev --namespace=kube-system --default-repo=ghcr.io/kuoss

#go-build:
#	go mod download -x && go build -ldflags '$(LDFLAGS)' -o /app/venti

docker-build:
	docker build -t ${IMAGE_REPO}/venti:${VENTI_VERSION} --build-arg VENTI_VERSION=${VENTI_VERSION} . && docker push ${IMAGE_REPO}/venti:${VENTI_VERSION} 
