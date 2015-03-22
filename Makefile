PKG=github.com/dullgiulio/teamgress
BINDIR=bin
BINS=tg-agent-deploy-log \
	 tg-agent-test \
	 tg-controller-term \
	 tg-controller-websocket

all: clean vet fmt build

build: libteamgress $(BINS)

clean:
	rm -f $(BINDIR)/*

fmt:
	go fmt $(PKG)/...

vet:
	go vet $(PKG)/...

libteamgress:
	go build $(PKG)/libteamgress

$(BINS):
	go build -o $(BINDIR)/$@ $(PKG)/$@

.PHONY: all build clean fmt vet libteamgress $(BINS)
