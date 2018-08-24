# ============================================================
# Dockerfile used to build the various microservices
#
# To build run something like:
#
# docker build -t mytag --build-arg service=darwinref .
#
# where the value of service is:
#   darwinref     The Darwin Reference API
# ============================================================

ARG arch=amd64
ARG goos=linux

# ============================================================
# Build container containing our pre-pulled libraries.
# As this changes rarely it means we can use the cache between
# building each microservice.
FROM golang:alpine as build

# The golang alpine image is missing git so ensure we have additional tools
RUN apk add --no-cache \
      curl \
      git \
      tzdata

# Our build scripts
ADD scripts/ /usr/local/bin/

# go-bindata
RUN go get -v github.com/kevinburke/go-bindata &&\
    go build -o /usr/local/bin/go-bindata github.com/kevinburke/go-bindata/go-bindata

# Ensure we have the libraries - docker will cache these between builds
RUN get.sh

# ============================================================
# source container contains the source as it exists within the
# repository.
FROM build as source
WORKDIR /go/src/github.com/peter-mount/uktransport
ADD . .

# Import sql so we can build as needed
RUN go-bindata -o lib/sqlassets.go -pkg lib sql/

# ============================================================
# Compile the source.
FROM source as compiler
ARG arch
ARG goos
ARG goarch
ARG goarm

# Build the microservice.
# NB: CGO_ENABLED=0 forces a static build
RUN CGO_ENABLED=0 \
    GOOS=${goos} \
    GOARCH=${goarch} \
    GOARM=${goarm} \
    compile.sh /dest

# ============================================================
# Finally build the final runtime container
FROM alpine

# The golang alpine image is missing git so ensure we have additional tools
RUN apk add --no-cache \
      curl \
      tzdata

COPY --from=compiler /dest/ /usr/local/bin/
