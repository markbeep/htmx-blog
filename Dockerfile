FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY internal internal
COPY components components
RUN go run github.com/a-h/templ/cmd/templ@v0.2.747 generate
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/htmx-blog .

FROM node:24-alpine AS styles
WORKDIR /src
COPY package.json package-lock.json uno.config.ts ./
RUN npm ci --ignore-scripts
COPY components components
COPY static/uno.css static/uno.css
RUN CI=1 npx unocss -c uno.config.ts --minify

FROM alpine:3.22 AS images
RUN apk add --no-cache libwebp-tools
WORKDIR /src
COPY content content
COPY scripts/optimize-images.sh ./
RUN sh optimize-images.sh content

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /out/htmx-blog ./htmx-blog
COPY --from=build /src/components/index.xml ./components/index.xml
COPY static ./static
COPY --from=styles /src/static/main.css ./static/main.css
COPY --from=images /src/content ./content
ENV ROOT_PATH=/app
EXPOSE 3000
CMD ["./htmx-blog"]
