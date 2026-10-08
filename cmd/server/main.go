package main

import (
	"log"
	"net"

	prayerv1 "rushd-backend/gen/prayer/v1"
	"rushd-backend/internal/prayer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcAddress = "0.0.0.0:50051"

func main() {
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
