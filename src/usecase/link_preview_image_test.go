package usecase

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"testing"

	"github.com/aldinokemal/go-whatsapp-web-multidevice/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyCustomLinkPreviewImagePreservesPageMetadata(t *testing.T) {
	var encoded bytes.Buffer
	require.NoError(t, png.Encode(&encoded, image.NewRGBA(image.Rect(0, 0, 1, 1))))
	metadata := utils.Metadata{
		Title:       "Page title",
		Description: "Page description",
		Image:       "https://example.com/og.png",
		ImageThumb:  []byte("old image thumbnail"),
		JPEGThumb:   []byte("old jpeg thumbnail"),
	}

	err := applyCustomLinkPreviewImage(&metadata, base64.StdEncoding.EncodeToString(encoded.Bytes()))
	require.NoError(t, err)
	assert.Equal(t, "Page title", metadata.Title)
	assert.Equal(t, "Page description", metadata.Description)
	assert.Equal(t, "https://example.com/og.png", metadata.Image)
	assert.NotEqual(t, []byte("old image thumbnail"), metadata.ImageThumb)
	assert.NotEqual(t, []byte("old jpeg thumbnail"), metadata.JPEGThumb)
	require.NotNil(t, metadata.Width)
	require.NotNil(t, metadata.Height)
	assert.Equal(t, uint32(1), *metadata.Width)
	assert.Equal(t, uint32(1), *metadata.Height)
}

func TestApplyCustomLinkPreviewImageRejectsMalformedSource(t *testing.T) {
	err := applyCustomLinkPreviewImage(&utils.Metadata{}, "not-base64")
	assert.Error(t, err)
}
