with import <nixpkgs> {};
mkShellNoCC {
  packages = with pkgs; [ pre-commit gcc go ];
  env.CGO_ENABLED = "1";
}
