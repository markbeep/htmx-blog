FROM golang:1.21-alpine3.18
WORKDIR /app

# Download tailwindcss
RUN apk add --no-cache curl
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.5/tailwindcss-linux-arm64
RUN chmod +x tailwindcss-linux-arm64
RUN mv tailwindcss-linux-arm64 /bin/tailwindcss
COPY static/tw.css tw.css
COPY templates templates
COPY tailwind.config.js .

COPY go.mod go.sum main.go ./
RUN go get
RUN go build

COPY static static
COPY content content
COPY templates templates


CMD ./htmx-blog
