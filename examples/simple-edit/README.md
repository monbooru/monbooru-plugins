# simple-edit

Rotate and crop images inside monbooru.

It exists as a reference plugin, showing the whole surface working at once.
You can copy the directory to start your own.

- **pairing** - offers itself on startup, waits for you to approve, and stores
  the two secrets it ends up holding.
- **a relay button** - `rotate 90`, on both the image page and the gallery
  batch bar. monbooru posts it the images in scope, the plugin turns each one
  and writes it back;
- **an open-mode button** - `crop`, a link to a page this plugin serves, opened
  in a pop-in over monbooru. Saving either replaces the image or pushes a new
  one and declares it a newer version of the original.

All three declare `media = "image"`, so monbooru keeps them off videos, comic
archives and animations rather than offering an edit that could only fail.

Every change it makes goes through monbooru's REST API. It never touches the
database or the files on disk directly.

## The three files

| file | what is in it |
|---|---|
| `plugin.go` | common mechanics for any plugin: pairing and its teardown, the credentials it leaves behind, calling the API, and checking that a request really is monbooru |
| `main.go` | specifics of this plugin: the buttons it declares, the rotate, the crop |
| `crop.html` | the crop page, compiled in with `go:embed` |

Standard library apart from `golang.org/x/image/webp`, because monbooru takes
webp files and the standard library cannot read one.

## Building it

```
go build -o simple-edit .
```

Ship each binary next to `plugin.toml` and `README.md`. Nothing else is
required in the folder.

## Running it

Put this folder in `<configdir>/plugins/` (next to `monbooru.toml`) with the
binary built, start it from its row under **Settings -> Plugins**, and approve
the pairing request that follows.

You can also run it by hand, without putting it in `plugins/`:

```
./simple-edit -monbooru http://127.0.0.1:8080 -addr 127.0.0.1:7301
```

## Changing the buttons

The buttons a plugin declares are what the operator approved, so changing
`buttons` in `main.go` means pairing again: remove the pairing under
**Settings -> Plugins** and approve the card that comes straight back. 

## What to copy

The contract is plain HTTP and JSON, so any language works. In Go, `plugin.go`
carries over as it stands. 