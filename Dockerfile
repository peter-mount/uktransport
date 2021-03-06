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
# The base golang environment with curl, git, uptodate tzdata
# and go-bindata installed
FROM golang:1.12.6-alpine as build
RUN apk add --no-cache \
      curl \
      git \
      tzdata \
      zip

# Fetch the go-bindata utility
RUN go get -v github.com/kevinburke/go-bindata &&\
    go build -o /usr/local/bin/go-bindata \
      github.com/kevinburke/go-bindata/go-bindata &&\
    mkdir -p /dest/bin

# ============================================================
# This stage contains the sources.
# It also generates any .go files, e.g. the sql
FROM build as source
WORKDIR /work

# This stage installs the required libraries
ADD go.mod .
RUN go mod download

ADD lib/ lib/
ADD naptanimport/ naptanimport/
ADD nptgimport/ nptgimport/
ADD publishmq/ publishmq/
ADD sql/ sql/

# Import sql so we can build as needed
RUN go-bindata -o lib/sqlassets.go -pkg lib sql/

ADD dbrest/ dbrest/

# ============================================================
# Now compile our binaries
FROM source as compiler
ARG arch
ARG goos
ARG goarch
ARG goarm

# Build the microservice.
# NB: CGO_ENABLED=0 forces a static build
RUN for bin in dbrest naptanimport nptgimport publishmq; \
    do \
      echo "Building ${bin}"; \
      CGO_ENABLED=0 \
      GOOS=${goos} \
      GOARCH=${goarch} \
      GOARM=${goarm} \
      go build -o /dest/bin/${bin} \
        ./${bin}/bin; \
    done

# ============================================================
# This stage retrieves prebuilt binaries from other containers
# that we want to include in this image
FROM compiler as bins

# cifimport
COPY --from=area51/nrod-cif:latest /bin/cifimport /dest/bin/
COPY --from=area51/nrod-cif:latest /bin/cifretrieve /dest/bin/
COPY --from=area51/dataretriever:latest /usr/local/bin/dataretriever /dest/bin/

# ============================================================
# Optional stage, upload the binaries as a tar file
FROM bins AS upload
ARG uploadPath=
ARG uploadCred=
ARG uploadName=
RUN if [ -n "${uploadCred}" -a -n "${uploadPath}" -a -n "${uploadName}" ] ;\
    then \
      cd /dest/bin; \
      tar cvzpf /tmp/${uploadName}.tgz * && \
      zip /tmp/${uploadName}.zip * && \
      curl -u ${uploadCred} --upload-file /tmp/${uploadName}.tgz ${uploadPath}/ && \
      curl -u ${uploadCred} --upload-file /tmp/${uploadName}.zip ${uploadPath}/; \
    fi

# ============================================================
# Finally build the final runtime container
FROM alpine

RUN apk add --no-cache \
      curl \
      tzdata

COPY --from=bins /dest/ /usr/local/
