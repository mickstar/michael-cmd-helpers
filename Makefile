BINARY := michael-cmd
INSTALL_DIR := $(HOME)/.local/bin

.PHONY: build install

build:
	go build -o $(BINARY) .

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "installed to $(INSTALL_DIR)/$(BINARY)"
