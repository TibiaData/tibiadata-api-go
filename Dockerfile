# get golang container
FROM golang:1.18.3

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


# get latest alpine container
FROM alpine:latest

# add ca-certificates
RUN apk --no-cache add ca-certificates tzdata

# create workdir
WORKDIR /root/

# copy binary from first container
COPY --from=0 /go/src/app .

# expose port 8080
EXPOSE 8080

# run application
CMD ["./app"]
