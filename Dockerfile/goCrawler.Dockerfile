FROM golang:1.14 AS build-env
LABEL maintainer="kroos.chen" \
      build-date={BUILD-DATE} \
      description="A service for golang crawler of ptt." \
      distribution-scope="private" \
      name={IMAGE} \
      release="0" \
      summary="A service for golang crawler of ptt." \
      vcs-ref={VCS-REF} \
      vcs-type="git" \
      vendor="kroos.chen" \
      version={VERSION}
ADD . /src
RUN cd /src && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goCrawler main.go
RUN  mkdir /etc/GoCrawler && cp /src/config/config.yaml /etc/GoCrawler/config.yaml
CMD ["/src/goCrawler"]
