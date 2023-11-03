FROM golang:1.21-alpine3.18
WORKDIR /app

# Download tailwindcss
RUN apk add --no-cache curl
RUN curl -o /bin/tailwindcss -sL https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.5/tailwindcss-linux-x64
RUN chmod +x /bin/tailwindcss

COPY tailwind.config.js .
COPY go.mod go.sum main.go ./
COPY internal internal

RUN go get
RUN go build

COPY static static
COPY content content
COPY templates templates

CMD ./htmx-blog
