# monbooru-plugins

<picture>
  <img src=".github/banner-mascot.png" alt="" width="160">
</picture>

Registry of plugins, themes and auto-tagger models for [monbooru](https://github.com/monbooru/monbooru). Themes also work in [monloader](https://github.com/monbooru/monloader).

- **[PLUGINS.md](PLUGINS.md)**: plugins
- **[THEMES.md](THEMES.md)**: themes
- **[TAGGERS.md](TAGGERS.md)**: auto-tagger models

## Before you install a plugin or a theme

Being listed here is not a code review. Nothing has been verified: read the source, then build it yourself from the source you read. That's why every row has to link a public repo.

## Installing

`<configdir>` is the folder that holds `monbooru.toml` (`/config` on Docker).

- **Plugin**: drop its folder (the one with a `plugin.toml`) in `<configdir>/plugins/` and press **Enable** on its row in **Settings -> Plugins**, or run it yourself on your network, pointed at monbooru. Either way, press **Approve** under **Pairing requests** in **Settings -> Plugins**. monbooru provides no runtime, so a plugin that isn't a self-contained binary should include in its README what it needs.
- **Theme**: drop its folder in `<configdir>/themes/` (a bare `.css` file works for CSS-only themes) and pick it in **Settings -> Plugins**, or in monloader's **Settings**.
- **Tagger**: copy its folder from [`taggers/`](taggers) into monbooru's models folder, then press **Install** on its row in **Settings -> Auto-Tagger** to fetch the weights.

## In this repo

- [`themes/dark`](themes/dark): starter dark theme with every variable, commented; copy the folder and start recoloring.
- [`themes/light`](themes/light): the shipped light theme.
- [`themes/legacy-dark`](themes/legacy-dark), [`themes/legacy-light`](themes/legacy-light): the legacy themes, with the previous logo and no drawings.
- [`taggers/`](taggers): the auto-tagger models.
- [`examples/simple-edit`](examples/simple-edit): an example plugin in Go, with rotate as a relay button that supports batch actions, and crop as an open-mode page.

## Reporting problems

See [Report it](https://github.com/monbooru/monbooru/blob/main/CONTRIBUTING.md#report-it) in monbooru's CONTRIBUTING: it says which repo each kind of problem goes to, plugins, themes and taggers included.

## Getting listed

Open a PR that adds one row to [PLUGINS.md](PLUGINS.md), [THEMES.md](THEMES.md) or [TAGGERS.md](TAGGERS.md) (with the tagger's folder). See [CONTRIBUTING](CONTRIBUTING.md).
