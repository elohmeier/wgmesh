{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;

          overlays = [ ];
        };
        buildDeps = with pkgs; [ ];
        devDeps = with pkgs; buildDeps ++ [
        ];
      in
      {
        devShell = pkgs.mkShell {
          packages = with pkgs; [
            go
            gnumake
            golangci-lint
            gotestsum
            golint
            goreleaser
            protobuf
            protoc-gen-go
            protoc-gen-go-grpc
            openapi-generator-cli
            ent
          ];
        };
      });
}
