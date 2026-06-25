{
  description = "Bark flake";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system}; in
      {
        devShells.default = pkgs.mkShell {
          name = "bark-devShell";
          packages = with pkgs; [
            go
            gopls
          ];
          env.LOG_LEVEL = "error";
          shellHook = ''
            export PATH="$PATH:$(go env GOPATH)/bin"
          '';
        };
      }
    );
}
