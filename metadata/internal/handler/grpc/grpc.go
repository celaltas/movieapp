package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller/metadata"
	"movieexample.com/metadata/pkg/model"
)

type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetMetadata(ctx context.Context, in *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	if in == nil || in.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	m, err := h.ctrl.Get(ctx, in.Id)
	if err != nil && errors.Is(err, metadata.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetMetadataResponse{Metadata: model.MetadataToProto(m)}, nil
}


// PutMetadata puts movie metadata to repository.
func (h *Handler) PutMetadata(ctx context.Context, req *gen.PutMetadataRequest) (*gen.PutMetadataResponse, error) {
	if req == nil || req.Metadata == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or metadata")
	}
	if err := h.ctrl.Put(ctx, model.MetadataFromProto(req.Metadata)); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.PutMetadataResponse{}, nil
}