
{ pkgs }: {
  deps = [
    pkgs.go_1_21
    pkgs.gopls
    pkgs.nodejs_20
    pkgs.gotools
  ];
}
