FROM golang:1.17 as build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./server ./src/main.go


# Stage 2: build the image
FROM alpine:latest
RUN apk --no-cache add ca-certificates libc6-compat curl
WORKDIR /app/
COPY --from=build-stage /app/server .
EXPOSE 8080
CMD ["./server"]