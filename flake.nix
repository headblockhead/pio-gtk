{
  description = "RP2040 PIO emulator frontend";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/26.05";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
    };
    gotk4-nix = {
      url = "github:diamondburned/gotk4-nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.follows = "flake-utils";
      inputs.gomod2nix.follows = "gomod2nix";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      gomod2nix,
      gotk4-nix,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        pname = "pio-gtk";
        pio-gtk = (pkgs.extend gomod2nix.overlays.default).buildGoApplication {
          inherit pname;
          version = self.shortRev or self.dirtyShortRev or "dev";
          src = ./.;
          modules = ./gomod2nix.toml;
          nativeBuildInputs = with pkgs; [
            pkg-config
            wrapGAppsHook4
            git
          ];
          buildInputs = with pkgs; [
            gtk4
            libadwaita
            glib
            gdk-pixbuf
            librsvg
            gobject-introspection
            hicolor-icon-theme
          ];
          doCheck = false;
        };

        gotk4 = gotk4-nix.lib.mkLib {
          inherit pkgs;
          base = {
            inherit pname;
            buildInputs = pkgs: [ pkgs.libadwaita ];
          };
        };
      in
      {
        packages = {
          inherit pio-gtk;
          default = pio-gtk;
        };

        apps.default = {
          type = "app";
          program = "${pio-gtk}/bin/${pname}";
        };

        devShells.default = gotk4.mkShell {
          buildInputs = [
            gomod2nix.packages.${system}.default
            pkgs.libadwaita
          ];
          clangdPackages =
            pkgs: with pkgs; [
              gtk4
              glib
              libadwaita
            ];
        };
      }
    );
}
