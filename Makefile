# Builder image
builder-image:
	docker build -t kai5263499/image-ddns-builder -f Dockerfile .

exec-interactive:
	docker run -it \
	-e CLOUDFLARE_API_KEY=${CLOUDFLARE_API_KEY} \
	-e CLOUDFLARE_API_EMAIL=${CLOUDFLARE_API_EMAIL} \
	-e NAME=${NAME} \
	-e ZONE=${ZONE} \
	-e IMAGE_URL=${IMAGE_URL} \
	-e LOG_LEVEL=${LOG_LEVEL} \
	-v /home/wes/code/deproot/src/github.com/kai5263499:/go/src/github.com/kai5263499 \
	--tmpfs /tmp:exec \
	-w /go/src/github.com/kai5263499/image-ddns/cmd/image-ddns \
	kai5263499/image-ddns-builder

image-ddns-image: builder-image
	docker build -t kai5263499/image-ddns -f cmd/image-ddns/Dockerfile .

image-ddns: image-ddns-image
	docker run -it --rm \
	-e CLOUDFLARE_API_KEY=${CLOUDFLARE_API_KEY} \
	-e CLOUDFLARE_API_EMAIL=${CLOUDFLARE_API_EMAIL} \
	-e NAME=${NAME} \
	-e ZONE=${ZONE} \
	-e IMAGE_URL=${IMAGE_URL} \
	-e LOG_LEVEL=${LOG_LEVEL} \
	--tmpfs /tmp \
	kai5263499/image-ddns
