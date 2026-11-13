package validations

import (
	"context"
	"strings"
	"testing"

	domainSend "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/send"
	pkgError "github.com/aldinokemal/go-whatsapp-web-multidevice/pkg/error"
	"github.com/aldinokemal/go-whatsapp-web-multidevice/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestValidateSendLinkPreviewImageSize(t *testing.T) {
	baseRequest := domainSend.LinkRequest{
		BaseRequest: domainSend.BaseRequest{Phone: "1728937129312@s.whatsapp.net"},
		Caption:     "description",
		Link:        "https://example.com",
	}

	t.Run("accepts a small optional image", func(t *testing.T) {
		request := baseRequest
		request.LinkPreviewImage = "aGVsbG8="
		assert.NoError(t, ValidateSendLink(context.Background(), request))
	})

	t.Run("rejects an oversized encoded image before fetching the link", func(t *testing.T) {
		request := baseRequest
		request.LinkPreviewImage = strings.Repeat("A", utils.LinkPreviewImageMaxEncodedBytes+1)
		assert.Equal(t,
			pkgError.ValidationError("link_preview_image exceeds the maximum encoded size"),
			ValidateSendLink(context.Background(), request),
		)
	})
}
