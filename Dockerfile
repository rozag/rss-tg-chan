# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.12
ARG GO_VERSION=1.12

# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user
RUN echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd
RUN echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates

# Set the environment variables for the go command:
# * CGO_ENABLED=0 to build a statically-linked executable
# * GOFLAGS=-mod=vendor to force `go build` to look into the `/vendor` folder.
ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor

# Define app sources path
ENV APP_PATH=/go/src/github.com/rozag/rss-tg-chan

# Assuming the source code is collocated to this Dockerfile
COPY ./config.ini /
COPY . ${APP_PATH}

# Jump to sources
WORKDIR ${APP_PATH}

# Build the executable to `/app`. Mark the build as statically linked.
RUN go build -installsuffix 'static' -o /app .

# The second and final stage
FROM scratch

# Import the user and group files from the first stage
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the first stage
COPY --from=builder /app /app

# Import the config file from the first stage
COPY --from=builder /config.ini /config.ini

# Perform any further action as an unprivileged user
USER nobody:nobody

# Run the compiled binary
ENTRYPOINT ["/app"]