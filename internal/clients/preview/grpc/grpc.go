package grpc

import (
	"context"
	"fmt"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	previewv1 "github.com/sha1sof/grpc_Echelon/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	api previewv1.ThumbnailServiceClient
}

// New конструктор клиента.
func New(ctx context.Context, address string, timeout time.Duration, retriesCount int) (*Client, error) {
	const op = "Client.grpc.New"

	retryOptions := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	cc, err := grpc.DialContext(ctx, address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOptions...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		api: previewv1.NewThumbnailServiceClient(cc),
	}, nil
}

// GetThumbnail возвращает данные в байтах
func (c *Client) GetThumbnail(ctx context.Context, videoURL string) ([]byte, error) {
	const op = "Client.GetThumbnail"

	resp, err := c.api.GetThumbnail(ctx, &previewv1.ThumbnailRequest{
		VideoUrl: videoURL,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp.ThumbnailData, nil
}
