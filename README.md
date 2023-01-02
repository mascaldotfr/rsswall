# rsswall
*A minimal Wall of Text feed aggregator*

**rsswall** fetches a list of RSS/ATOM feeds specified in input and generates
a static web page from them to stdout as output, with the different feeds being
like post-its on a wallboard. It's similar in spirit to
[upstract/popurls](https://upstract.com/?ref=pop), but it's leaner, fully customisable
and open source.

It has been meant to be hosted on a pre-existing webserver.

You can see a sample generated page at https://mascaldotfr.github.io/rsswall/

## Features

Note that rsswall is considered feature complete because i love its simplicity,
it's so simple that you better fork it for your own customisations. Bugs and
naivety fixes are welcome though.

- Fully progressive web page, it will try to occupy most of the screen size given
  that CSS flex is supported, may it be a small smartphone or a 4K screen. It does
  also work on a text-based browser, in a more linear way indeed.
- Fetch and parse feeds asynchronously
- Single HTML file output.
- No javascript or images, just text. But it has light/dark themes out of the
  box, and the page refreshes itself every 10 minutes.
- Can display only `x` items, per feed, if needed
- One shot program, no service or database to deal with. All resources
  (templates, favicon etc.) are self contained in the generated binary,
  excepted the feedlist, no dependency hell.

## Quickstart and usage

**rsswall** only requires `golang` *(1.19 or later, well, 1.16 should do the job
but has not been tested)* and `make`:

```shell
$ git clone https://github.com/mascaldotfr/rsswall.git
$ cd rsswall
$ make
$ ./rsswall feeds.example > /tmp/rsswall.html
```

For the next runs it's just:

```shell
$ ./rsswall feeds.example > /tmp/rsswall.html
```

**rsswall** operates silently and will only report error/warnings.

## Feedlist

**rsswall** uses a simple feedlist format, similar to `newsboat` (there are
compatible actually as i'm writing this), but with an extension.

I recommend reading the provided `feeds.example`, but here are the rules:

- Blank lines are ignored, whitespaces are trimmed
- All lines starting with `#` are ignored
- All other lines will be treated as a feed URL *line*
- If whitespace(s) and a `number` are present after an URL, **rsswall** will
  only display the last `number` items from that feed.
- If there is no `number` or if it is invalid, then the last 5 (by default)
  items from that feed will be displayed

## Install

Note that there are install/uninstall `make` targets for convenience if you want a
system-wide install.

## Cronjob

**rsswall** needs to be run periodically through cron (or via a systemd timer),
to do so, as user:

```shell
$ crontab -e
```

And add the following line, here we'll run **rsswall** every 30 minutes:
```crontab
*/30 * * * * /where/is/rsswall /where/is/your/feedlist.txt > /where/to/put/the/html/file.html
```

## Customisation

Because **rsswall** is self contained, any change i'm detailing here will
indeed require to rebuild the binary.

### The "views" directory

There you can find the HTML template used, `layout.html`. There is also the
favicon as `favicon.png`. Note that if you want to change the favicon it *has*
to be a PNG file unless you modify `layout.html`'s and `main.go`.

Removing totally the favicon "feature" will require to modify `main.go` as well.

### main.go

While the default settings are pretty sensible, you may want to change some of
them. There are some `const`s in `main.go`, that are self documented, notably
the default timezone used (let's say you're in Europe but your server is in the
USA), time and date formats, and the default number of feed items to display.

## BUGS

This readme is almost as long as all the Go code used for this project.

More seriously, feel free to report them.
