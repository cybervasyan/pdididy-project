package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	inventoryV1 "github.com/cybervasyan/pdididy-project/inventory/internal/api/inventory/v1"
	repoPart "github.com/cybervasyan/pdididy-project/inventory/internal/repository/part"
	servPart "github.com/cybervasyan/pdididy-project/inventory/internal/service/part"
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load .env: %v", err)
		return
	}

	mongoURI := os.Getenv("MONGO_URI")
	databaseName := os.Getenv("MONGO_DATABASE")
	collectionName := os.Getenv("MONGO_COLLECTION")

	mongoClient, err := mongo.Connect(
		options.Client().ApplyURI(mongoURI),
	)
	if err != nil {
		log.Printf("failed to create MongoDB client: %v", err)
		return
	}

	defer func() {
		disconnectCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		if disconnectErr := mongoClient.Disconnect(disconnectCtx); disconnectErr != nil {
			log.Printf("failed to disconnect from MongoDB: %v", disconnectErr)
		}
	}()

	pingCtx, pingCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer pingCancel()

	if err = mongoClient.Ping(pingCtx, readpref.Primary()); err != nil {
		log.Printf("failed to connect to MongoDB: %v", err)
		return
	}

	collection := mongoClient.
		Database(databaseName).
		Collection(collectionName)

	seedCtx, seedCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer seedCancel()

	if err = seedParts(seedCtx, collection); err != nil {
		log.Printf("failed to seed MongoDB: %v", err)
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	s := grpc.NewServer()

	repository := repoPart.NewRepository(collection)
	service := servPart.NewPartService(repository)
	invAPI := inventoryV1.NewAPI(service)

	inventoryv1.RegisterInventoryServiceServer(s, invAPI)

	reflection.Register(s)

	go func() {
		log.Printf("InventoryService gRPC server listening on %d\n", 50051)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
