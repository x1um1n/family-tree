FROM golang:1.12.0 as build

WORKDIR /go/src/github.com/x1um1n/family-tree
COPY . .

## Get 3rd party golang packages
RUN go get gopkg.in/yaml.v2
RUN go get github.com/x1um1n/checkerr

RUN go build family-tree.go

##RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo family-tree.go

##FROM archlinux/base as runtime

##COPY --from=build go/src/github.com/x1um1n/family-tree/family-tree /usr/local/family-tree/app
##COPY --from=build go/src/github.com/x1um1n/family-tree/web /usr/local/family-tree/web

EXPOSE 8080
CMD ["/go/src/github.com/x1um1n/family-tree/family-tree"]
