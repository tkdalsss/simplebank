# Build stage
FROM golang:1.19.7-alpine3.17 AS builder
WORKDIR /app
# will build from the root of our project / current working directory image
COPY . . 
RUN go build -o main main.go


# Run stage
FROM alpine:3.17
WORKDIR /app
# --from -> argument to tell docker where to copy the file from
# /app/main . -> dot represents the WORKDIR that we set above /app
COPY --from=builder /app/main .


EXPOSE 8080
# excecutable
CMD ["/app/main"]


# need original code
