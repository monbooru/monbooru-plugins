# Plugin registry

see [CONTRIBUTING](CONTRIBUTING.md).

> **Assume nothing in this table has been reviewed.** 
> Read the source and build it yourself before you run it:


| name        | what it does                               | maintainer | source                                       | contract | monbooru |
|-------------|--------------------------------------------|------------|----------------------------------------------|----------|----------|
| simple-edit | rotate and crop images (Go, single binary) | leqwin     | [examples/simple-edit](examples/simple-edit) | 1        | v1.22.1  |

`source` is a public repo. `contract` is the monbooru plugin contract version the plugin sends in its pairing request; a monbooru that uses an older one refuses it, so a row older than the running monbooru is not by itself a warning. `monbooru` is the last monbooru version the maintainer tested against (for information only).

## Row rules

- Link public source. Binary-only listings are refused. 
- One row per plugin, one line, description short enough to read at a
  glance.
- If it is not a self-contained program, say what has to be installed
  first. monbooru provides no runtime, and the Docker image has no
  interpreter and no package manager.
- An entry whose repo is gone will be removed.