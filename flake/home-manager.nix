self: { config, lib, pkgs, ... }:
let
  cfg = config.programs.bark;
  tomlFormat = pkgs.formats.toml { };
in
{
  options.programs.bark = {
    enable = lib.mkEnableOption "bark";

    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.stdenv.hostPlatform.system}.default;
      defaultText = lib.literalExpression "bark";
    };

    settings = lib.mkOption {
      type = tomlFormat.type;
      default = { };
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [ cfg.package ];

    xdg.configFile."bark/bark.toml" = lib.mkIf (cfg.settings != { }) {
      source = tomlFormat.generate "bark.toml" cfg.settings;
    };
  };
}
