package main

import (
	"bufio"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"os"
	"sort"
	"strings"
	"strconv"
	"time"

	"github.com/mmcdole/gofeed"
)

// Default number of feed items displayed when none are present in the feedlist
const max_items = 5
// Default date/time format, see the layout at https://pkg.go.dev/time#pkg-constants
const datetimeformat = "02/01 15:04"
// Same but time only
const timeformat = "15:04"
// Your timezone (if different from the machine running rsswall)
// Use `timedatectl` or `cat /etc/timezone` on your own machine
// to get it.
const mytz = "Local"
// Generated page refresh time (in seconds)
const refreshtime int = 10 * 60

// configuration is ending here

//go:embed views/*.html
//go:embed views/favicon.png
var views embed.FS

type Page struct {
	Updated 	string
	RefreshTime	int
	Favicon		string
	Feeds 		[]Feed
}
type Feed struct {
	Title		string
	Items		[]Item
	Link		string
}
type Item struct {
	Datetime 	string
	Link		string
	Title		string
}

func main() {
	var rsspage Page

	tzlocation, err := time.LoadLocation(mytz)
	if err != nil {
		log.Println("Wrong timezone, defaulting to Local")
		tzlocation, _ = time.LoadLocation("Local")
	}

	now := time.Now().In(tzlocation)
	rsspage.Updated = now.Format(timeformat)
	rsspage.RefreshTime = refreshtime

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s feedfile\n", os.Args[0])
		os.Exit(1)
	}
	feedfile, err := os.Open(os.Args[1])
	if err != nil {
        	panic(err)
    	}
	defer feedfile.Close()

	favicon, err := views.ReadFile("views/favicon.png")
	rsspage.Favicon = base64.StdEncoding.EncodeToString(favicon)

	linescanner := bufio.NewScanner(feedfile)
	linescanner.Split(bufio.ScanLines)
	
	fp := gofeed.NewParser()
	for linescanner.Scan() {
		feedline := strings.TrimSpace(linescanner.Text())
		if len(feedline) == 0 || strings.HasPrefix(feedline, "#") {
			continue
		}
		
		words := strings.Fields(feedline)
		url := words[0]
		maxitems := max_items
		if len(words) > 1 {
			maxitems, err = strconv.Atoi(words[1])
			if err != nil {
				maxitems = max_items
				log.Printf(`Feedline parsing failing "%s" Reason: %s.`,
					feedline, err.Error())
				log.Printf("Corrected, using %d as the maxitems value.",
					max_items)
			}
		}
		
		feed, err := fp.ParseURL(url)
		if err != nil {
			log.Println(url + " => " + err.Error())
			continue
		}
		items := feed.Items
		// some feed items have no pubdate, generate one for them
		if items[0].PublishedParsed == nil {
			for i, _ := range items {
				items[i].PublishedParsed = &now
			}
		} else {
			sort.SliceStable(items, func(i, j int) bool {
				return items[i].PublishedParsed.Unix() > items[j].PublishedParsed.Unix()
			})
		}

		var parsedfeed Feed
		parsedfeed.Title = feed.Title
		parsedfeed.Link = feed.Link
		for itemnr, item := range items {
			var parseditem Item
			if (itemnr >= maxitems) {
				break
			}
			parseditem.Link = item.Link
			parseditem.Datetime = item.PublishedParsed.In(tzlocation).Format(datetimeformat)
			parseditem.Title = item.Title
			parsedfeed.Items = append(parsedfeed.Items, parseditem)
		}
		rsspage.Feeds = append(rsspage.Feeds, parsedfeed)
	}
	tmpl, err := template.ParseFS(views, "views/layout.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(os.Stdout, rsspage)
}
