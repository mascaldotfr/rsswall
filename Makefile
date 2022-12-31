all: rsswall

rsswall: *.go go.mod views/*
	go build -o rsswall

clean:
	rm -f rsswall

test: 
	go run . feeds.example > /dev/shm/test.html

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp rsswall $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/rsswall

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/rsswall

.PHONY: all clean install uninstall
