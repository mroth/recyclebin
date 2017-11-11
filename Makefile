TARGET := recyclebin
.PHONY: clean

$(TARGET): *.go
	go build -o $@

clean:
	rm $(TARGET)
