# Infinity Dialog
![](https://img.shields.io/badge/go-65A2BE2?logo=go&style=for-the-badge&logoColor=grey)
[![](https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black)](https://github.com/dark0dave/infinity_dialog/releases/latest)
[![](https://img.shields.io/badge/Windows-0078D6?&style=for-the-badge&logoColor=white&logo=git-for-windows)](https://github.com/dark0dave/infinity_dialog/releases/latest)
[![](https://img.shields.io/badge/mac%20os-grey?style=for-the-badge&logo=apple&logoColor=white)](https://github.com/dark0dave/infinity_dialog/releases/latest)
[![](https://img.shields.io/github/actions/workflow/status/dark0dave/infinity_dialog/main.yaml?style=for-the-badge)](https://github.com/dark0dave/infinity_dialog/actions/workflows/main.yaml)
[![](https://img.shields.io/github/license/dark0dave/infinity_dialog?style=for-the-badge)](./LICENSE)


## Demo

![](./docs/example.gif)

## Run

```sh
go run ./main.go
```

## Build

```sh
go build ./...
```

## Test

```sh
go test ./...
```

## Debug

```sh
dlv debug --headless --api-version=2 --listen=127.0.0.1:2345
```
Then attach in vscodium
