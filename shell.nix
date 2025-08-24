{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    git
    go
    golangci-lint
    pre-commit
  ];
}
