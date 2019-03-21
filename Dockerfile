# Base build image
FROM golang:1.11-alpine AS build_base
 
# Install some dependencies needed to build the project
RUN apk add bash git gcc g++ libc-dev

COPY /utils/go.mod /utils/go.sum /chuck/utils/
COPY /storage/go.mod /storage/go.sum /chuck/storage/
COPY /handlers/go.mod /handlers/go.sum /chuck/handlers/
COPY /cmds/go.mod /cmds/go.sum /chuck/cmds/
COPY go.mod go.sum /chuck/

WORKDIR /chuck/
 
# Force the go compiler to use modules
ENV GO111MODULE=on

# This is the ‘magic’ step that will download all the dependencies that are specified in 
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download 
# command will _ only_ be re-run when the go.mod or go.sum file change 
# (or when we add another docker instruction this line)
RUN go mod download
 
# This image builds the chuck
FROM build_base AS server_builder

COPY /utils/ ./utils
COPY /storage/ ./storage
COPY /handlers/ ./handlers
COPY /cmds/ ./cmds
COPY main.go .

# And compile the project
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' 
 
#In this last stage, we start from a fresh Alpine image, to reduce the image size and not ship the Go compiler in our production artifacts.
FROM alpine AS chuck-intg

# Finally we copy the statically compiled Go binary.
COPY --from=server_builder ./chuck /bin/chuck/
COPY ca.pem key.pem /bin/chuck/

EXPOSE 8123

WORKDIR /bin/chuck/
CMD ["./chuck", "intg", "-address=0.0.0.0", "-port=8123", "-folder=intg"]