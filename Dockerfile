FROM golang:1.16-alpine

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

# Fix timezone
RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime
RUN echo "Asia/Jakarta" >  /etc/timezone
RUN date

COPY . .
RUN go get -d -v


# Build the Go app linux
# RUN go build -o ./out/go-sample-app .

# Build the Go app mac
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o  ./out/go-sample-app .


# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./out/go-sample-app"]