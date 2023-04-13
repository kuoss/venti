FROM golang:1.19-alpine AS base1
WORKDIR /temp/
RUN apk add --no-cache git npm make gcc musl-dev
COPY . ./
RUN go mod download -x
RUN go build -o /app/venti

FROM node:lts-alpine AS base2
COPY --from=base1 /app/venti /app/
WORKDIR /temp/
COPY . ./
RUN cd web && npm install --force
RUN cd web && npm run build
RUN mkdir -p             /app/web/
RUN cp -a /temp/web/dist /app/web/
RUN cp -a ./tools        /app/
RUN mkdir -p             /app/data

FROM alpine:3.17
COPY --from=base2 /app /app
WORKDIR /app
ENTRYPOINT ["/app/venti"]
