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

To get the blog running locally you first need a bunch of tools. The easiest way to get up and running is with [mise](https://mise.jdx.dev/)
where you can install and activate the required tools with `mise activate`. If you want to manually install the tools
look into the `mise.toml` file for the exact versions required.

Justfile makes it easier to _just_ run the required programs. Look into the `justfile` file to see the full commands.

1. Run templ to generate the go files that are then required:

   ```sh
   templ generate
   ```

2. Run the webserver:

   ```sh
   just d
   ```

3. In a separate terminal, run tailwindcss to continously update the css files as required:
   ```sh
   just tw
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
just new content/posts/name_of_post
```
