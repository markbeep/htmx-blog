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
        packages = {
          htmx-blog = pkgs.buildGo121Module {
            name = "htmx-blog";
            src = gitignore.lib.gitignoreSource ./.;
            vendorHash = "sha256-MpBSY9I18hSw37E5ipfBEg6ZmIDyU70qV9zRB1LmK1E=";

            preBuild = ''
              ${templ-pkg}/bin/templ generate
            '';
          };

          tailwindcss = pkgs.tailwindcss;
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
