############################
# STEP 1 build the ui
############################
FROM node:14 as builderNode

RUN mkdir /webapp
COPY ./webapp/ /webapp/webapp/
COPY ./build/*.js /webapp/build/
COPY ./*.json /webapp/
COPY ./*.js /webapp/
WORKDIR /webapp
# install node packages
RUN npm set progress=false
RUN npm ci --no-optional
# Build the web app
RUN npm run build
############################
# STEP 2 build executable binary
############################
FROM golang:1.16-alpine as builderGo
ARG BUILD_VERSION
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates
# Create appuser
RUN adduser -D -g '' appuser
# Copy the go source
COPY ./pkg/ $GOPATH/src/github.com/stevenweathers/wakita-retro-tool/pkg/
COPY ./*.go $GOPATH/src/github.com/stevenweathers/wakita-retro-tool/
COPY ./go.mod $GOPATH/src/github.com/stevenweathers/wakita-retro-tool/
COPY ./go.sum $GOPATH/src/github.com/stevenweathers/wakita-retro-tool/
# Copy SQL file
COPY ./schema.sql $GOPATH/src/github.com/stevenweathers/wakita-retro-tool/
# Copy our static assets
COPY --from=builderNode /webapp/dist $GOPATH/src/github.com/stevenweathers/wakita-retro-tool/dist
# Set working dir
WORKDIR $GOPATH/src/github.com/stevenweathers/wakita-retro-tool/
# Fetch dependencies
RUN go mod download
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -ldflags "-X main.version=$BUILD_VERSION" -o /go/bin/wakita
############################
# STEP 3 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=builderGo /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builderGo /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builderGo /go/bin/wakita /go/bin/wakita
# Use an unprivileged user.
USER appuser

# Run the wakita binary.
ENTRYPOINT ["/go/bin/wakita"]