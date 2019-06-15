TARGET=cloud-ocr

all: clean build

clean:
	rm -rf $(TARGET)

build:
	go build && go install