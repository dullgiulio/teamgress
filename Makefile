PKG=github.com/dullgiulio/teamgress
BINDIR=bin
BINS=tg-agent-deploy-log \
	 tg-agent-test \
	 tg-controller-term \
	 tg-controller-websocket
PKGDEPS=github.com/ActiveState/tail \
		golang.org/x/net/websocket
#	github.com/GeertJohan/go.rice \
#	github.com/elazarl/go-bindata-assetfs \
#	github.com/jteeuwen/go-bindata

all: clean vet fmt build

build: libteamgress $(BINS)

clean:
	rm -f $(BINDIR)/*

deps: $(PKGDEPS)

fmt:
	go fmt $(PKG)/...

vet:
	go vet $(PKG)/...

libteamgress:
	go build $(PKG)

$(BINS):
	go build -o $(BINDIR)/$@ $(PKG)/$@

$(PKGDEPS):
	go get -u $@

.PHONY: all deps build clean fmt vet libteamgress $(BINS) $(PKGDEPS)
