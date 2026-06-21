package grpcserver

import (
	"context"
	"errors"

	"github.com/mereska0/cliplink/api/gen/linkpb"
	"github.com/mereska0/cliplink/internal/domain"
	"github.com/mereska0/cliplink/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LinkServer struct {
	linkpb.UnimplementedLinkServiceServer

	linkService *service.LinkService
}

func NewLinkServer(linkService *service.LinkService) *LinkServer {
	return &LinkServer{
		linkService: linkService,
	}
}

func (s *LinkServer) CreateLink(ctx context.Context, req *linkpb.CreateLinkRequest) (*linkpb.CreateLinkResponse, error) {
	link, err := s.linkService.CreateLink(ctx, service.CreateLinkInput{
		OriginalURL: req.GetOriginalUrl(),
		CustomAlias: req.GetCustomAlias(),
	})
	if err != nil {
		return nil, mapError(err)
	}

	return &linkpb.CreateLinkResponse{
		Link: toProtoLink(link),
	}, nil
}

func (s *LinkServer) ListLinks(ctx context.Context, req *linkpb.ListLinksRequest) (*linkpb.ListLinksResponse, error) {
	links, err := s.linkService.ListLinks(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	result := make([]*linkpb.Link, 0, len(links))

	for _, link := range links {
		result = append(result, toProtoLink(&link))
	}

	return &linkpb.ListLinksResponse{
		Links: result,
	}, nil
}

func (s *LinkServer) DeleteLink(ctx context.Context, req *linkpb.DeleteLinkRequest) (*linkpb.DeleteLinkResponse, error) {
	err := s.linkService.DeleteLink(ctx, req.GetShortCode())
	if err != nil {
		return nil, mapError(err)
	}

	return &linkpb.DeleteLinkResponse{}, nil
}

func (s *LinkServer) GetLink(ctx context.Context, req *linkpb.GetLinkRequest) (*linkpb.GetLinkResponse, error) {
	link, err := s.linkService.GetLink(ctx, req.GetShortCode())
	if err != nil {
		return nil, mapError(err)
	}

	return &linkpb.GetLinkResponse{
		Link: toProtoLink(link),
	}, nil
}

func toProtoLink(link *domain.Link) *linkpb.Link {
	return &linkpb.Link{
		Id:          link.ID,
		ShortCode:   link.ShortCode,
		OriginalUrl: link.OriginalURL,
		Clicks:      link.Clicks,
		CreatedAt:   link.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidURL):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrLinkNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrAliasTaken):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrDeletedLink):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
