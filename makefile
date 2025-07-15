APP_NAME = habit-tracker
APP_VERSION = 0.1.0
.PHONY: all linux mac mac-arm64 windows clean
all: linux mac mac-arm64 windows
linux:
	GOOS=linux GOARCH=amd64 go build -o $(APP_NAME)-$(APP_VERSION)-linux
mac:
	GOOS=darwin GOARCH=amd64 go build -o $(APP_NAME)-$(APP_VERSION)-mac
mac-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(APP_NAME)-$(APP_VERSION)-mac-arm64
windows:
	GOOS=windows GOARCH=amd64 go build -o $(APP_NAME)-$(APP_VERSION).exe
clean:
	rm -f $(APP_NAME)-$(APP_VERSION)-linux $(APP_NAME)-$(APP_VERSION)-mac $(APP_NAME)-$(APP_VERSION)-mac-arm64 $(APP_NAME)-$(APP_VERSION).exe