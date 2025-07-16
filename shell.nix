with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "go-env";
  buildInputs = [
    delve
    git
    gnupg
    pre-commit
    go
    go-critic
  ];
}
