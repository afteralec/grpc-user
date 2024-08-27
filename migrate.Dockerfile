FROM rust:1.80

ENV GROUP_ID=65535
ENV GROUP_NAME=noroot
ENV USER_ID=65535
ENV USER_NAME=noroot

WORKDIR /opt/mirror

RUN cargo install sqlx-cli

VOLUME /var/db

COPY ./migrations ./migrations

COPY ./scripts/migrate.sh /usr/local/bin/migrate.sh

RUN chmod +x /usr/local/bin/migrate.sh

RUN groupadd -g $GROUP_ID $GROUP_NAME && useradd -l -g $GROUP_ID -u $USER_ID $USER_NAME

CMD ["/usr/local/bin/migrate.sh"]
