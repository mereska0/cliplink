package grpcclient

import (
	"context"

	"github.com/mereska0/cliplink/api/gen/linkpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LinkClient struct {
	conn   *grpc.ClientConn
	client linkpb.LinkServiceClient
}

func NewLinkClient(addr string) (*LinkClient, error) {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &LinkClient{
		conn:   conn,
		client: linkpb.NewLinkServiceClient(conn),
	}, nil
}

func (c *LinkClient) Close() error {
	return c.conn.Close()
}

func (c *LinkClient) CreateLink(ctx context.Context, originalURL string, customAlias string) (*linkpb.Link, error) {
	resp, err := c.client.CreateLink(ctx, &linkpb.CreateLinkRequest{
		OriginalUrl: originalURL,
		CustomAlias: customAlias,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetLink(), nil
}

func (c *LinkClient) ListLinks(ctx context.Context) ([]*linkpb.Link, error) {
	resp, err := c.client.ListLinks(ctx, &linkpb.ListLinksRequest{})
	if err != nil {
		return nil, err
	}

	return resp.GetLinks(), nil
}

func (c *LinkClient) DeleteLink(ctx context.Context, shortCode string) error {
	_, err := c.client.DeleteLink(ctx, &linkpb.DeleteLinkRequest{
		ShortCode: shortCode,
	})

	return err
}

func (c *LinkClient) GetLink(ctx context.Context, shortCode string) (*linkpb.Link, error) {
	resp, err := c.client.GetLink(ctx, &linkpb.GetLinkRequest{
		ShortCode: shortCode,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetLink(), nil
}
