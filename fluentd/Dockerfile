FROM fluent/fluentd:v1.16

USER root

RUN apk add --no-cache --update \
    build-base \
    ruby-dev \
    ruby-bundler \
    ruby-rake \
    ruby-irb \
    librdkafka-dev

#COPY ./fluentd.conf  /fluentd/etc/fluent.conf

RUN gem install fluent-plugin-kafka --no-document
