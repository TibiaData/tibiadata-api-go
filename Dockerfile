# get golang container
FROM golang:1.22.4 AS builder

# get args
ARG TibiaDataBuildBuilder=dockerfile
ARG TibiaDataBuildRelease=-
ARG TibiaDataBuildCommit=-

# create and set workingfolder
WORKDIR /go/src/

# copy go mod files
COPY go.mod go.sum ./

# copy all sourcecode
COPY src/ ./src/

# download go mods
RUN go mod download

# compile the program
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s -X 'main.TibiaDataBuildBuilder=${TibiaDataBuildBuilder}' -X 'main.TibiaDataBuildRelease=${TibiaDataBuildRelease}' -X 'main.TibiaDataBuildCommit=${TibiaDataBuildCommit}'" -o app ./...


# get alpine container
FROM alpine:3.20.1 AS app

# create workdir
WORKDIR /opt/app

# add ca-certificates and tzdata
RUN apk --no-cache add ca-certificates tzdata

# create nonroot user and group
RUN addgroup -S nonroot && \
  adduser -S nonroot -G nonroot && \
  chown -R nonroot:nonroot .

# set user to nonroot
USER nonroot:nonroot

# copy binary from builder
COPY --from=builder --chown=nonroot:nonroot --chmod=544 /go/src/app .

# expose port 8080
EXPOSE 8080

# run application
CMD ["./app"]
