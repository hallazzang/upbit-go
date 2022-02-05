package main

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/hallazzang/upbit-go/upbit"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	opts := upbit.NewClientOptions()
	if secretKey != "" {
		opts.SetSecretKey(secretKey)
	}
	c, err := upbit.NewClient(accessKey, opts)
	if err != nil {
		panic(err)
	}

	resp, err := c.GET("candles/minutes/1", url.Values{"market": {"KRW-BTC"}, "count": {"1"}})
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(content))
}
