{
  # Only development environment since building docker images from Nix would be
  # inconvenient for the team
  description = "L0 development flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = {
    self,
    nixpkgs,
  }: let
    system = "x86_64-linux";
    pkgs =
      import nixpkgs {inherit system;};

    version = "0.0.1";

    scripts = let
      colors = {
        reset = "\\033[0m";
        red = "\\033[0;31m";
        green = "\\033[0;32m";
      };
      echo_success = ''echo -e "${colors.green}Success!${colors.reset}"'';
      echo_failed = ''echo -e "${colors.red}Failed${colors.reset}"'';
    in
      builtins.attrValues (builtins.mapAttrs pkgs.writeShellScriptBin {
        run-tests = ''
          go test ./... -coverprofile=cover.out
          go tool cover -html=cover.out -o cover.html
          rm cover.out
        '';
      });
    back = with pkgs; [go_1_20 gopls delve];
    utils = with pkgs; [insomnia dbeaver plantuml];
  in {
    devShells.${system}.default = pkgs.mkShell {
      name = "WB-L0";
      buildInputs = back ++ scripts ++ utils;
      CGO_ENABLED = 0; # delve from nixpkgs refuses to work otherwise
    };
  };
}
