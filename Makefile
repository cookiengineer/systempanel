.PHONY: build install uninstall clean run

APP_NAME := systempanel
GO := go
PREFIX := /usr

build:
	$(GO) build -o $(APP_NAME) .

install: build
	install -Dm755 $(APP_NAME) $(DESTDIR)$(PREFIX)/bin/$(APP_NAME)
	install -Dm644 assets/systempanel.desktop $(DESTDIR)$(PREFIX)/share/applications/systempanel.desktop

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/$(APP_NAME)
	rm -f $(DESTDIR)$(PREFIX)/share/applications/systempanel.desktop

clean:
	rm -f $(APP_NAME)

run: build
	./$(APP_NAME)
