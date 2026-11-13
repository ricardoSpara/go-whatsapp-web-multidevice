# Custom image in a link preview

`POST /send/link` can use an image supplied by the caller instead of the
image found in the page's Open Graph metadata. The link still supplies its
title and description.

Use the optional `link_preview_image` field with either bare Base64 or a
`data:image/...;base64,...` URI. The decoded image is limited to 10 MiB and
4000 pixels on either side. JPEG, PNG, GIF, and WebP are supported. Invalid
images return HTTP 400 and the message is not sent. When the field is omitted,
GOWA keeps using the page's Open Graph image.

## Try it with REST

Replace the sample phone, link, caption, and image path. `APP_BASIC_AUTH` must
be configured on the GOWA server.

```sh
base64 -i banner.png | tr -d '\n' > /tmp/banner.base64
jq -n --rawfile image /tmp/banner.base64 '{phone:"5511999999999",link:"https://example.com",caption:"Confira este link",link_preview_image:$image}' \
  | curl -u "$GOWA_USER:$GOWA_PASS" -H 'Content-Type: application/json' \
      --data-binary @- http://localhost:3000/send/link
```

If the server has multiple paired devices, add
`-H 'X-Device-Id: <device-id>'` to the `curl` command.

The same field is available to the MCP `whatsapp_send` tool when `type` is
`link`. Scheduled link sends preserve the image in their stored request.
