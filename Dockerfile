FROM golang:1.20.4 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd cmd/multi-output-data-processor \
  && go build \
  && chmod +x multi-output-data-processor \
  && mv multi-output-data-processor /app

EXPOSE 8080

ENTRYPOINT ["/app/multi-output-data-processor"]
