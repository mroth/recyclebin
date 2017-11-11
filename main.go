package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

// environment variables for Twitter Streaming API access tokens
var (
	consumerKey    = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
	accessToken    = os.Getenv("ACCESS_TOKEN")
	accessSecret   = os.Getenv("ACCESS_TOKEN_SECRET")
)

// track some global state for status logging (note: not threadsafe!)
var (
	startTime   time.Time
	tracked     = 0
	skipped     = 0
	trackedLast = 0
	skippedLast = 0
)

func initTwitterAPI() *anaconda.TwitterApi {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessSecret)
	return api
}

func initStreamFilter(terms string) *anaconda.Stream {
	api := initTwitterAPI()
	v := url.Values{}
	v.Set("track", terms)
	v.Set("stall_warnings", "true")
	return api.PublicStreamFilter(v)
}

func tracker(terms string, duration time.Duration) (r results) {
	// init somewhere to store results
	r = newResults()

	// setup streaming from the API
	stream := initStreamFilter(terms)

	// a timer for knowing when we are done sampling
	startTime = time.Now()
	done := time.NewTimer(duration)

	// but also capture ctrl-c interrupts, in case we wants to end early,
	// we should still see the summary of what we already analyzed
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	for {
		select {
		case <-done.C:
			stream.Stop()
			return
		case <-sigchan:
			stream.Stop()
			return
		case m := <-stream.C:
			switch m.(type) {
			case anaconda.Tweet:
				tracked++
				t := m.(anaconda.Tweet)
				// as a quick experiment, try to normalize to text without URL,
				// because t.co fucks with us otherwise, for now just grab text up to before http
				if len(t.Entities.Urls) >= 1 {
					// firstUrl := t.Entities.Urls[0] // this is too unreliable, because of difference in counting multibyte
					mi := strings.Index(t.Text, "http")
					part1 := t.Text[:mi]
					r.phrases.Increment(part1)
				} else {
					r.phrases.Increment(t.Text)
				}

				r.users.Increment(t.User.ScreenName)
				r.lang.Increment(t.Lang)

				for _, url := range t.Entities.Urls {
					r.urls.Increment(url.Expanded_url)
				}
			case anaconda.StallWarning:
				fmt.Println("Got a stall warning! falling behind!")
			case anaconda.DisconnectMessage:
				fmt.Println("Got disconnected!")
			default:
				fmt.Printf("got something else! %T\n", m)
				os.Exit(1)
			}
		}
	}
}

func startLogger(rf time.Duration) *time.Ticker {
	ticker := time.NewTicker(rf)
	go func() {
		for {
			<-ticker.C
			period := tracked - trackedLast
			periodRate := float64(period) / rf.Seconds()
			log.Printf("Tweets tracked: %5d (↑%v, +%.1f/sec.)\n", tracked, period, periodRate)
			trackedLast = tracked
		}
	}()
	return ticker
}

func main() {
	// default flags
	var term = flag.String("term", "♻️", "term to monitor")
	var sampleDuration = flag.Duration("sample", time.Minute*5, "sample length")
	var reportFrequency = flag.Duration("report", 0, "periodically report on progress")
	flag.Parse()

	// start progress reports if desired
	var logger *time.Ticker
	if *reportFrequency > 0 {
		logger = startLogger(*reportFrequency)
	}

	// do the monitoring for the length of sample
	fmt.Printf("🚀 Starting to monitor Twitter for term: [ %v ]...\n", *term)
	results := tracker(*term, *sampleDuration)

	// monitoring is done, cleanup
	if *reportFrequency > 0 {
		logger.Stop()
		time.Sleep(time.Millisecond * 250) // allow logger time to finish any flush to stdout
	}

	// produce the report!
	reportDuration := time.Since(startTime).Truncate(time.Second)
	rate := float64(tracked) / (reportDuration).Seconds()
	fmt.Printf("\n\n ✨ DONE ✨ - time monitored: %v, total tweets tracked: %v, rate: %.1f/sec.\n", reportDuration, tracked, rate)
	results.PrintReport()
}
