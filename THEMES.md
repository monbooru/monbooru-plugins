# Theme registry

see [CONTRIBUTING](CONTRIBUTING.md).

A theme is a folder holding `theme.css` and optionally `logo.png` and
`favicon.png`, dropped in `<configdir>/themes/` and picked in
Settings > Plugins.

| name    | the look                                              | logo | favicon | maintainer | source                           | monbooru |
|---------|-------------------------------------------------------|------|---------|------------|----------------------------------|----------|
| light   | a simple light theme | yes   | yes     | leqwin     | [themes/light](themes/light)     | v1.20.1  |
| dark | the default shipped dark theme, every variable commented     | no   | no      | leqwin     | [themes/dark](themes/dark) | v1.20.1  |
| lainbooru | serial experiments lain theme, indigo grounds with magenta accent | yes  | no      | gary-host-laptop | [lainbooru-theme](https://github.com/gary-host-laptop/ghost-themes/tree/main/lainbooru) | v1.20.1  |
| old_steam | olive palette inspired by the old steam client | yes  | no      | gary-host-laptop | [old_steam](https://github.com/gary-host-laptop/ghost-themes/tree/main/old_steam) | v1.20.1  |

`logo` and `favicon` say whether picking the theme also changes the topbar logo and the tab icon. 
`monbooru` is the last version the maintainer tested against.

## Row rules

- Link public source, same as the plugin registry.
- Rule-level selectors are explicitly
  unstable, so a theme that overrides more than the `:root` variables could break
  on any release.
- One row per theme, one line for the look.
- An entry whose repo is gone will be removed.
