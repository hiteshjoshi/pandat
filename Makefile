
BINARY=pandat
VERSION=1.0.0
BUILD=`git rev-parse HEAD`

CGO_ENABLED=0

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}" -a -installsuffix cgo


build: 	
	 GOOS=linux GOARCH=amd64 go build ${LDFLAGS}  -o dist/pandat-linux-amd64

mac: 
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/pandat-darwin-amd64

windows: 
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/pandat-windows-amd64

docker:
	docker build -t pandat .
	
install: 
	go install ${LDFLAGS}

clean: 
	if [-f ${BINARY}] ; then rm ${BINARY}; fi

.PHONY: clean install build