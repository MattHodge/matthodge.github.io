FROM golang:1.12.5 as builder
ENV GOPROXY https://proxy.golang.org

# Download Vale
ENV VALE_VERSION 1.7.1
ENV VALE_SHA 239765ab91992fdc9e7482f73de087fd038222bec05b44e7c1f4fff637110b3e
RUN mkdir /output
RUN mkdir /temp-download && \
    cd /temp-download && \
    curl -OL https://github.com/errata-ai/vale/releases/download/v${VALE_VERSION}/vale_${VALE_VERSION}_Linux_64-bit.tar.gz && \
    echo "${VALE_SHA}  vale_${VALE_VERSION}_Linux_64-bit.tar.gz" | sha256sum -c && \
    tar -zxvf vale_${VALE_VERSION}_Linux_64-bit.tar.gz && \
    chmod +x vale && \
    mv vale /output/vale

ENV GO111MODULE on
WORKDIR /go/src/github.com/MattHodge/.github/actions/vale
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
# static compile artifact-generator
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /output/vale-entrypoint && \
    chmod +x /output/vale-entrypoint



# Output Container
FROM alpine:3.9

# Default inputs defined here due to defaults not working in GitHub Actions
# BUG: https://github.community/t5/GitHub-Actions/Metadata-Syntax-for-Inputs-Defaults-Not-Working-for-Actions/m-p/31067#M660
ENV INPUT_FILEGLOB="*.{md,rst}"
ENV INPUT_LINTDIRECTORY="."
ENV INPUT_LINTUNCHANGEDFILES="true"
ENV INPUT_CONFIGFILEPATH=".vale.ini"
RUN apk update \
  && apk upgrade \
  && apk add --no-cache \
  ca-certificates \
  libc6-compat \
  && update-ca-certificates 2>/dev/null || true
COPY --from=builder /output/vale-entrypoint /usr/local/bin/entrypoint
COPY --from=builder /output/vale /usr/local/bin/vale
COPY .vale.ini /etc/vale/.vale.ini
#COPY entrypoint.sh /entrypoint.sh
#RUN chmod +x /entrypoint.sh
#ENTRYPOINT ["/entrypoint.sh"]
ENTRYPOINT ["/usr/local/bin/entrypoint"]
