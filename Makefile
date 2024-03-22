ifeq ($(OS),Windows_NT)
 PLATFORM="win"
else
 ifeq ($(shell uname),Darwin)
  PLATFORM="mac"
 else
  PLATFORM="linux"
 endif
endif

ifeq ($(PREFIX),)
    PREFIX := /usr/local/bin
endif

BINARY=cert-deployer
B_MAC=${BINARY}-darwin
B_LINUX=${BINARY}-linux
B_WIN=${BINARY}-windows.exe

default:
	# outs dir should be synced with gh-actions:release
	GOARCH=amd64 GOOS=darwin go build -o "bin/${B_MAC}" ./cmd/app
	GOARCH=amd64 GOOS=linux go build -o "bin/${B_LINUX}" ./cmd/app
	GOARCH=amd64 GOOS=windows go build -o "bin/${B_WIN}" ./cmd/app

clean:
	@rm -rf bin
	@go clean

install:
	@if [ ${PLATFORM} == "mac" ]; then \
	    cp "bin/${B_MAC}" "${PREFIX}/${BINARY}"; \
	elif [ ${PLATFORM} == "linux" ]; then \
	    cp "bin/${B_LINUX}" "${PREFIX}/${BINARY}"; \
	else \
	  	echo "Only support mac/unix-like system!"; \
	fi

uninstall:
	@if [ ${PLATFORM} == "win" ]; then \
	    echo "Only support mac/unix-like system!"; \
	else \
	  	rm "${PREFIX}/${BINARY}"; \
	fi

help:
	@echo "make build \tbuild binary file"
	@echo "make clean \tclean binary file"
	@echo "make install \tinstall to ${PREFIX}"
	@echo "make uninstall"
