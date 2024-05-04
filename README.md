# Alpaca

Trading exploration done with the [Alpaca API](https://docs.alpaca.markets/), [Reddit API](https://www.reddit.com/dev/api/), and [OpenAI API](https://platform.openai.com/docs/overview).

The script scrapes for the top reddit posts from r/stocks from the previous 24 hours and performs "sentiment analysis" using an LLM. These sentiments are aggregated into positions which are executed directly into the user's Alpaca brokerage. 

Note that this is not financial advice and this strategy should not be used in a non-paper trading scenario.

## To run

To run the script, run:

```bash
go run main/main.go
```

## env

The script reads secrets from a `.env` file (not tracked via Git for obvious reasons), so ensure that you have a `.env` file with the following keys:

```
AlpacaAPIKey=""
AlpacaAPISecret=""
OpenAIAPIKey=""
```