GOSRC = $(shell find . -type f -name '*.go')

version=v0.0.1

build: dhcp_client4

dhcp_client4: $(GOSRC)
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dhcp_client4 cmd/dhcp4/dhcp4.go

dhcp4-arm:
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o dhcp_client4 cmd/dhcp4/dhcp4.go

dhcp_client6:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dhcp_client6 cmd/dhcp6/dhcp6.go

clean4:
	rm -rf dhcp_client4

clean6:
	rm -rf dhcp_client6

.PHONY: clean dhcp4