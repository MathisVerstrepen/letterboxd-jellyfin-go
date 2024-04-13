FROM golang:1.22.1-bookworm

RUN apt-get update && apt-get -y install cron vim

WORKDIR /app

COPY crontab /etc/cron.d/crontab
RUN chmod 0644 /etc/cron.d/crontab && /usr/bin/crontab /etc/cron.d/crontab

COPY go.mod go.sum ./

RUN go mod download

COPY . .
COPY ./.env /app/.env

RUN go build -o main .

CMD ["cron", "-f"]