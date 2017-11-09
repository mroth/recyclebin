package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
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

type results struct {
	phrases, users, urls, lang *termCounter
}

func tracker(terms string, duration time.Duration) (r results) {
	r.phrases = NewTermCounter()
	r.users = NewTermCounter()
	r.urls = NewTermCounter()
	r.lang = NewTermCounter()

	api := initTwitterAPI()
	v := url.Values{}
	v.Set("track", terms)
	v.Set("stall_warnings", "true")
	stream := api.PublicStreamFilter(v)

	done := time.NewTimer(duration)
	for {
		select {
		case <-done.C:
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
			default:
				fmt.Println("got something else!")
			}
		}
	}
}

func logger() *time.Ticker {
	var lograte = time.Second * 5
	ticker := time.NewTicker(lograte)
	go func() {
		for {
			<-ticker.C
			period := tracked - trackedLast
			periodRate := float64(period) / lograte.Seconds()
			log.Printf("Tweets tracked: %v (↑%v, +%v/sec.)\n", tracked, period, periodRate)
			trackedLast = tracked
		}
	}()
	return ticker
}

func main() {
	term := "♻️"
	fmt.Printf("🚀 Starting to monitor Twitter for term: [ %v ]...\n", term)

	logger := logger()
	duration := time.Second * 10 // TODO: paramaterize
	results := tracker("♻️", duration)

	logger.Stop()
	time.Sleep(time.Millisecond * 250) // allow logger to catch up

	rate := float64(tracked) / duration.Seconds()
	fmt.Printf("\n\n✨ DONE ✨ - time monitored: %v, total tweets tracked: %v, rate: %.1f/sec.\n", duration, tracked, rate)

	fmt.Println("\nACCOUNTS")
	userScores := results.users.Scores()
	multiTweeters := userScores.GreaterThan(1)
	fmt.Printf("Total distinct accounts: %d, amount who tweeted more than once: %d\n",
		userScores.Len(),
		multiTweeters.Len(),
	)
	fmt.Println("Most active:", userScores.GreaterThan(1).Sorted().First(10))

	fmt.Println("\nLANGUAGE")
	langScores := results.lang.Scores()
	fmt.Printf("Language distribution: %v\n", langScores.Sorted())

	fmt.Println("\nURLS")
	urlScores := results.urls.Scores()
	reusedUrls := urlScores.GreaterThan(1)
	fmt.Printf("Total distinct URLs: %d, appeared more than once: %d\n", urlScores.Len(), reusedUrls.Len())
	fmt.Println("Most active:", urlScores.GreaterThan(1).Sorted().First(10))

	fmt.Println("\nTEXT")
	phraseScores := results.phrases.Scores()
	reusedPhrases := phraseScores.GreaterThan(1)
	fmt.Printf("Total distinct text tweets: %d, appeared more than once: %d\n", phraseScores.Len(), reusedPhrases.Len())
	fmt.Println("Most common:", phraseScores.GreaterThan(1).Sorted().First(10))

}

// TODO: catch early interrupt and show results
