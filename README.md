<h1>GI Gen - Gitignore File Generator</h1>

<details>
<summary>Table of Contents</summary>

- [Usage](#usage)
- [Features](#features)
- [Known Bugs](#known-bugs)
- [Contribute](#contribute)

</details>

<hr />

GI Gen is an open source CLI to generate `.gitignore` files. It is completely cross-platform, and
standalone (no dependencies - other than `git` itself), so you may literally use it for any project
on any platform.

You can run this CLI program to create or append a `.gitignore` file from a chosen list of template
easily.

You may choose more than one template to generate.

## Usage

Download the file for your platform in the [Releases page][releases].

Put it anywhere that you can run an executable from. It is completely portable to any directory, but
it is preferable you put it somewhere that is in your `PATH`.

Just run `gi_gen` in the directory you wish to add to and follow the prompts.

```shell
$ gi_gen
```

## Features

GI Gen supports the following features:

- `.gitignore` discovery:
  - Auto-discover any gitignore templates that might be related to your project
  - Optionally list all available templates instead (see [github/gitignore][gh-gi] for the complete
    list of templates)
- `.gitignore` clean: Clean up results from any patterns that aren't in your project before
  outputting (optional)
- Writes to `.gitignore` file in current directory (you may overwrite/skip/append if already exists)

Credits to [open-source-ideas][osi] for the idea for the tool.

## Known Bugs

If the binary doesn't work for you, try installing via `go install`, if you have Go lang tools
installed.

```shell
go install github.com/chenasraf/gi_gen
```

Or clone this repository and install directly from source.

## Contribute

Please feel free to open PRs or issues with bug fixes/reports, or feature requests.

This project was built using Go, and should run easily with the normal Go tools with no further
configuration.

Testing was only done on macOS using ARM architecture, so feel free to report any issues on your
platform if you have any, or are missing your platform and cannot/don't want to build from source (I
tried building for the most common platforms).

If you are feeling incredibly generous and appreciate the time &amp; effort I put into developing
this tool, kindly consider donating any amount to help me make up for the work hours. It is really
very much appreciated! 🙏🏼

<a href='https://ko-fi.com/casraf' target='_blank'>
  <img height='36' style='border:0px;height:36px;'
    src='https://cdn.ko-fi.com/cdn/kofi1.png?v=3'
    alt='Buy Me a Coffee at ko-fi.com' />
</a>

[releases]: https://github.com/chenasraf/gi_gen/releases/latest
[osi]: https://github.com/open-source-ideas/ideas/issues/296
[gh-gi]: https://github.com/github/gitignore
