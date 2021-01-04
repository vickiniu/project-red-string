FROM golang:latest
RUN mkdir /app
ADD . /app
WORKDIR /app
## Add this go mod download command to pull in any dependencies
RUN go mod download
## Our start command which kicks off
## our newly created binary executable
CMD ["go", "run", "main.go"]