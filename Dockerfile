FROM pandoc/core:3.1-alpine AS build-pages

WORKDIR /app
COPY content content
COPY templates/pandoc.html pandoc.html
COPY static/pandoc.css pandoc.css


# Creates the same directory structure
RUN mkdir -p generated && find content -type d -exec mkdir -p -- generated/{} \;
# Turns all markdown files into the equivalent html file
RUN find content -iname "*.md" -type f -exec sh -c 'pandoc "${0}" -s --template pandoc.html --css pandoc.css --highlight-style breezedark -o "./generated/${0%.md}.html"' {} \;
# Change all mathjax inline to inline-block, so they're actually inline
RUN find generated -iname "*.html" -type f -exec sh -c 'sed -i "s/math inline/math inline-block/g" ${0}' {} \;


FROM alpine:3.18 AS build-css
WORKDIR /app
RUN apk add --no-cache curl
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.5/tailwindcss-linux-arm64
RUN chmod +x tailwindcss-linux-arm64
COPY static/tw.css tw.css
COPY --from=build-pages /app/generated generated
COPY templates templates
COPY tailwind.config.js .
RUN ./tailwindcss-linux-arm64 -i tw.css -o main.css --minify


FROM golang:1.21-alpine3.18
WORKDIR /app

COPY go.mod go.sum main.go ./
RUN go get
RUN go build

COPY --from=build-pages /app/generated generated
COPY static static
COPY --from=build-css /app/main.css static/main.css
COPY content content
COPY templates templates


CMD ./htmx-blog
