{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    templ.url = "github:a-h/templ";
    gitignore = {
      url = "github:hercules-ci/gitignore.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, flake-utils, nixpkgs, gitignore, templ }:
    flake-utils.lib.eachDefaultSystem (system: 
      let
        pkgs = import nixpkgs { inherit system; };
        templ-pkg = templ.packages.${system}.templ;
      in
      {
        packages = rec {
          htmx-blog = pkgs.buildGo121Module {
            name = "htmx-blog";
            src = gitignore.lib.gitignoreSource ./.;
            vendorHash = "sha256-HRrW8wRYqnyDtcCWVt9YM2DmlutCV2KWvVM4jM4DqdY=";

            preBuild = ''
              ${templ-pkg}/bin/templ generate
            '';
            postBuild = ''
              ${pkgs.tailwindcss}/bin/tailwindcss -- -i static/tw.css -o static/main.css --minify
              mkdir -p $out/static $out/content $out/components
              cp -r static/* $out/static
              cp -r content/* $out/content
              cp -r components/index.xml $out/components/index.xml
            '';
          };

          tailwindcss = pkgs.tailwindcss;

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "htmx-blog";
            config = {
              Cmd = ["${htmx-blog}/bin/htmx-blog"];
              ExposedPorts = {
                "3000/tcp" = {};
              };
              Env = [
                "ROOT_PATH=${htmx-blog}"
              ];
            };
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ 
              go
              tailwindcss
              nodejs_20
              playwright-test
              k6
              templ-pkg
          ];
        };
      }
    );
}
