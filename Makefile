GO		:= go
SOURCES := $(shell find . -name "*.go")
GOOS	?= linux
GOARCH	?= amd64
TARGET	:= nsdp-$(GOOS)-$(GOARCH)

.PHONY: all clean

all: bin/$(TARGET)

bin/$(TARGET): $(SOURCES)
	-mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o bin/$(TARGET) main.go

clean:
	-rm -rvf bin
