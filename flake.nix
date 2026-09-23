{
  description = "RP2040 PIO emulator frontend";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-26.05";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gotk4-nix.url = "github:diamondburned/gotk4-nix";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
      gomod2nix,
      gotk4-nix,
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
            gobject-introspection
            wrapGAppsHook4
          ];

          buildInputs = with pkgs; [
            gtk4
            libadwaita
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
        packages.default = pio-gtk;

        devShells.default = gotk4.mkShell {
          buildInputs = [
            gomod2nix.packages.${system}.default
          ];
        };
      }
    );
}
