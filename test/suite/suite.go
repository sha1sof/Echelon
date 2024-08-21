package suite

import (
	"github.com/sha1sof/Echelon-/internal/config"
	previewv1 "github.com/sha1sof/grpc_Echelon/gen"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

type Suite struct {
	*testing.T
	Cfg     *config.Config
	TClient previewv1.ThumbnailServiceClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath("../config/prod.yaml")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(context.Background(),
		"localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	return ctx, &Suite{
		T:       t,
		Cfg:     cfg,
		TClient: previewv1.NewThumbnailServiceClient(cc),
	}
}
