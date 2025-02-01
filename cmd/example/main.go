package main

import (
	"io"
	"log"
	"net"

	"github.com/bufbuild/protovalidate-go"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/moguchev/grpc/pkg/api/example"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Валидатор protobuf сообщений
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal(err)
	}

	// Создаем gRPC сервер (aka HTTP сервер)
	server := grpc.NewServer(
		/*
		   Когда применять?
		   * UnaryInterceptor: Если нужен единственный перехватчик.
		   * ChainUnaryInterceptor: Когда нужно объединить несколько перехватчиков в один (логирование, мониторинг, аутентификация).
		*/
		// grpc.UnaryInterceptor( // Это обычный перехватчик, который оборачивает обработчик запроса (handler) в дополнительную логику.
		// 	WithProtoValidate(validator), // наша реализация
		// ),
		grpc.ChainUnaryInterceptor( // Когда нужно добавить несколько перехватчиков (например, логирование, валидацию, трассировку), используется цепочка перехватчиков.
			recovery.UnaryServerInterceptor(),
			protovalidate_middleware.UnaryServerInterceptor(validator), // готовая реализация
		),
		// Больше готовых интерсепторов: https://github.com/grpc-ecosystem/go-grpc-middleware/tree/v2.1.0/interceptors
	)

	// Регистрируем gRPC обработчик рефлексии (роутер, мультиплексор)
	// Нужен для того, чтобы можно было узнать API нашего gRPC сервиса и
	// использовать Postman/grpcli/grpcurl с JSON
	reflection.Register(server)

	// Наша имплементация обработчика gRPC
	service := &ExampleService{
		storage:   make(map[uint64]*Post, 1),
		validator: validator,
	}

	// Регистрируем наш gRPC обработчик (роутер, мультиплексор)
	example.RegisterExampleServiceServer(server, service)

	// Создаем TCP листенер
	lis, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("gRPC server listen on :8085")

	// Начинаем обрабатывать входящие запросы
	if err := server.Serve(lis); err != io.EOF {
		log.Fatal(err)
	}
}
