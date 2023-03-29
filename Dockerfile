# Build stage
FROM golang:1.19.7-alpine3.17 AS builder
WORKDIR /app
# will build from the root of our project / current working directory image
COPY . . 
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.17
WORKDIR /app
# --from -> argument to tell docker where to copy the file from
# /app/main . -> dot represents the WORKDIR that we set above /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
# excecutable
CMD ["/app/main"]
# the entry point of the docker image
ENTRYPOINT [ "/app/start.sh" ]

# need original code
