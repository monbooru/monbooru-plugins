# Plugin registry

see [CONTRIBUTING](CONTRIBUTING.md).

> **Assume nothing in this table has been reviewed.** 
> Read the source and build it yourself before you run it:


| name        | what it does                               | maintainer | source                                       | monbooru |
|-------------|--------------------------------------------|------------|----------------------------------------------|----------|
| simple-edit | rotate and crop images (Go, single binary) | leqwin     | [examples/simple-edit](examples/simple-edit) | v1.17.1  |

`source` is a public repo. `monbooru` is the last monbooru version the maintainer tested against (for information only, non blocking).

## Row rules

- Link public source. Binary-only listings are refused. 
- One row per plugin, one line, description short enough to read at a
  glance.
- If it is not a self-contained program, say what has to be installed
  first. monbooru provides no runtime, and the Docker image has no
  interpreter and no package manager.
- An entry whose repo is gone will be removed.