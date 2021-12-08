# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build \
    		-mod=vendor \
    		-o /bin/urlshortener \
    		cmd/urlshortener/main.go

EXPOSE 8080

CMD [ "/bin/urlshortener" ]
