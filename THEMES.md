# Theme registry

see [CONTRIBUTING](CONTRIBUTING.md).

A theme is a folder holding `theme.css` and optionally `logo.png`, dropped
in `<configdir>/themes/` and picked in Settings > Plugins.

| name    | the look                                              | logo | maintainer | source                           | monbooru |
|---------|-------------------------------------------------------|------|------------|----------------------------------|----------|
| light   | a simple light theme | yes   | leqwin     | [themes/light](themes/light)     | v1.17.1  |
| dark | the default shipped dark theme, every variable commented     | no   | leqwin     | [themes/dark](themes/dark) | v1.17.1  |
| lainbooru | serial experiments lain theme, indigo grounds with magenta accent | yes  | gary-host-laptop | [lainbooru-theme](https://github.com/gary-host-laptop/lainbooru-theme) | v1.17.1  |

`logo` says whether picking the theme also changes the topbar logo.  
`monbooru` is the last version the maintainer tested against.

## Row rules

- Link public source, same as the plugin registry.
- Rule-level selectors are explicitly
  unstable, so a theme that overrides more than the `:root` variables could break
  on any release.
- One row per theme, one line for the look.
- An entry whose repo is gone will be removed.
