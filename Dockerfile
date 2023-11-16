FROM golang:1.21 as builder

## Create appuser.
# ENV USER=appuser
# ENV UID=10001

## See https://stackoverflow.com/a/55757473/12429735RUN
# RUN adduser \
#   --disabled-password \
#   --gecos "" \
#   --home "/nonexistent" \
#   --shell "/sbin/nologin" \
#   --no-create-home \
#   --uid "${UID}" \
#   "${USER}"

WORKDIR /tmp/go

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -o server
RUN go build -o ./seeder-runner ./seeder
RUN go build -o ./migrator-runner ./migrator

FROM golang:1.21

WORKDIR /

ARG IMAGE_VERSION_TAG
ENV IMAGE_VERSION_TAG=${IMAGE_VERSION_TAG}

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /tmp/go/server /server
COPY --from=builder /tmp/go/seeder-runner /seeder
COPY --from=builder /tmp/go/migrator-runner /migrator
COPY --from=builder /tmp/go/.env* /
COPY ./wait-for-it.sh /wait-for-it.sh

### Use an unprivileged user.
## USER appuser:appuser

EXPOSE 4040

CMD ["/server"]
