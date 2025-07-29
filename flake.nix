{
  description = "golang playground";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            golangci-lint

            mockgen
            sqlc

            # protobuf
            buf
            protoc-gen-go
            protoc-gen-connect-go
          ];
        };
      }
    );
}
