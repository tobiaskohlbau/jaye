FROM golang:alpine
LABEL maintainer="Tobias Kohlbau <tobias@kohlbau.de>"

RUN apk add --update ffmpeg

VOLUME [ "/go/src/kohlbau.de/x/jaye", "/videos" ]

COPY . /go/src/kohlbau.de/x/jaye
WORKDIR /go/src/kohlbau.de/x/jaye

RUN go-wrapper install

CMD ["go-wrapper", "run", "-config", "./config/config.prod.json"]

EXPOSE 8080
