FROM golang:1.22-alpine as go-builder
WORKDIR /app

COPY server .
RUN mkdir -p ./dist && GO111MODULE=on go mod download
RUN go build -o ./dist/traefik-authserver main.go

FROM node:22-alpine as node-builder
WORKDIR /app

COPY web .
RUN npm i -g pnpm
RUN pnpm i
RUN pnpm build

FROM alpine
WORKDIR /app

# default config file
ENV CONF_FILE_PATH=/etc/traefik-authserver.yaml

COPY --from=go-builder /app/dist/traefik-authserver /app/traefik-authserver
COPY --from=node-builder /app/dist /app/public

CMD ["/app/traefik-authserver"]
