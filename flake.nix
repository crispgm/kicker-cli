{
  inputs = {
    #flake-utils.url = "github:numtide/flake-utils";
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs = {
    self,
    nixpkgs,
    flake-parts,

  } @ inputs:
  flake-parts.lib.mkFlake {inherit inputs;} {
    systems = [
      "aarch64-darwin"
      "aarch64-linux"
      "x86_64-darwin"
      "x86_64-linux"
    ];
    
    perSystem = {
      pkgs,
      system,
      ...
    }:
    {
      packages = {
        default = pkgs.callPackage ./pkg.nix {};
                     
      };
    };
  };
}
