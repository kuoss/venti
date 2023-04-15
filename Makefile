VERSION := v0.1.13
IMAGE := ghcr.io/kuoss/venti:$(VERSION)

install-dev:
	go mod tidy
	cd web && npm install
	go install github.com/cosmtrek/air@latest

mock-prometheus:
	docker rm -f prometheus; docker run -d -p9090:9090 --name prometheus prom/prometheus


# dev server (port 5173)
run-dev: run-dev-go run-dev-web
run-dev-go:
	air
run-dev-web:
	cd web && npm run dev --clearScreen=false

# gin server (port 8080)
run-watch: run-watch-go run-watch-web
run-watch-go:
	cd web && npm run watch
run-watch-web:
	sleep 15 && air

# gin server (port 8080)
run-air:
	cd web && npm run build
	air


docker:
	docker build -t $(IMAGE) --build-arg VERSION=$(VERSION) . && docker push $(IMAGE)

pre-checks:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

checks:
	./hack/checks.sh

test-cover:
	./hack/test-cover.sh

go-licenses:
	./hack/go-licenses.sh
