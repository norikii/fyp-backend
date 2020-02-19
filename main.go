package main

import (
	"context"
	"fmt"
	"github.com/tatrasoft/fyp-backend/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/tatrasoft/fyp-backend/config"
	"github.com/tatrasoft/fyp-backend/database"
	"github.com/tatrasoft/fyp-backend/handlers"
)

func main() {
	ctx := context.Background()

	fmt.Println("Starting server on port :50051...")
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption

	s := grpc.NewServer(opts...)
	// services
	itemService := &handlers.ItemsServerService{}
	proto.RegisterItemServiceServer(s, itemService)

	// obtaining config from config file
	dbConfig := config.DBConfig{}
	dbConf, err := dbConfig.GetConfig("config.json")
	if err != nil {
		panic(err)
	}

	// getting the database client
	dbClient, err := database.NewClient(&dbConfig)
	if err != nil {
		panic(err)
	}

	// connecting to the mongo client
	fmt.Println("Connecting to database...")
	err = dbClient.Connect(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected")
	//defer dbClient.CloseConnection(ctx)

	// test the connection before connecting
	fmt.Println("Checking connection...")
	err = dbClient.Ping(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Check successful")

	dbHelper := dbClient.Database(dbConf.DatabaseName)
	handlers.SetColHelper(dbHelper)

	go func() {
		if err := s.Serve(listener); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Server successfully started on port: 50051")

	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user doesn't press CTRL+C a message is not passed and our main routine keeps running
	<-c

	fmt.Println("\nstopping the server...")
	s.Stop()
	listener.Close()

	fmt.Println("Closing MongoDB connection")
	err = dbClient.CloseConnection(ctx)
	if err != nil {
		log.Fatalf("unable to close db connection: %v", err)
	}
	fmt.Println("Done.")
}
