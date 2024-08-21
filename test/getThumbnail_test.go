package test

import (
	"github.com/sha1sof/Echelon-/test/suite"
	preview "github.com/sha1sof/grpc_Echelon/gen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestGetThumbnail_HappyPath загрузка в кеш и получение из кеша.
func TestGetThumbnail_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	// Нужно менять все время но новую ссылку, которой нет в БД.
	url := "https://www.youtube.com/watch?v=e_pY0btswmk"

	res, err := st.TClient.GetThumbnail(ctx, &preview.ThumbnailRequest{
		VideoUrl: url,
	})
	require.NoError(t, err)
	assert.Empty(t, res.ThumbnailData)

	res, err = st.TClient.GetThumbnail(ctx, &preview.ThumbnailRequest{
		VideoUrl: url,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, res.ThumbnailData)
}

// TestGetThumbnail_Bad проверка на валидность ссылки.
func TestGetThumbnail_Bad(t *testing.T) {
	ctx, st := suite.New(t)

	url := "https://www.youtsadube.com/watcsdwash?v=r2sLp9ZOq0I"

	_, err := st.TClient.GetThumbnail(ctx, &preview.ThumbnailRequest{
		VideoUrl: url,
	})
	require.Error(t, err)
}

// TestGetThumbnail_Get получения видео уже которое есть к кеше.
// Делается после TestGetThumbnail_HappyPath.
func TestGetThumbnail_Get(t *testing.T) {
	ctx, st := suite.New(t)

	url := "https://www.youtube.com/watch?v=e_pY0btswmk"

	res, err := st.TClient.GetThumbnail(ctx, &preview.ThumbnailRequest{
		VideoUrl: url,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, res.ThumbnailData)
}
