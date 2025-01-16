FROM node:18-slim

RUN apt-get update \
  && apt-get install --no-install-recommends --yes \
    ca-certificates \
    curl \
    jq \
    procps \
    tini \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Get Home Assistan add-on Bashio
ARG BASHIO_VERSION

RUN mkdir -p /usr/lib/bashio \
  && curl -L "https://github.com/hassio-addons/bashio/archive/v${BASHIO_VERSION}.tar.gz" | tar -xz -C /usr/lib/bashio --strip-components=1 \
  && ln -s /usr/lib/bashio/lib/bashio /usr/bin/bashio

COPY --chmod=555 docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh

ARG user=joplin
RUN useradd --create-home --shell /bin/bash $user

USER $user

ENV NODE_ENV=production
ENV RUNNING_IN_DOCKER=1

EXPOSE $APP_PORT

WORKDIR /home/${user}/packages/server
ENTRYPOINT ["tini", "--"]
CMD [ "/usr/local/bin/docker-entrypoint.sh" ]

ARG APP_REVISION
ARG APP_VERSION
ARG BUILD_ARCH
ARG BUILD_DATE
ARG BUILD_VERSION
LABEL \
    io.hass.name="Joplin Server" \
    io.hass.description="Home Assistant Add-on for Joplin Server" \
    io.hass.arch="${BUILD_ARCH}" \
    io.hass.type="addon" \
    io.hass.version=${BUILD_VERSION} \
    org.opencontainers.image.title="Joplin Server" \
    org.opencontainers.image.description="Docker image for Joplin Server" \
    org.opencontainers.image.url="https://joplinapp.org/" \
    org.opencontainers.image.source="https://github.com/laurent22/joplin.git" \
    org.opencontainers.image.created=${BUILD_DATE} \
    org.opencontainers.image.revision=${APP_REVISION} \
    org.opencontainers.image.version=${APP_VERSION}

ENV HEALTH_PORT="" \
    HEALTH_URL=""

HEALTHCHECK \
    --interval=5s \
    --retries=5 \
    --start-period=30s \
    --timeout=25s \
    CMD ps aux | grep "joplin" | grep -v "grep" || exit 1
