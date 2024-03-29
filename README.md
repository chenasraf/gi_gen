<h1>GI Gen - Gitignore File Generator</h1>

<details>
<summary>Table of Contents</summary>

- [Features](#features)
- [Command Line Usage](#command-line-usage)
  - [Command Line Flags](#command-line-flags)
  - [Examples](#examples)
- [Contribute](#contribute)

</details>

<hr />

GI Gen is an open source CLI to generate `.gitignore` files for most project types.

Simply run the command and follow the prompts, and your project will have a `.gitignore` file to
match it.

It is completely cross-platform, and standalone (no dependencies - other than `git` itself), so you
may literally use it for **any project** on **any platform**.

You can run this CLI program to create or append a `.gitignore` file from a chosen list of template
easily.

You may choose more than one template to generate.

**Check out the [VS Code Extension](https://github.com/chenasraf/gi-gen-vscode-extension)!**

## Features

GI Gen supports the following features:

- `.gitignore` discovery:
  - Auto-discover any gitignore templates that might be related to your project
    - Can confidently discover over
      [50 project languages](https://github.com/chenasraf/gi_gen/issues/2) using your project
      structure
    - Can fall back on process of elimination using patterns in the template
  - Optionally list all available templates instead (see [github/gitignore][gh-gi] for the complete
    list of templates)
- `.gitignore` clean: Clean up results from any patterns that aren't in your project before
  outputting (optional)
- Writes to `.gitignore` file in current directory (you may overwrite/skip/append if already exists)

## Command Line Usage

Download the file for your platform in the [Releases page][releases].

Put it anywhere that you can run an executable from. It is completely portable to any directory, but
it is preferable you put it somewhere that is in your `PATH`.

Just run `gi_gen` in the directory you wish to add to and follow the prompts.

```shell
$ gi_gen
```

### Command Line Flags

You may pass additional flags to `gi_gen`. These are the currently available flags:

| Usage                    | Description                                                                                                                                                      |
| ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `-languages` \| `-l`     | List the languages you want to use as templates.<br />To add multiple templates, use commas as separators, e.g.: `-languages Node,Python`                        |
| `-auto-discover` \| `-d` | Use auto-discovery for project, detecting the project type and using the result as the pre-selected template list.                                               |
| `-clean-output` \| `-c`  | Perform cleanup on the output .gitignore file, removing any unused patterns                                                                                      |
| `-keep-output` \| `-k`   | Do not perform cleanup on the output .gitignore file, keep all the original contents                                                                             |
| `-append` \| `-a`        | Append to .gitignore file if it already exists                                                                                                                   |
| `-overwrite` \| `-w`     | Overwrite .gitignore file if it already exists                                                                                                                   |
| `-detect-languages`      | Outputs the automatically-detected languages, separated by newlines, and exits. Useful for outside tools detection.                                              |
| `-all-languages`         | Outputs all the available languages, separated by newlines, and exits. Useful for outside tools detection.                                                       |
| `-clear-cache`           | Clear the .gitignore cache directory, for troubleshooting or for removing trace files of this program.<br />Exits after running, so other flags will be ignored. |
| `-help` \| `-h`          | Display help message                                                                                                                                             |

### Examples

- Run normally, prompting for every step:

  ```shell
  gi_gen
  ```

- Pre-select languages (skip prompt):

  ```shell
  gi_gen -languages Node          # One language
  gi_gen -languages Node,Python   # Multiple languages
  ```

- Perform clean up (skip prompt):

  ```shell
  gi_gen -clean-output    # clean up
  gi_gen -keep-output     # skip clean up
  ```

- Use auto-discovery (skip prompt):

  ```shell
  gi_gen -auto-discover
  ```

- Existing file handlers (skip prompt):

  ```shell
  gi_gen -append      # if file exists, add to end of it
  gi_gen -overwrite   # if file exists, replace the existing content
  ```

- Combined (skip all prompts):

  ```shell
  gi_gen -d -c -a       # Use auto discovery, append and clean unused patterns
  gi_gen -l Node -k -w  # Use passed languages, overwrite and keep unused patterns
  ```

- Clean cache directory and exit:

  ```shell
  gi_gen -clear-cache
  ```

- Detect languages and output the results, then exit:

  ```shell
  gi_gen -detect-languages
  ```

## Contribute

Credits to [open-source-ideas][osi] for the idea for the tool.

Please feel free to open PRs or issues with bug fixes/reports, or feature requests.

This project was built using Go, and should run easily with the normal Go tools with no further
configuration.

Tested on all major platforms, but feel free to report any issues on your platform if you have any,
or are missing your platform and cannot/don't want to build from source (I tried building for the
most common platforms).

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
