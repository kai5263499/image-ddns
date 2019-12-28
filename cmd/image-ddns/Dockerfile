FROM kai5263499/image-ddns-builder as builder

COPY / /go/src/github.com/kai5263499/image-ddns

RUN cd /go/src/github.com/kai5263499/image-ddns/cmd/image-ddns && \
    go build && \
    ldd image-ddns | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname deps%); cp % deps%;'

FROM scratch

LABEL MAINTAINER="Wes Widner <kai5263499@gmail.com>"

ENV CLOUDFLARE_API_KEY=""
ENV CLOUDFLARE_API_EMAIL=""
ENV NAME=""
ENV ZONE=""
ENV IMAGE_URL=""
ENV LOG_LEVEL="debug"
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/4.00/tessdata

COPY --from=builder /usr/share/tesseract-ocr/4.00/tessdata /usr/share/tesseract-ocr/4.00/tessdata
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/kai5263499/image-ddns/cmd/image-ddns/deps /
COPY --from=builder /go/src/github.com/kai5263499/image-ddns/cmd/image-ddns/image-ddns /image-ddns

ENTRYPOINT [ "/image-ddns" ]