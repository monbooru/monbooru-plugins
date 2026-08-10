# monbooru-plugins

Registry for [monbooru](https://github.com/monbooru/monbooru)
plugins and themes.

- **[PLUGINS.md](PLUGINS.md)** - plugins registry
- **[THEMES.md](THEMES.md)** - themes registry

## Before you run one

Listing here is not code review. Assume nothing has been verified, 
so read the source, then build it yourself from the source you read. This
is why every row has to link a public repo. 

Bugs in a listed plugin go to that plugin's own tracker. The monbooru
tracker takes contract issues only.

## Installing

- **app**: drop the plugin folder (one carrying a `plugin.toml`) in
  `<configdir>/plugins/` and start it from its row in Settings > Plugins,
  or run it yourself on your LAN pointed at monbooru. Either way, approve the pairing card in Settings > Plugins.  A dropped
  folder brings whatever it needs to run with it (monbooru installs
  nothing and provides no runtime) so a plugin that is not a
  self-contained binary says in its README what needs to be on the machine
  first.
- **theme**: drop the theme folder in `<configdir>/themes/` (a bare
  `.css` works for CSS-only themes) and pick it in Settings > Plugins. 
  Your own `custom_css` and `server.logo` still win over any theme.

## In this repo

- [`themes/dark`](themes/dark) - starter dark theme with every variable and comments; copy the folder and start recoloring.
- [`themes/light`](themes/light) - a light theme example.
- [`examples/simple-edit`](examples/simple-edit) - an example plugin in go:
  rotate as a relay button with batch action support, crop as an open-mode page.

## Getting listed

PR one row to [PLUGINS.md](PLUGINS.md) or [THEMES.md](THEMES.md).  
See [CONTRIBUTING](CONTRIBUTING.md).
