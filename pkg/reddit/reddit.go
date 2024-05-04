package reddit

import (
	"context"
	"fmt"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Reddit struct {
	Client *reddit.Client
}

func NewReddit() (*Reddit, error) {
	redditClient, err := reddit.NewReadonlyClient()
	if err != nil {
		return nil, fmt.Errorf("error initializing reddit client: %v", err)
	}
	return &Reddit{
		Client: redditClient,
	}, nil
}

func (r Reddit) GetTopSubredditPosts(subreddit string, limit int) ([]reddit.Post, error) {
	posts, resp, err := r.Client.Subreddit.TopPosts(context.Background(), subreddit, &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: limit,
		},
		Time: "day",
	})
	if err != nil {
		return nil, fmt.Errorf("error getting top subreddit posts with response: %v; %v", resp, err)
	}
	return convertPosts(posts), nil
}

func convertPosts(posts []*reddit.Post) []reddit.Post {
	convertedPosts := make([]reddit.Post, len(posts))
	for i, p := range posts {
		convertedPosts[i] = *p
	}
	return convertedPosts
}
