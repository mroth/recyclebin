package main

import "fmt"

type results struct {
	phrases, users, urls, lang *termCounter
}

func newResults() results {
	return results{
		phrases: NewTermCounter(),
		users:   NewTermCounter(),
		urls:    NewTermCounter(),
		lang:    NewTermCounter(),
	}
}

func (r results) PrintReport() {
	fmt.Println("\n👨‍👨‍👦 ACCOUNTS 👨‍👨‍👦")
	userScores := r.users.Scores()
	multiTweeters := userScores.GreaterThan(1)
	fmt.Printf("Total distinct accounts: %d, amount who tweeted more than once: %d\n",
		userScores.Len(),
		multiTweeters.Len(),
	)
	fmt.Println("Most active:", userScores.GreaterThan(1).Sorted().First(10))

	fmt.Println("\n📣 LANG 📣")
	langScores := r.lang.Scores()
	fmt.Printf("Language distribution: %v\n", langScores.Sorted())

	fmt.Println("\n🔗 URLS 🔗")
	urlScores := r.urls.Scores()
	reusedUrls := urlScores.GreaterThan(1)
	fmt.Printf("Total distinct URLs: %d, appeared more than once: %d\n", urlScores.Len(), reusedUrls.Len())
	fmt.Println("Most active:", urlScores.GreaterThan(1).Sorted().First(10))

	fmt.Println("\n📃 TEXT 📃")
	phraseScores := r.phrases.Scores()
	reusedPhrases := phraseScores.GreaterThan(1)
	fmt.Printf("Total distinct text phrases (before URL): %d, appeared more than once: %d\n", phraseScores.Len(), reusedPhrases.Len())
	topPhrases := phraseScores.GreaterThan(1).Sorted().First(20)
	fmt.Printf("Top %v most common phrases:\n", len(topPhrases))
	for _, phrase := range topPhrases {
		fmt.Printf("%v: %q\n", phrase.Value, phrase.Key)
	}
}
