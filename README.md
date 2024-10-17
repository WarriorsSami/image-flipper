# image-flipper

Image flipper implemented as a Flutter GUI and a Go TUI with support for:
- horizontal, vertical and horizontal-vertical (both) flipping

### Input images folder

![Before flipping](assets/before-flipper.png)

### Output images folder (after flipping)

![After flippinng](assets/after-flipper.png)

## GUI

Implemented using Flutter Cubit (BLoC related state management package)

![Flipper GUI](assets/flipper-gui.png)

## TUI

Implemented using Golang Bubbletea (framework for building TUI apps), Huh forms, errgroups and pipeline pattern for image processing.

![Flipper TUI](assets/flipper-tui.png)

## CLI

Implemented using Golang Cobra (framework for building CLI apps), errgroups and pipeline pattern for image processing.

![Flipper CLI](assets/flipper-cli.png)
