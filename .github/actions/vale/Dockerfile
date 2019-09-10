FROM alpine:3.9
ENV VALE_SHA=239765ab91992fdc9e7482f73de087fd038222bec05b44e7c1f4fff637110b3e
# Default inputs defined here due to defaults not working in GitHub Actions
# BUG: https://github.community/t5/GitHub-Actions/Metadata-Syntax-for-Inputs-Defaults-Not-Working-for-Actions/m-p/31067#M660
ENV INPUT_FILEGLOB="*.{md,rst}"
ENV INPUT_LINTDIRECTORY="."
ENV INPUT_LINTUNCHANGEDFILES="true"
RUN apk update \
  && apk upgrade \
  && apk add --no-cache \
  ca-certificates \
  curl \
  libc6-compat \
  && update-ca-certificates 2>/dev/null || true
RUN mkdir /temp-download && \
    cd /temp-download && \
    curl -OL https://github.com/errata-ai/vale/releases/download/v1.7.1/vale_1.7.1_Linux_64-bit.tar.gz && \
    echo "${VALE_SHA}  vale_1.7.1_Linux_64-bit.tar.gz" | sha256sum -c && \
    tar -zxvf vale_1.7.1_Linux_64-bit.tar.gz && \
    chmod +x vale && \
    mv vale /usr/local/bin/ && \
    rm -rf /temp-download
COPY .vale.ini /root
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
