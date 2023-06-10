FROM golang as build
WORKDIR /usr/src/app

COPY go.mod go.sum .
RUN go mod download

COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./run


FROM debian as main
ENV GIN_MODE=release

RUN apt update && apt install -y ca-certificates libc6-dev && apt clean && rm -rf /var/lib

RUN mkdir /app
COPY --from=build /usr/src/app/run /app/run
EXPOSE 8080

CMD ["/app/run"]
