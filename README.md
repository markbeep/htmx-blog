# Mark's Blog

![Playwright Tests](https://img.shields.io/github/actions/workflow/status/markbeep/htmx-blog/playwright.yml?logo=playwright&logoColor=%23FFFFFF&label=Playwright%20Tests&link=https%3A%2F%2Fgithub.com%2Fmarkbeep%2Fhtmx-blog%2Factions)
![Deploy Dockerfile](https://img.shields.io/github/actions/workflow/status/markbeep/htmx-blog/build.yml?logo=kubernetes&logoColor=%23FFFFFF&label=Deploy&link=https%3A%2F%2Fgithub.com%2Fmarkbeep%2Fhtmx-blog%2Factions)
![Go version](https://img.shields.io/github/go-mod/go-version/markbeep/htmx-blog?logo=go)

This is my [personal blog](https://markc.su). It started out as a simple project to practice using [htmx](https://htmx.org/) to make
it an extremely lightweight and efficient website and [tailwindcss](https://tailwindcss.com/) for the styling.

One important requirement was that all my blog posts were to be written in markdown (in `/content`). It then turned out to be
a bigger challenge to get the markdown file generation working than actually creating the blog. I first started out using [pandoc](https://pandoc.org/)
to generate the markdown externally, but pandoc does some very weird styling and the syntax highlighting colors they support look
really bad. I instead then opted for [goldmark](https://github.com/yuin/goldmark) which is also in use by the [hugo](https://gohugo.io/).

## Installation

To get the blog running locally you first need `go 1.21`. You'll also need [templ](https://templ.guide/quick-start/installation) which has three quick ways to install:

```bash
go install github.com/a-h/templ/cmd/templ@latest # install via go
curl -o https://github.com/a-h/templ/releases/download/v0.2.432/templ_Linux_x86_64.tar.gz && chmod +x templ # install binary
nix run github:a-h/templ # run using nix
```

Once templ and go is set, you can generate the template go files and then start the server:

```bash
templ generate
go run .
```

### Hot Reload

For the best developer experience I recommend getting [air](https://github.com/cosmtrek/air) which enables hot-reloading on save. Install air and run it:

```bash
go install github.com/cosmtrek/air@latest
air
```

### Styling

If you're working on the css or styling, it is recommended to also get the [tailwindcss CLI](https://tailwindcss.com/blog/standalone-cli) tool and run it on the side to keep on updating
the `main.css` stylesheet:

```bash
tailwindcss -i static/tw.css -o static/main.css --watch
```

### Docker

It's also possible to start up the website using [Docker](https://www.docker.com/):

```bash
docker compose up --build
```

This will build and start up the website as well as run the playwright tests. If you're building for arm64, you'll need to modify the Dockerfile and swap out
the tailwindcss CLI tool installation for the arm64 version.

## Testing

For extensive testing, [Playwright](https://playwright.dev/) is used. To make the install simpler, there's a Dockerfile
in `e2e/Dockerfile.playwright`. Build and run it from the **root directory**.

```bash
docker build -t htmx-playwright -f e2e/Dockerfile.playwright .
docker run --rm -it --network host htmx-playwright /bin/bash
```

## Nix

With the `flake.nix` and `flake.lock` files you can set up local development with the exact same dependency versions.

```bash
nix develop # downloads all required dependencies
nix build .#htmx-blog # build the blog
nix run .#htmx-blog # build and run the blog locally
nix run .#tailwindcss -- -i static/tw.css -o static/main.css --watch
```

**Note:** Tailwindcss still needs to be executed manually on the side to generate the `static/main.css` file

# Creating Posts

Posts can be created by executing the `post.sh` script. Simply execute the following:

```sh
./post.sh content/posts/name-of-post
```
