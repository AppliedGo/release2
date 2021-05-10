package main

/*
<!--
Copyright (c) 2019 Christoph Berger. Some rights reserved.

Use of the text in this file is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

Use of the code in this file is governed by a BSD 3-clause license that can be found
in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Packaging a project release (goreleaser part 2)"
description = "Create OS install packages from your Go release with goreleaser"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2021-05-08"
draft = "false"
categories = ["Go Ecosystem"]
tags = ["deployment", "packaging", "distribution"]
articletypes = ["Tutorial"]
+++

In the previous post, I used goreleaser to add binaries to a project release. Now let's have goreleaser build a Homebrew formula as well. Automatically, and for macOS and Linux alike.

<!--more-->

## Convenience comes in packages

There was a tiny problem with the macOS binary in the [previous post]({{< ref "/release" >}} . As it is not signed by an Apple developer, it get quarantined when downloaded through a browser. So this time let's package the binary for delivery via Homebrew. Using a package manager is more convenient for the users anyway, and they also get updates delivered and can easily uninstall the package if they don't want to use it anymore.


## Homebrew

![Beer mugs with apple and penguin](homebrew-linuxbrew.png)

*(Mug images licensed under [CC-BY](https://creativecommons.org/licenses/by/4.0/) the [Homebrew project](https://github.com/Homebrew/brew) owners.)*

macOS users, at least the tech-savvy ones, know and love [Homebrew](https://brew.sh/), the Missing Package Manager. Moreover, Homebrew is also [available for Linux](https://docs.brew.sh/Homebrew-on-Linux). And I mention that here because gorleeaser can build Homebrew formulae for both macOS and Linux.

Homebrew has a core repository of packages but also supports providing your own "[tap](https://docs.brew.sh/Taps)" for dispensing your freshly brewed binary, which is what I will  be using here. (If you really want your project listed in the Homebrew core repository, you can do so but the [process](https://docs.brew.sh/Formula-Cookbook) is quite involved, and `goreleaser` also does not generate valid `homebrew-core` formulas at the moment.)

### Let's start brewing

To create a formula, I add a `brews` section to `.goreleaser.yml`. Note the plural – it is possible to generate multiple formulae.

I am going to create the most minimal config possible. (After all, I don't want to replicate the manual but only demonstrate you what is possible and how to get up and running in the easiest way.) That means I rely on default values wherever possible. For example, the recipe name template defaults to the project name, which should be fine in most cases.

You can add customization where necessary, as listed in the [Homebrew section ](https://goreleaser.com/customization/homebrew/)of the `goreleaser` docs.

The first necessary entries are the name and owner of the GitHub or GitLab repository that shall serve as your Homebrew tap. Yes, a tap is nothing else but a GitHub or GitLab repo. A Homebrew user just needs to tell their Homebrew client the name of your tap, and then they can immediately install and upgrade formulae from that tap.

```yaml
brews:
  -
    tap:
      owner: appliedgocode
      name: homebrew-tools
```

`goreleaser` will now publish my formulae to a repository named "appliedgocode/homebrew-tools". The name of the tap itself is the same minus the "homebrew-" part: "appliedgocode/tools".

If I wanted to have a tap for `goman` only, I could have named the repo "homebrew-goman". However, as I plan to provide more tools from the same tap, I chose the quite generic name "tools" here.

Also, I want to avoid unnecessary repetition: \
`    brew install appliedgocode/tools/goman`  \
sounds better than \
`    brew install appliedgocode/goman/goman`.

But that's just a matter of taste. I have come across quite a few such taps with the same name for tap and tool, and that's ok. We're not in a style contest here, are we?

The next entry without useful default values is the Git author who commits to the repository. Add your GitHub/GitLab username and email here. I added mine:

```yaml
    commit_author:
      name: christophberger
      email: my@github.email
```

Together with the GITHUB_TOKEN that I already created in the previous post, `goreleaser` now can publish the formula to the tap repo under my name.

The following values are not strictly necessary but I would highly recommend setting them. Especially if your tool has some caveats, ensure to list them in the caveats entry. Homebrew displays this entry after installing the tool.

For the license name, ensure to use the [SPDX identifier](https://spdx.org/licenses/) of your license.

```yaml
    homepage: "https://github.com/appliedgocode/goman"
    description: "The missing man pages for Go binaries"
    license: "BSD-3-Clause"
    caveats: "Returns strange results at full moon"
```

The following entry controls when to upload (commit) the formula to the Homebrew tap. Default is `false`, which means the formula is always uploaded to the tap. Use `true` to always skip the upload, and `auto` for skipping the upload if the latest Git tag contains an indicator for a pre-release (like, for example, v1.1.0-beta1). I am using `auto` here, as I don't want to pester all `goman` users with pre-release updates unnecessarily.

```yaml
    skip_upload: auto
```

If your tool has dependencies on other tools that are available via Homebrew, you can add a `dependencies` subsection. Homebrew will then install the dependencies as well. And if there are any Homebrew packages that are in conflict with yours, ensure to list them in `conflicts`. Homebrew warns the user if any of these conflicting packages is installed already.

```yaml
    # just examples, not required for goman:
    dependencies:
        - name: git
        - name: zsh
            type: optional

        conflicts:
        - svn
        - bash
```

As `goman` has neither any dependencies nor any known conflicts, I omit both.

So that's my complete but fairly minimal `brews` section in `.goreleaser.yml`:

```yaml
brews:
  - tap:
      owner: appliedgocode
      name: tap
    commit_author:
      name: christophberger
      email: my@github.email
    homepage: "https://github.com/appliedgocode/goman"
    description: "The missing man pages for Go binaries"
    license: "BSD-3-Clause"
    skip_upload: auto
```

(Yes, the caveat is not included. It does not apply to `goman`, or so I hope!)

### And now: release!

To prepare the release, I only need two last steps.

1. I create my new Homebrew tap at github.com/appliedgocode/homebrew-tools.
2. I create and push a new Git tag

Now here comes the magic: I just need to push a new Git tag and run

```sh
goreleaser release --rm-dist
```

and `goreleaser` builds and uploads a `goman` formula to my tap repository. Now it is ready for installing!

Soooo convenient.

I now just need to switch my brains from "tool publisher" to "Homebrew user" and install `goman`.

```sh
brew tap appliedgocode/tools
brew install goman
```

Worked! Awesome!

With Homebrew, your releases already cover macOS, Linux, and even Windows to some extent, via Windows Subsystem for Linux (WSL).

But you can get even more specific and create pure Linux packages (.deb, .rpm, .apk, or a Snapcraft snap), pure Windows packages (for Scoop), and even a Docker image.  And more. I leave this as a homework exercise.


**Happy coding!**

*/
