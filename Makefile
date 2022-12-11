appname := formatyaya

sources := $(wildcard *.go)

build = GOOS=$(1) GOARCH=$(2) go build -o ./out/$(appname)_$(1)_$(2)$(3)

.PHONY: all windows linux

all: windows linux 

linux: installer_linux
installer_linux: $(sources)
	$(call build,linux,amd64,)
	$(call build,linux,arm,)
	$(call build,linux,arm64,)

windows: installer_windows
installer_windows: $(sources)
	$(call build,windows,amd64,.exe)
	$(call build,windows,arm,.exe)
	$(call build,windows,arm64,.exe)
