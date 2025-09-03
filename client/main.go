package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/yourname/urlshortener/gen"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect:", err)
	}
	defer conn.Close()

	client := pb.NewURLShortenerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Call Shorten
	resp, err := client.Shorten(ctx, &pb.ShortenRequest{
		OriginalUrl: "https://github.com/sufiyanshaikh01",
	})
	if err != nil {
		log.Fatal("Shorten failed:", err)
	}
	fmt.Printf("Shortened: %s -> %s\n", resp.Url.OriginalUrl, resp.Url.ShortUrl)

	// Call Resolve
	res, err := client.Resolve(ctx, &pb.ResolveRequest{
		Id: resp.Url.Id,
	})
	if err != nil {
		log.Fatal("Resolve failed:", err)
	}
	fmt.Printf("Resolved: %s -> %s\n", res.Url.ShortUrl, res.Url.OriginalUrl)
}
