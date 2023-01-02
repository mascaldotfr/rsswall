all: rsswall

rsswall: *.go go.mod views/*
	go build -o rsswall

clean:
	rm -f rsswall

test:
	go run . feeds.example > /dev/shm/test.html

stresstest: all
	curl -s https://raw.githubusercontent.com/simevidas/web-dev-feeds/master/feeds.opml | \
		xml2 | awk -F= '/@xmlUrl/ {print $$2}' > /dev/shm/big_feed_list
	./rsswall /dev/shm/big_feed_list > /dev/shm/big.html

install: all
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp rsswall $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/rsswall

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/rsswall

.PHONY: all clean stresstest install uninstall
