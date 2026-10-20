# Mark's Blog

![Playwright Tests](https://img.shields.io/github/actions/workflow/status/markbeep/htmx-blog/playwright.yml?logo=playwright&logoColor=%23FFFFFF&label=Playwright%20Tests&link=https%3A%2F%2Fgithub.com%2Fmarkbeep%2Fhtmx-blog%2Factions)
![Deploy Dockerfile](https://img.shields.io/github/actions/workflow/status/markbeep/htmx-blog/build.yml?logo=kubernetes&logoColor=%23FFFFFF&label=Deploy&link=https%3A%2F%2Fgithub.com%2Fmarkbeep%2Fhtmx-blog%2Factions)
![Go version](https://img.shields.io/github/go-mod/go-version/markbeep/htmx-blog?logo=go)

This is my [personal blog](https://rkma.dev). It is a lightweight and efficient static site built with Go, [templ](https://templ.guide/),
and [UnoCSS](https://unocss.dev/).

One important requirement was that all my blog posts were to be written in markdown (in `/content`). It then turned out to be
a bigger challenge to get the markdown file generation working than actually creating the blog. I first started out using [pandoc](https://pandoc.org/)
to generate the markdown externally, but pandoc does some very weird styling and the syntax highlighting colors they support look
really bad. I instead then opted for [goldmark](https://github.com/yuin/goldmark) which is also in use by the [hugo](https://gohugo.io/).

## Installation

To get the blog running locally you first need a bunch of tools. The easiest way to get up and running is with [mise](https://mise.jdx.dev/)
where you can install and activate the required tools with `mise activate`. If you want to manually install the tools
look into the `mise.toml` file for the exact versions required.

Justfile makes it easier to _just_ run the required programs. Look into the `justfile` file to see the full commands.

1. Install the frontend dependencies:

   ```sh
   npm install
   ```

2. Run templ to generate the go files that are then required:

   ```sh
   templ generate
   ```

3. Run the webserver:

   ```sh
   just d
   ```

4. In a separate terminal, run UnoCSS to continuously update the CSS file:
   ```sh
   just css-watch
   ```

Blog post images are converted to WebP at build time. To do the same locally (requires `cwebp`/`gif2webp` from libwebp or ImageMagick):

```sh
just webp
```

### Docker

It's also possible to start up the website using [Docker](https://www.docker.com/):

```bash
docker compose up --build
```

The image is built in multiple Alpine stages: Go builds the server, Node runs UnoCSS, and libwebp converts post images.
The container serves brotli when the client accepts it, falls back to gzip, and otherwise sends uncompressed responses.

## Testing

For extensive testing, [Playwright](https://playwright.dev/) is used. To make the install simpler, there's a Dockerfile
in `e2e/Dockerfile.playwright`. Build and run it from the **root directory**.

```bash
docker build -t htmx-playwright -f e2e/Dockerfile.playwright .
docker run --rm -it --network host htmx-playwright /bin/bash
```

# Creating Posts

Posts can be created with the `just new` recipe. Simply execute the following:

```sh
just new content/posts/name_of_post
```
