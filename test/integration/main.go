package main

import (
	"context"
	"log"
	"net"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"movieexample.com/gen"
	"movieexample.com/metadata/pkg/testutils"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/memory"
)

const (
	metadataServiceName = "metadata"
	ratingServiceName   = "rating"
	movieServiceName    = "movie"
	metadataServiceAddr = "localhost:8081"
	ratingServiceAddr   = "localhost:8082"
	movieServiceAddr    = "localhost:8083"
)

func main() {
	log.Println("Starting the integration test")
	ctx := context.Background()
	registry := memory.NewRegistry()
	log.Println("Setting up service handlers and clients")
	metadataSrv := startMetadataService(ctx, registry)
	defer metadataSrv.GracefulStop()


	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	metadataConn, err := grpc.Dial(metadataServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer metadataConn.Close()
	metadataClient := gen.NewMetadataServiceClient(metadataConn)

	log.Println("Saving test metadata via metadata service")

	m := &gen.Metadata{
		Id:          "the movie",
		Title:       "The Movie",
		Description: "The Movie, the one and only",
		Director:    "Mr D.",
	}

	if _, err := metadataClient.PutMetadata(ctx, &gen.PutMetadataRequest{Metadata: m}); err != nil {
		log.Fatalf("put metadata: %v", err)
	}

	log.Println("Retrieving test metadata via metadata service")

	getMetadataResp, err := metadataClient.GetMetadata(ctx, &gen.GetMetadataRequest{Id: m.Id})
	if err != nil {
		log.Fatalf("get metadata: %v", err)
	}

	if diff := cmp.Diff(getMetadataResp.Metadata, m, cmpopts.
		IgnoreUnexported(gen.Metadata{})); diff != "" {
		log.Fatalf("get metadata after put mismatch: %v", diff)
		}

}

func startMetadataService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("Starting metadata service on " + metadataServiceAddr)
	handler := testutils.NewTestMetadataGRPCServer()
	listener, err := net.Listen("tcp", metadataServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, handler)
	go func() {
		if err := srv.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v\n", err)
		}
	}()

	id := discovery.GenerateInstanceID(metadataServiceName)
	if err := registry.Register(ctx, id, metadataServiceName, metadataServiceAddr); err != nil {
		panic(err)
	}
	return srv

}
