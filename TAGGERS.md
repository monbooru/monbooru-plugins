# Tagger registry

see [CONTRIBUTING](CONTRIBUTING.md).

Every auto-tagger model monbooru supports, one folder per model under
[`taggers/`](taggers). monbooru already carries some of them on its own
(wd-swinv2, animetimm-eva02, joytag, pixai-v1), so they already show in monbooru's catalog in Settings > Auto-Tagger. For any other, copy its
folder into monbooru's models folder.

| name              | what it tags                           | size   | login | maintainer | source                                                                          | monbooru |
|-------------------|----------------------------------------|--------|-------|------------|---------------------------------------------------------------------------------|----------|
| wd-swinv2         | anime, the WD v3 label set             | 330 MB | no    | leqwin     | [SmilingWolf](https://huggingface.co/SmilingWolf/wd-swinv2-tagger-v3)           | v1.22.2      |
| animetimm-eva02   | anime, a newer danbooru tag set        | 1.3 GB | yes   | leqwin     | [animetimm](https://huggingface.co/animetimm/eva02_large_patch14_448.dbv4-full) | v1.22.2      |
| joytag            | photos and realistic images            | 430 MB | no    | leqwin     | [fancyfeast](https://huggingface.co/fancyfeast/joytag)                          | v1.22.2     |
| camie-v2          | anime, with artists and copyrights     | 789 MB | no    | leqwin     | [Camais03](https://huggingface.co/Camais03/camie-tagger-v2)                     | v1.22.2      |
| pixai-v1          | anime, the newest tag set with artists | 2 GB   | no    | leqwin     | [noaione](https://huggingface.co/noaione/pixai-tagger-v1.0-onnx)                | v1.22.2      |
| wd-eva02-large-v3 | anime, the WD v3 label set             | 1.3 GB | no    | leqwin     | [SmilingWolf](https://huggingface.co/SmilingWolf/wd-eva02-large-tagger-v3)      | v1.22.2      |

`login` says whether the model page needs a Hugging Face login first.
`monbooru` is the last version the maintainer tested against.

## Row rules

- One row here and one folder in `taggers/`.
- Test it first: drop the folder and the weights in a model folder, run
  it on a few images, and tune the thresholds in `model.json`.
- The weights and label files stay where their authors publish them;
  `files` pins them to a revision and `source` links the model page.
- An entry whose model is gone will be removed.
