package main

import (
	"context"
	"log"
	"math/rand/v2"
	"sync"

	"github.com/bufbuild/protovalidate-go"
	"github.com/moguchev/grpc/pkg/api/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Post struct {
	ID       uint64
	Title    string
	Content  string
	AuthorID string
}

type ExampleService struct {
	example.UnimplementedExampleServiceServer

	validator *protovalidate.Validator
	storage   map[uint64]*Post
	mx        sync.RWMutex
}

const clientHeaderName = "x-client-header-key"

func (s *ExampleService) CreatePost(ctx context.Context, req *example.CreatePostRequest) (*example.CreatePostResponse, error) {
	// Все прешедшие заголовки будут лежать в контексте
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		// Metadata - это map[Header][]Value
		log.Println(md.Get(clientHeaderName))
	}
	// Вариант 2.
	log.Println(metadata.ValueFromIncomingContext(ctx, clientHeaderName))

	id := rand.Uint64()
	post := &Post{
		ID:       id,
		Title:    req.GetTitle(),
		Content:  req.GetContent(),
		AuthorID: req.GetAuthorId(),
	}
	s.mx.Lock()
	s.storage[id] = post
	s.mx.Unlock()

	// Мы можем также создать свои заголовки
	header := metadata.Pairs(
		"x-header-key1", "value1",
		"x-header-key2", "value2",
	)

	// и отправить их клиенту на его запрос либо в заголовках либо в конце (Trailer)
	grpc.SetHeader(ctx, header)  //  регистрирует заголовки (Header), которые будут отправлены в начале ответа до основного тела.
	grpc.SetTrailer(ctx, header) //  регистрирует трейлеры (Trailer), которые будут отправлены в конце ответа после основного тела.

	return &example.CreatePostResponse{
		PostId: id,
	}, nil
}

func (s *ExampleService) ListPosts(context.Context, *example.ListPostsRequest) (*example.ListPostsResponse, error) {
	// RPC ошибки подобно HTTP имеют свои коды по которым клиент может ориентироваться.
	// Пакет status и codes необходимы для создания rpc ошибок.
	rpcError := status.Error(codes.Unimplemented, codes.Unimplemented.String())

	return nil, rpcError
}
