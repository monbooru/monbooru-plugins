# Theme registry

see [CONTRIBUTING](CONTRIBUTING.md).

A theme is a folder holding `theme.css` and optionally `logo.png`,
`favicon.png` and the mascot drawings
(`mascot-hello.png`, `mascot-lost.png`, `mascot-sleep.png`, `mascot-panic.png`,
`mascot-sparkle.png`, `mascot-sparkle-head.png`, `mascot-sleep-head.png`),
dropped in `<configdir>/themes/` and picked in monbooru's Settings ->
Plugins or monloader's Settings. 
A theme without drawings shows none, and a theme without favicon and/or logo shows the default one.
The same folder works in both monbooru and monloader.

| name         | the look                                                            | logo | favicon | drawings | maintainer       | source                                                                                   | monbooru |
|--------------|---------------------------------------------------------------------|------|---------|----------|------------------|------------------------------------------------------------------------------------------|----------|
| light        | the shipped light theme, red accent                                 | yes  | yes     | yes      | leqwin           | [themes/light](themes/light)                                                             | v1.23.0  |
| dark         | the default shipped dark theme, every variable commented            | yes  | yes     | yes      | leqwin           | [themes/dark](themes/dark)                                                               | v1.23.0  |
| legacy-dark  | the legacy dark theme                                               | yes  | yes     | no       | leqwin           | [themes/legacy-dark](themes/legacy-dark)                                                 | v1.23.0  |
| legacy-light | the legacy light theme                                              | yes  | yes     | no       | leqwin           | [themes/legacy-light](themes/legacy-light)                                               | v1.23.0  |
| lainbooru    | serial experiments lain theme, indigo grounds with magenta accent   | yes  | no      | no       | gary-host-laptop | [lainbooru-theme](https://github.com/gary-host-laptop/ghost-themes/tree/main/lainbooru)  | v1.20.1  |
| old_steam    | olive palette inspired by the old steam client                      | yes  | no      | no       | gary-host-laptop | [old_steam](https://github.com/gary-host-laptop/ghost-themes/tree/main/old_steam)        | v1.20.1  |

`logo` and `favicon` say whether picking the theme also changes the topbar logo and the tab icon;
`drawings` whether it ships the mascot drawings the apps show on empty pages.
`monbooru` is the last version the maintainer tested against.

## Row rules

- Link public source, same as the plugin registry.
- Rule-level selectors are explicitly
  unstable, so a theme that overrides more than the `:root` variables could break
  on any release.
- One row per theme, one line for the look.
- An entry whose repo is gone will be removed.
