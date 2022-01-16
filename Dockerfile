FROM amd64/debian:stable-slim
WORKDIR /
RUN export DEBIAN_FRONTEND=noninteractive && apt-get update && apt-get install -yq \
  tzdata \
  curl \
  ca-certificates \
  && rm -rf /var/lib/apt/lists/*


ENV TBOTKEY telegramkey
ENV DSN postgres_db_dsn
ENV APIKEY apy_key
ENV USERPASS userpass

COPY ./bot /bot
CMD "/bot"
