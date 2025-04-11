
{ pkgs }: {
  deps = [
    pkgs.nodejs_22
    pkgs.go_1_21
    pkgs.nodejs_20
    pkgs.gopls
    pkgs.gotools
  ];
}
