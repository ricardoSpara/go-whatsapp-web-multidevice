package utils_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"testing"

	"github.com/aldinokemal/go-whatsapp-web-multidevice/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func encodePreviewPNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Set(0, 0, color.RGBA{R: 120, G: 60, B: 200, A: 255})
	var encoded bytes.Buffer
	require.NoError(t, png.Encode(&encoded, img))
	return encoded.Bytes()
}

func TestBuildLinkPreviewThumbnailsFromBase64(t *testing.T) {
	imageData := encodePreviewPNG(t, 64, 32)
	source := base64.StdEncoding.EncodeToString(imageData)

	for _, input := range []struct {
		name   string
		source string
	}{
		{name: "bare base64", source: source},
		{name: "data URI", source: fmt.Sprintf("data:image/png;base64,%s", source)},
	} {
		t.Run(input.name, func(t *testing.T) {
			imageThumb, jpegThumb, width, height, err := utils.BuildLinkPreviewThumbnailsFromBase64(input.source)
			require.NoError(t, err)
			assert.NotEmpty(t, imageThumb)
			assert.NotEmpty(t, jpegThumb)
			assert.Equal(t, uint32(64), width)
			assert.Equal(t, uint32(32), height)
		})
	}
}

func TestBuildLinkPreviewThumbnailsFromBase64RejectsInvalidImages(t *testing.T) {
	for _, input := range []struct {
		name   string
		source string
	}{
		{name: "malformed base64", source: "%%%"},
		{name: "non-image data", source: base64.StdEncoding.EncodeToString([]byte("plain text"))},
		{name: "non-image data URI", source: "data:text/plain;base64,aGVsbG8="},
	} {
		t.Run(input.name, func(t *testing.T) {
			_, _, _, _, err := utils.BuildLinkPreviewThumbnailsFromBase64(input.source)
			assert.Error(t, err)
		})
	}
}

func TestBuildLinkPreviewThumbnailsFromBase64EnforcesLimits(t *testing.T) {
	t.Run("decoded bytes", func(t *testing.T) {
		source := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0}, utils.LinkPreviewImageMaxBytes+1))
		_, _, _, _, err := utils.BuildLinkPreviewThumbnailsFromBase64(source)
		assert.ErrorContains(t, err, "exceeds")
	})

	t.Run("image dimensions", func(t *testing.T) {
		source := base64.StdEncoding.EncodeToString(encodePreviewPNG(t, utils.LinkPreviewImageMaxDimension+1, 1))
		_, _, _, _, err := utils.BuildLinkPreviewThumbnailsFromBase64(source)
		assert.ErrorContains(t, err, "exceeds")
	})
}
