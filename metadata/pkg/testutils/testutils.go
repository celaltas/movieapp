package testutils

import (
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller/metadata"
	"movieexample.com/metadata/internal/handler/grpc"
	"movieexample.com/metadata/internal/repository/memory"
)

func NewTestMetadataGRPCServer() gen.MetadataServiceServer {
	repository := memory.New()
	controller := metadata.New(repository)
	handler := grpc.New(controller)
	return handler
}
