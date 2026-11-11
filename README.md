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

# then, when already on main with local changes you want to keep
latest-main --force
```

### Commands

| Command | Description |
|---|---|
| `checkout-main-and-update [--force]` | Checkout `main`/`master` and pull latest. Errors if working tree is dirty. `--force` skips the clean-tree check and pulls in place, but only when already on `main`/`master` — it never switches branches. |
| `setup-gitignore-local` | Point the repo's `core.excludesFile` at a root `.gitignore.local` and create it (self-ignoring) so local-only ignore patterns never push to the remote. |

## Build

```sh
make build
```
