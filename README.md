# koro

Pac-Man inspired prototype built with Go + Ebiten. The current gameplay already includes:

- Tile-based maze with pellets, power pellets, warp tunnels.
- Player character with grid-snapped movement and keyboard/touch controls.
- Randomised ghost AI with frightened mode and scoring/life rules.

## Run locally

```bash
go run ./cmd/game
```

Controls: Arrow keys or a gamepad D-pad.

## Mobile builds

Prerequisites:

- Go toolchain (1.25+).
- `ebitenmobile` installed (`go install github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile@latest`).
- Android: Android SDK/NDK configured in `ANDROID_HOME`/`ANDROID_NDK_HOME`.
- iOS: Xcode command line tools installed.

Targets:

```bash
# Generate Android AAR (outputs to build/koro.aar)
make mobile-android

# Generate iOS framework (outputs to build/Koro.framework)
make mobile-ios

# Remove build artifacts
make mobile-clean
```

The resulting bindings can be dropped into a native Android/iOS project to create an installable package.

## Next steps

- Replace placeholder colors with sprite art + sounds.
- Flesh out level progression and menus.
