FROM node:alpine

LABEL maintainer "Tobias Kohlbua <tobias@kohlbau.de>"

RUN apk update && \
    apk add --update alpine-sdk su-exec && \
    npm install -g @angular/cli && \
    ng set --global packageManager=yarn && \
    apk del alpine-sdk && \
    rm -rf /tmp/* /var/cache/apk/* *.tar.gz ~/.npm && \
    npm cache clear --force && \
    sed -i -e "s/bin\/ash/bin\/sh/" /etc/passwd

RUN mkdir /frontend
COPY . /frontend
COPY entrypoint.sh /usr/local/bin/

WORKDIR /frontend

RUN npm install && \
    ng build -op /dist --prod --env=prod


FROM abiosoft/caddy:latest
WORKDIR /srv
COPY --from=0 /dist /srv
