FROM golang:1.18
WORKDIR /app
RUN go mod init myapi
COPY *.go ./
RUN go mod tidy
RUN go mod download
RUN go build -o /api
CMD ["/api"]