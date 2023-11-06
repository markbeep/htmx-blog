FROM golang:1.21-alpine3.18
WORKDIR /app

# Download tailwindcss
RUN apk add --no-cache curl
RUN curl -o /bin/tailwindcss -sSL https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.5/tailwindcss-linux-x64
# RUN curl -o /bin/tailwindcss -sL https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.5/tailwindcss-linux-arm64
RUN chmod +x /bin/tailwindcss
RUN curl -o templ -sSL https://github.com/a-h/templ/releases/download/v0.2.432/templ_Linux_x86_64.tar.gz && chmod +x templ

COPY tailwind.config.js .
COPY go.mod go.sum main.go ./
COPY internal internal
COPY components components
RUN ./templ generate components

RUN go get
RUN go build

COPY static static
RUN tailwindcss -i static/tw.css -o static/main.css --minify

COPY content content


CMD ./htmx-blog
