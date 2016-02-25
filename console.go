package main

/*
Create a simple console program that monitors HTTP traffic on your machine:

Consume an actively written-to w3c-formatted HTTP access log

Every 10s, display in the console the sections of the web site with the most hits (a section is defined as being what's before the second '/' in a URL. i.e. the section for "http://my.site.com/pages/create' is "http://my.site.com/pages"), as well as interesting summary statistics on the traffic as a whole.

Make sure a user can keep the console app running and monitor traffic on their machine

Whenever total traffic for the past 2 minutes exceeds a certain number on average, add a message saying that “High traffic generated an alert - hits = {value}, triggered at {time}”

Whenever the total traffic drops again below that value on average for the past 2 minutes, add another message detailing when the alert recovered

Make sure all messages showing when alerting thresholds are crossed remain visible on the page for historical reasons.

Write a test for the alerting logic

Explain how you’d improve on this application design

*/

import (
	"container/ring"
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	// implements tail -f on a file
	"github.com/hpcloud/tail"
)

// LogRegexp is the main regexp used to parse a w3c access log file, to extract "sections".
const (
	LogRegexp           = "\"GET (/[^/]+)/"
	HighThresholdAlert  = "High traffic generated an alert - hits = %.2f, triggered at %s"
	UnderThresholdAlert = "High traffic alert recovered, triggered at %s"
)

var Alerts []string

// Used to display the top N section hits
type topN struct {
	section string
	hits    int64
}

// Processes a single line of log. Extract just the request; returns the section.
func processLine(rx *regexp.Regexp, line string) string {
	res := rx.FindStringSubmatch(line)
	if res == nil {
		return ""
	}
	return res[1]
}

func displayInfo(gh int64, top [3]topN, avg float64) {
	fmt.Printf("Total hits: %d, Avg hits over last 2m: %.2f\n", gh, avg)
	for i, h := range top {
		fmt.Printf("Top %d hit: %s, %d hits\n", i+1, h.section, h.hits)
	}
	fmt.Printf("Alerts:\n")
	fmt.Printf(strings.Join(Alerts, "\n"))
	fmt.Printf("\n\n")
}

func averageSample(sample *ring.Ring) float64 {
	var sum int
	sample.Do(func(x interface{}) {
		if x != nil {
			sum += x.(int)
		}
	})
	return float64(sum) / float64(sample.Len())
}

func checkHighTraffic(sample *ring.Ring, thresh int, alertStatus bool) bool {
	trafAvg := averageSample(sample)
	if int(trafAvg) > thresh && !alertStatus {
		alertStatus = true
		Alerts = append(Alerts, fmt.Sprintf(HighThresholdAlert, trafAvg, time.Now().Format(time.UnixDate)))
	}
	if int(trafAvg) <= thresh && alertStatus {
		alertStatus = false
		Alerts = append(Alerts, fmt.Sprintf(UnderThresholdAlert, time.Now().Format(time.UnixDate)))
	}
	return alertStatus
}

func main() {
	Rx := regexp.MustCompile(LogRegexp)

	var SectionCount = struct {
		sync.RWMutex
		m map[string]int64
	}{m: make(map[string]int64)}
	var GlobalCounter = struct {
		sync.RWMutex
		totalHits     int64
		twoSecondHits int
		topThree      [3]topN
	}{}

	// sample hits every 2s for 2m
	hitSample := ring.New(60)
	alertStatus := false

	// Flags
	var fname string
	var trafficThreshold int
	flag.StringVar(&fname, "fname", "access_log", "File name to parse")
	flag.IntVar(&trafficThreshold, "avgThreshold", 100, "Traffic threshold for sending alerts")
	flag.Parse()

	scanner, err := tail.TailFile(fname, tail.Config{ReOpen: true, Follow: true})
	if err != nil {
		log.Fatalf("Could not open file %s: %q", fname, err)
	}

	// Every 2s, update sampling
	go func() {
		for {
			time.Sleep(2 * time.Second)
			GlobalCounter.Lock()
			hitSample.Value = GlobalCounter.twoSecondHits
			GlobalCounter.twoSecondHits = 0
			GlobalCounter.Unlock()
			hitSample = hitSample.Next()
		}
	}()

	// Every 10s, display useful information
	go func() {
		for {
			GlobalCounter.RLock()
			th := GlobalCounter.totalHits
			t3 := GlobalCounter.topThree
			GlobalCounter.RUnlock()
			displayInfo(th, t3, averageSample(hitSample))
			time.Sleep(10 * time.Second)
		}
	}()

	// Every 2 minutes, compare traffic to threshold
	go func() {
		for {
			alertStatus = checkHighTraffic(hitSample, trafficThreshold, alertStatus)
			time.Sleep(2 * time.Minute)
		}
	}()

	for line := range scanner.Lines {
		section := processLine(Rx, line.Text)
		if section != "" {
			SectionCount.Lock()
			SectionCount.m[section]++
			hits := SectionCount.m[section]
			SectionCount.Unlock()
			GlobalCounter.Lock()
		outerloop:
			for i, tt := range GlobalCounter.topThree {
				if hits > tt.hits {
					GlobalCounter.topThree[i].section = section
					GlobalCounter.topThree[i].hits = hits
					break outerloop
				}
			}
			GlobalCounter.totalHits++
			GlobalCounter.twoSecondHits++
			GlobalCounter.Unlock()
		}
	}
}
