EXEC = flwlnet-pat

.PHONY: clean tools lint format

all: bin

clean:
	rm -rf $(EXEC)* vendor
	go clean

tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/daixiang0/gci@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2
	go install fyne.io/tools/cmd/fyne@latest

lint:
	golangci-lint run

format:
	gci write --skip-vendor --skip-generated --section standard --section default --section dot .
	find . -name '*.go' -not -path './vendor/*' -exec goimports -w {} \+

bin:
	go build -trimpath -v -o $(EXEC)

release:
	fyne build --release --output $(EXEC)-$(shell go env GOOS)-$(shell go env GOARCH)

podman-init-%:
	$(eval parts := $(subst -, ,$@))
	$(eval os := $(word 3,$(parts)))
	$(eval arch := $(word 4,$(parts)))
	podman build --platform $(os)/$(arch) --tag gobuilder-$(os)-$(arch) --file builder.docker

podman-build-%:
	$(eval parts := $(subst -, ,$@))
	$(eval os := $(word 3,$(parts)))
	$(eval arch := $(word 4,$(parts)))
	go mod vendor
	podman run --volume ${CURDIR}:/build --network none --tty --interactive --rm --platform $(os)/$(arch) gobuilder-$(os)-$(arch) make -C /build release
