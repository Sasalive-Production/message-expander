FROM golang:1.23.3-alpine AS go

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /build

CMD ["/build"]
