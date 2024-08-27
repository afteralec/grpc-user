FROM golang:1.22-alpine3.20 AS build

ENV CGO_ENABLED=1
ENV GROUP_ID=65535
ENV GROUP_NAME=noroot
ENV USER_ID=65535
ENV USER_NAME=noroot

RUN apk add --no-cache \
  gcc=13.2.1_git20240309-r0 \
  musl-dev=1.2.5-r0

WORKDIR /workdir

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags='-s -w -extldflags "-static"' -o /out/user .

RUN addgroup -g $GROUP_ID -S $GROUP_NAME && adduser -D -H -S -G $GROUP_NAME -u $USER_ID $USER_NAME

FROM scratch AS run

COPY --from=build /etc/group /etc/group
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /out/user /usr/local/bin/user

USER noroot:noroot

EXPOSE 8009

ENTRYPOINT ["/usr/local/bin/user"]
