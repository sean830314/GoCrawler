FROM fluent/fluentd:v1.3

# add mongo plugin
RUN apk add --no-cache bash make gcc libc-dev ruby-dev \
    && gem install fluent-plugin-mongo \
    && gem install  fluent-plugin-redis-store \
    && apk del make gcc libc-dev ruby-dev \
    && rm -rf /var/cache/apk/* \
    && rm -rf /tmp/* /var/tmp/* /usr/lib/ruby/gems/*/cache/*.gem
