package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/sufiyanshaikh01/gRPC-Remote-Procedure-Calls/gen" // <-- update with your module path
	"google.golang.org/grpc"
)

// URL model (similar to your old struct)
type URL struct {
	ID           string
	OriginalURL  string
	ShortURL     string
	CreationDate time.Time
}

// In-memory DB
type urlStore struct {
	mu sync.RWMutex
	db map[string]URL
}

func newURLStore() *urlStore {
	return &urlStore{db: make(map[string]URL)}
}

func generateShortURL(original string) string {
	hasher := md5.New()
	hasher.Write([]byte(original))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}

func (s *urlStore) create(original string) URL {
	short := generateShortURL(original)
	id := short
	s.mu.Lock()
	defer s.mu.Unlock()
	url := URL{
		ID:           id,
		OriginalURL:  original,
		ShortURL:     short,
		CreationDate: time.Now(),
	}
	s.db[id] = url
	return url
}

func (s *urlStore) get(id string) (URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url, ok := s.db[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

// gRPC server
type urlShortenerServer struct {
	pb.UnimplementedURLShortenerServer
	store *urlStore
}

func (s *urlShortenerServer) Shorten(ctx context.Context, req *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	url := s.store.create(req.OriginalUrl)
	return &pb.ShortenResponse{
		Url: &pb.URL{
			Id:           url.ID,
			OriginalUrl:  url.OriginalURL,
			ShortUrl:     url.ShortURL,
			CreationUnix: url.CreationDate.Unix(),
		},
	}, nil
}

func (s *urlShortenerServer) Resolve(ctx context.Context, req *pb.ResolveRequest) (*pb.ResolveResponse, error) {
	url, err := s.store.get(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ResolveResponse{
		Url: &pb.URL{
			Id:           url.ID,
			OriginalUrl:  url.OriginalURL,
			ShortUrl:     url.ShortURL,
			CreationUnix: url.CreationDate.Unix(),
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterURLShortenerServer(grpcServer, &urlShortenerServer{store: newURLStore()})

	fmt.Println("gRPC server running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
