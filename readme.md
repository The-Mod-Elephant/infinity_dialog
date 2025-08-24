# Infinity Dialog
![](https://img.shields.io/badge/go-65A2BE2?logo=go&style=for-the-badge&logoColor=grey)
[![](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)](https://github.com/The-Mod-Elephant/infinity_dialog/releases/latest)
[![](https://img.shields.io/badge/Windows-0078D6?&style=for-the-badge&logoColor=white&logo=git-for-windows)](https://github.com/The-Mod-Elephant/infinity_dialog/releases/latest)
[![](https://img.shields.io/badge/mac%20os-grey?style=for-the-badge&logo=apple&logoColor=white)](https://github.com/The-Mod-Elephant/infinity_dialog/releases/latest)
[![](https://img.shields.io/github/actions/workflow/status/The-Mod-Elephant/infinity_dialog/main.yaml?style=for-the-badge)](https://github.com/The-Mod-Elephant/infinity_dialog/actions/workflows/main.yaml)
[![](https://img.shields.io/github/license/The-Mod-Elephant/infinity_dialog?style=for-the-badge)](./LICENSE)

This tool has 3 features, [Missing](#missing),[Traverse](#traverse), [Discover](#discover) and [View](#view).

- [Missing](#missing) shows missing string for tra, files, ie if a string is missing for english but not in french. This can also fix by placing extra strings.
- [Discover](#discover) shows all the strings for all languages and where they are set.
- [Traverse](#traverse) walks over a bg mod and shows you the tree of paths through the files.
- [View](#view) can read a collection of infinity file formats and text files involved in weidu mods.

### Traverse

This feature lets you see the paths through your mod by parsing the area + baf files.

#### Demo

![](./docs/traverse.gif)

### Discover

This feature interacts with .tra files used in infinity engine mods. Currently we walk a directory of a mod and allow you to browse the strings.

#### Demo

![](./docs/discover.gif)

### View

This feature read infinity engine files and text files. Here is a list of infinity engine file extensions is supported:
- .are
- .bam
- .baf
- .cre
- .dlg
- .eff
- .itm
- .spl

#### Demo

![](./docs/view.gif)


## Run

```sh
go run ./main.go
```

## Build

```sh
go build -o . ./...
```

## Test

```sh
go test ./...
```

## Debug

```sh
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug --headless --api-version=2 --listen=0.0.0.0:2345 || reset
```
Then attach in vscodium
