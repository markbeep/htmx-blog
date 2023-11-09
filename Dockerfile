FROM nixos/nix
WORKDIR /app

RUN nix-env -iA nixpkgs.nixFlakes nixpkgs.curl
RUN echo "experimental-features = flakes nix-command" >> /etc/nix/nix.conf

COPY tailwind.config.js .
COPY go.mod go.sum main.go ./
COPY internal internal
COPY components components
COPY flake.lock flake.nix ./

RUN nix build .#htmx-blog

COPY static static
RUN nix run .#tailwindcss -- -i static/tw.css -o static/main.css --minify

COPY content content

CMD ./result/bin/htmx-blog
