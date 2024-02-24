{
  description = "glsl learning project";
  inputs = {
    #flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs = { self, nixpkgs, nvim-conf,nixvim}:
  let
    pkgs = nixpkgs.legacyPackages.x86_64-linux.pkgs;
    
    
  in{
    devShells.x86_64-linux.default = pkgs.mkShell {
      buildInputs = [
        pkgs.go
        pkgs.gopls
      ];

      shellHook = ''
        export inNixShell=1
        echo "hello World" ...
      '';
    };
  };  
}
