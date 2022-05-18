# Gitignore File Generator

<details>
<summary>Table of Contents</summary>

- [Gitignore File Generator](#gitignore-file-generator)
  - [Usage](#usage)
  - [Features](#features)

</details>

<hr />

GI Gen is an open source CLI to generate .gitignore files. It is completely cross-platform, and
standalone (no dependencies).

You can run this CLI program to create or append a .gitignore template easily.

## Usage

Download the file for your platform in the [Releases page][releases].

Put it anywhere that you can run an executable from. It is completely portable to any directory, but
it is preferable you put it somewhere that is in your `PATH`.

Just run `gi_gen` in the directory you wish to add to and follow the prompts.

```shell
$ gi_gen
```

## Features

GI Gen does the following things:

- Discovers any gitignore templates that might be related to your project (optional)
- Optionally clean up results from any patterns that aren't in your project
- Writes to .gitignore (if a file already exists, you may append to it instead)

Credits to [open-source-ideas][osi] for the idea for the application.

[releases]: /releases/latest
[osi]: https://github.com/open-source-ideas/ideas/issues/296
