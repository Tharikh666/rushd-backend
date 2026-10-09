package main

import (
	"log"
	"net"
	"os"

	prayerv1 "rushd-backend/gen/prayer/v1"
	"rushd-backend/internal/prayer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const defaultPort = "50051"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	grpcAddress := ":" + port

	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		log.Fatalf(
			"failed to listen on %s: %v",
			grpcAddress,
			err,
		)
	}

	grpcServer := grpc.NewServer()

	// ------------------------------------------------------------
	// Prayer
	// ------------------------------------------------------------

	prayerCalculator := prayer.NewCalculator()

	prayerService := prayer.NewService(
		prayerCalculator,
	)

	prayerv1.RegisterPrayerServiceServer(
		grpcServer,
		prayerService,
	)

	reflection.Register(grpcServer)

	// ------------------------------------------------------------
	// Start server
	// ------------------------------------------------------------

	log.Printf(
		"RUSHD backend listening on %s",
		grpcAddress,
	)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf(
			"gRPC server stopped: %v",
			err,
		)
	}
}
