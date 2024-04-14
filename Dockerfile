FROM golang:1.22

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update \
    && apt-get -y install dos2unix \
    && dos2unix wait-for-postgres.sh \
    && chmod +x wait-for-postgres.sh

RUN apt-get -y install postgresql-client

RUN go mod download
RUN go mod tidy
RUN go build -o banner ./cmd/main.go

CMD ["./banner"]