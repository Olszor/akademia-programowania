package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reddit/fetcher"
	"sync"
	"time"
)

func main() {
	//var f fetcher.RedditFetcher // do not change
	var w io.Writer // do not change

	subreddits := []string{
		"golang",
		"java",
		"csharp",
	}

	var wg sync.WaitGroup
	wg.Add(len(subreddits))

	for _, subreddit := range subreddits {
		go func(sub string) {
			defer wg.Done()
			f := &fetcher.Fetcher{
				Client: http.Client{},
				Url:    fmt.Sprintf("https://www.reddit.com/r/%s.json", sub),
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			err := f.FetchWithContext(ctx)
			if err != nil {
				log.Printf(err.Error())
			}

			w, err = os.Create(fmt.Sprintf("./data/%s.txt", sub))
			if err != nil {
				log.Printf(err.Error())
			}

			err = f.Save(w)
			if err != nil {
				log.Printf(err.Error())
			}
		}(subreddit)
	}
	wg.Wait()

}
