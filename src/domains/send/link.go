package send

type LinkRequest struct {
	BaseRequest
	Caption          string `json:"caption"`
	Link             string `json:"link"`
	LinkPreviewImage string `json:"link_preview_image,omitempty"`
}
