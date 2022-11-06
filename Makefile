VENTI_VERSION=v0.1.13
IMAGE_REPO=ghcr.io/kuoss

LDFLAGS += -X "main.ventiVersion=$(VENTI_VERSION)"
MAKEFLAGS += -j2

# dev server(5173)
dev: go-dev web-dev
go-dev:
	API_ONLY=1 VENTI_VERSION=${VENTI_VERSION} air
web-dev:
	cd web && VITE_SERVER_HMR_HOST=localhost npm run dev --clearScreen=false

# gin server(8080)
watch: go-air web-watch
web-watch:
	cd web && npm run watch
go-air:
	sleep 15 && air

# on codespace
stage:
	skaffold dev --namespace=kube-system --default-repo=ghcr.io/kuoss

#go-build:
#	go mod download -x && go build -ldflags '$(LDFLAGS)' -o /app/venti

docker-build:
	docker build -t ${IMAGE_REPO}/venti:${VENTI_VERSION} --build-arg VENTI_VERSION=${VENTI_VERSION} . && docker push ${IMAGE_REPO}/venti:${VENTI_VERSION} 
