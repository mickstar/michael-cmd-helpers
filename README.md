# michael-cmd-helpers

Personal CLI utility for common git/dev workflows.

## Install

```sh
make install
```

Installs `michael-cmd` to `~/.local/bin`.

## Usage

```sh
# add to your shell config
alias latest-main="michael-cmd checkout-main-and-update"
```

### Commands

| Command | Description |
|---|---|
| `checkout-main-and-update` | Checkout `main`/`master` and pull latest. Errors if working tree is dirty. |

## Build

```sh
make build
```
