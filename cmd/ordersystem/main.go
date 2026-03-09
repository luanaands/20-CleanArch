package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"time"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/luanaands/clean-arch/configs"
	"github.com/luanaands/clean-arch/internal/event/handler"
	"github.com/luanaands/clean-arch/internal/infra/graph"
	"github.com/luanaands/clean-arch/internal/infra/grpc/pb"
	"github.com/luanaands/clean-arch/internal/infra/grpc/service"
	"github.com/luanaands/clean-arch/internal/infra/web/webserver"
	"github.com/luanaands/clean-arch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	rabbitMQChannel, err := getRabbitMQChannel()
	if err != nil {
		panic(err)
	}

	db, err := connectMySQL(configs)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		panic(err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://sql/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		panic(err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(err)
	}

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	getAllOrdersUseCase := NewGetAllOrdersUseCase(db)
	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/orders", webOrderHandler.GetAll)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go func() {
		if err := webserver.Start(); err != nil {
			panic(err)
		}
	}()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *getAllOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase:  *createOrderUseCase,
		GetAllOrdersUseCase: *getAllOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() (*amqp.Channel, error) {
	url := "amqp://guest:guest@rabbitmq:5672/"
	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				return ch, nil
			}
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return nil, fmt.Errorf("rabbitmq não ficou pronto: %w", err)
}

func connectMySQL(cfg *configs.Conf) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		return nil, err
	}

	for i := 0; i < 10; i++ {
		if err = db.Ping(); err == nil {
			return db, nil
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return nil, fmt.Errorf("mysql não ficou pronto: %w", err)
}
