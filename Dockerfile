FROM golang:1.19.2 as build-env
#RUN go install github.com/Ether-Security/leviathan@latest
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN go build -o leviathan

FROM ubuntu:22.10

COPY --from=build-env /app/leviathan /usr/bin/leviathan

ENTRYPOINT ["leviathan"]