FROM node:19-alpine as frontend
COPY frontend /app
WORKDIR /app
RUN npm i &&\
    npm run build

FROM golang:1.19-alpine as backend
COPY backend /app
WORKDIR /app
COPY --from=frontend /app/build pb_public
RUN apk update &&\
    apk add gcc musl-dev &&\
    go mod tidy &&\
    go build -ldflags "-s -w" -o upsnap main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=backend /app/upsnap .
RUN apk update &&\
    apk add --no-cache nmap curl samba-common-tools &&\
    rm -rf /var/cache/apk/*
CMD ["./upsnap", "serve", "--http", "0.0.0.0:8090"]
