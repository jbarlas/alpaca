package main

import (
	"fmt"
	"log"

	a "github.com/jbarlas/alpaca/pkg/alpaca"
	o "github.com/jbarlas/alpaca/pkg/openai"
	r "github.com/jbarlas/alpaca/pkg/reddit"

	"github.com/joho/godotenv"
	redditApi "github.com/vartanbeno/go-reddit/v2/reddit"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	alpaca, err := a.NewAlpaca()
	if err != nil {
		log.Fatal("Error initializing alpaca: ", err)
	}
	openai, err := o.NewOpenAI()
	if err != nil {
		log.Fatal("Error initializing openai: ", err)
	}
	reddit, err := r.NewReddit()
	if err != nil {
		log.Fatal("Error initializing reddit: ", err)
	}

	posts, err := reddit.GetTopSubredditPosts("stocks", 15)
	if err != nil {
		log.Fatal("Error getting top subreddit posts: ", err)
	}

	ingestPosts(posts, alpaca, openai)

	fmt.Println("-------------------------------------------------")
	fmt.Println("Closing all positions")
	err = alpaca.CloseAllPositions()
	if err != nil {
		log.Fatal("Error closing all positions: ", err)
	}
	fmt.Printf("Taking positions: %v\n", alpaca.Positions)
	alpaca.ExecutePositions()
}

func ingestPosts(posts []redditApi.Post, alpaca *a.Alpaca, openai *o.OpenAI) {
	for _, post := range posts {
		fmt.Println("-------------------------------------------------")
		fmt.Printf("Analyzing post: %s \nURL: %s\n", post.Title, post.URL)
		analysis, err := openai.PerformAnalysisOnPost(post)
		if err != nil {
			fmt.Println("Error performing analysis on post: ", err)
		}
		if analysis == nil || err != nil {
			fmt.Println("skipping this post...")
			continue
		}
		fmt.Println("\nAnalysis for post:")
		fmt.Printf("Security: %s\n", analysis.Security)
		fmt.Printf("Sentiment: %d\n", analysis.Sentiment)
		fmt.Printf("Summary: %s\n", analysis.Summary)
		err = alpaca.AddPosition(analysis.Security, analysis.Sentiment)
		if err != nil {
			fmt.Println("Error adding position to alpaca: ", err)
			continue
		}
		fmt.Printf("\nAdded position %s: %d\n", analysis.Security, analysis.Sentiment)
	}
}
