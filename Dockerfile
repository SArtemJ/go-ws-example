FROM golang:alpine
ENV PROJECT_PATH=github.com/SArtemJ/wstest
ARG _path
RUN apk add --no-cache --update \
    git

RUN mkdir -p ${GOPATH}/src/${PROJECT_PATH}
WORKDIR ${GOPATH}/src/${PROJECT_PATH}
COPY . .
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN go build
ENTRYPOINT [ "./wstest" ]
EXPOSE 8099