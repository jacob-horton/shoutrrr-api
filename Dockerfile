FROM golang
WORKDIR /usr/src/app

COPY go.mod go.sum .
RUN go mod download

COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./run

EXPOSE 8080

CMD ./run
