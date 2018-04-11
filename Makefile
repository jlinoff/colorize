# Build the colorize program on mac and linux for 64 bit systems.
# To add additional architectures see
#     https://golang.org/doc/install/source#environment.

all: bin/native/colorize bin/darwin/amd64/colorize bin/linux/amd64/colorize

clean:
	@-find . -type f -name '*~' -delete
	rm -rf bin

help: bin/native/colorize
	bin/native/colorize --help

bin/darwin/amd64/colorize: colorize.go
	GOOS=darwin GOARCH=amd64 go build -buildmode=exe -o $@ colorize.go
	cp $@ $@-darwin-amd64

bin/linux/amd64/colorize: colorize.go
	GOOS=linux GOARCH=amd64 go build -buildmode=exe -o $@ colorize.go
	cp $@ $@-linux-amd64

bin/native/colorize: colorize.go
	go build -buildmode=exe -o $@ colorize.go

