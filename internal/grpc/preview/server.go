package preview

import (
	"context"
	preview "github.com/sha1sof/grpc_Echelon/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type Thumbnail interface {
	GetThumbnail(
		ctx context.Context,
		url string,
	) (preVideo []byte, err error)
}

type serverApi struct {
	preview.UnimplementedThumbnailServiceServer
	thumbnail Thumbnail
}

// Register сервера
func Register(gRPC *grpc.Server, thumbnail Thumbnail) {
	preview.RegisterThumbnailServiceServer(gRPC, &serverApi{thumbnail: thumbnail})
}

// GetThumbnail
func (s *serverApi) GetThumbnail(ctx context.Context, req *preview.ThumbnailRequest) (*preview.ThumbnailResponse, error) {
	err := validUrl(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	preVideo, err := s.thumbnail.GetThumbnail(ctx, req.GetVideoUrl())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "InvalidArgument error")
	}

	return &preview.ThumbnailResponse{
		ThumbnailData: preVideo,
	}, nil
}

// validUrl проверка на валидность.
func validUrl(req *preview.ThumbnailRequest) error {
	videoUrl := req.GetVideoUrl()
	if videoUrl == "" {
		return status.Error(codes.InvalidArgument, "empty video url")
	}

	ytReg := regexp.MustCompile(`^(http(s)?:\/\/)?((w){3}.)?youtu(be|.be)?(\.com)?\/.+`)
	if !ytReg.MatchString(videoUrl) {
		return status.Error(codes.InvalidArgument, "invalid YouTube video url")
	}

	return nil
}
