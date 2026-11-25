{
  description = "Bark flake";

  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils }:
    (flake-utils.lib.eachDefaultSystem (system:
      let pkgs = nixpkgs.legacyPackages.${system}; in
      {
        packages.default = pkgs.buildGoModule {
          pname = "bark";
          version = "0-unstable";
          src = ./.;
          vendorHash = "sha256-38bWTAdoXO5qrnkPBKv5G5Qq4SJSKsoJg7Kv6FhCqPk=";
          nativeCheckInputs = [ pkgs.git ];

          nativeBuildInputs = [ pkgs.installShellFiles pkgs.git ];

          postInstall = pkgs.lib.optionalString
            (pkgs.stdenv.buildPlatform.canExecute pkgs.stdenv.hostPlatform) ''
              $out/bin/bark --generate man-pages
              installManPage *.1
            '';
        };

        devShells.default = pkgs.mkShell {
          name = "bark-devShell";
          packages = with pkgs; [
            go
            gopls
          ];
          env.LOG_LEVEL = "error";
          shellHook = ''
            export PATH="$(go env GOPATH)/bin:$PATH"
          '';
        };
      }
    )) // {
      homeManagerModules.default = import ./flake/home-manager.nix self;
    };
}
