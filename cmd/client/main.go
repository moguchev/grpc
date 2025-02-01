package main

import (
	"context"
	"fmt"
	"log"

	"github.com/moguchev/grpc/pkg/api/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const ExampleServiceHost = ":8085"

func main() {
	ctx := context.Background()

	// Создаем gRPC соединение
	//
	// Антипатерны: https://github.com/grpc/grpc-go/blob/master/Documentation/anti-patterns.md
	conn, err := grpc.NewClient(
		ExampleServiceHost,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // TLS или insecure
		grpc.WithUnaryInterceptor(DebugInterceptor(true)),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{PermitWithoutStream: true}),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.DefaultConfig}),
		// ...
	)
	if err != nil {
		log.Fatal(err)
	}

	// Создаем gRPC клиента
	client := example.NewExampleServiceClient(conn)

	if err := DoRequest(ctx, client); err != nil {
		log.Fatal(err)
	}
}

func DoRequest(ctx context.Context, client example.ExampleServiceClient) error {
	// Запрос
	req := &example.CreatePostRequest{
		Title:    "example",
		AuthorId: "1",
		Content:  "hello",
	}

	// Заголовки, которые передадим вместе с запросом
	md := metadata.Pairs("client-header-key", "val")
	octx := metadata.NewOutgoingContext(ctx, md)

	var headers, trailers = metadata.MD{}, metadata.MD{}
	// Выполняем gRPC запрос.
	// Для того чтобы получить заголовки из ответа,
	// требуется использовать опции grpc.Header/grpc.Trailer
	resp, err := client.CreatePost(octx, req,
		grpc.Header(&headers),
		grpc.Trailer(&trailers),
	)
	if err != nil {
		// Для того чтобы обработать код ошибки
		switch status.Code(err) {
		case codes.InvalidArgument:
			log.Println("некоректный запрос")
		default:
			log.Println("неожиданная ошибка")
			if st, ok := status.FromError(err); ok {
				log.Println(
					"code:", st.Code(),
					"message:", st.Message(),
					"details:", st.Details(),
				)
			} else {
				log.Println("not grpc") // такого не должно быть в реальности
			}
		}
		return err
	}

	fmt.Println(
		"headers:", headers,
		"trailers:", trailers,
	)

	log.Println(resp.GetPostId())

	// Для корректного Marshal/Unmarshal protobuf сообщений в JSON
	// следует использовать пакет protojson
	bytes, err := protojson.Marshal(resp)
	if err != nil {
		return err
	}

	log.Println(string(bytes))
	return nil
}
