---
title: "Mark's System Setups"
description: "How I have my systems optimized for programming"
date: 2026-02-19T09:10:50+01:00
showDate: true
tags: ["programming", "development"]
draft: true
---

Over the past years, during my studies, I've tried out a lot of tools and programs to find "the best setup" for myself. I've recently been asked about my setup a bunch and thought why not just write another blogpost about it I can direct people to.

Over time my mentality on what is important in a system has changed a lot. I used to enjoy procrastinating studies by spending absurd amounts of time configuring Neovim or my previous NixOS setup. But I now don't enjoy that anymore as it also distracts from the actual work I want to do (on my projects). I now only want things that _"just work"_ with mostly the default settings. I want to be able to get to work and not have to configure random things all the time. But I enjoy having the option to still customize my system however I want without having the OS restrict me.

# Table of Contents

- [Table of Contents](#table-of-contents)
- [System Setup](#system-setup)
  - [Systems](#systems)
  - [OS / Distro](#os--distro)
  - [Package Management](#package-management)
  - [Development Tools](#development-tools)
  - [Dotfiles](#dotfiles)
  - [Editor](#editor)
  - [Version Control](#version-control)
- [Conclusion](#conclusion)
- [Footnotes](#footnotes)

# System Setup

## Systems

I use two systems: A [Tuxedo](https://www.tuxedocomputers.com/) laptop and a full AMD CPU+GPU self-built desktop gaming PC.

On both systems I have 64GB ram. It's personally the minimum amount of RAM that I now get for my systems, not because I _need_ the RAM, but because it was a relatively cheap (before the RAM price spike) upgrade you could do to your PC and then it's a resource that you will never struggle with or be bottlenecked by. Granted I also don't do any ML.

I am fully satisfied with my desktop. I haven't done a single bit of ML since I have it so I can't tell you how good cuda(-replacements) are nowadays. But as for daily use and gaming the only issues I've had which were due to it being AMD were that the GPU sometimes crashed for three specific games. Underclocking the GPU to 95% fixed any crashes.

The Tuxedo laptop is a bit of a mixed bag. Performance is great, but the touchpad is actually ass and sometimes starts becoming a bit unresponsive on fine movements (reloading the touchpad kernel module fixes it but I haven't been able to pinpoint where this issue comes from exactly and how to fix it. If you have an idea please hit me up.). The battery is not the best, but to be honest I can just charge it wherever I am so it's not really something I mind a lot. I don't do any very long trips without a power bank on hand or a plug available in the train.

## OS / Distro

Both on my laptop and desktop I daily-drive Linux. On both I installed [EndeavourOS](https://endeavouros.com/). Here, I'll say it, EndeavourOS is the best Linux distro out there if you just want a working system without straight up bloating your system. On my desktop I used to run Windows. I've now replaced it with EndeavourOS running KDE Plasma and it all works without having to endlessly configure anything. Everything you need is preinstalled and the feel is extremely familiar if you come from Windows. Everything you expect to _just work_ from an OS just works; printers, drivers, usable settings UI, etc. The install processes is a straightforward UI where you easily click yourself through.

As a note, I mainly use my systems for programming and gaming, while most general work I do is in browsers (photopea, google docs, etc.) and therefore OS-agnostic. I don't really work with excel/docx files locally (other than maybe having to read one once in a blue moon) and mainly use google sheets and docs if required. Installing Steam and playing some games all just worked with no hiccups on my desktop.

EndeavourOS is Arch-based, so you have the powerful AUR (package repository) at your fingertips, which is superior to homebrew or Nix (Nix theoretically has more packages if you only look at the numbers, but if you've used Nix you'll know that it can often happen that some random library or tool simply doesn't exist on Nix yet, while on Arch you practically never face the problem of a package not existing.).

As for my laptop, I'm a big advocate for tiling window managers on laptops. On my desktop I would never install a tiling window manager. I have two monitors and a gaming mouse that I use to flick-shot all my windows to the correct positions. On a laptop, the screen real estate is a lot more limiting though and I also only ever use my laptop with the touchpad as I'm not going to carry around a mouse everywhere. There I want a system that saves me the work of having to tediously move around windows with the touchpad. The solution? A tiling window manager!

Out of the box EndeavourOS can install _standard_ window managers like Sway/i3. By "standard" I mean the type of tiling window managers that simply split your screen and allows you to move around and resize them on the screen. My main gripe with such window managers though, is that I strongly dislike all my windows constantly resizing when a browser or something like Anki opens a new temporary window. If I have a browser open and am on some website that cares too much about the size of my browser, the moment I try to add a new card to Anki (which opens a new window), the all my windows, including the browser, are then squished and messed up until I fix up the layout of my screen.

The solution to that is a different type of tiling window manager. The one I use is called [Niri](https://github.com/YaLTeR/niri). It is a tiling window manager that provides infinite horizontally-scrolling workspaces.
Opening a new window will never resize any of the other windows on the page. Instead, it simply shifts over the workspace as if you are on an infinitely long strip of windows.

So I instead also installed EndeavourOS with KDE on my laptop just to have a working system for the start. But I then used [DankMaterialShell](https://github.com/AvengeMedia/DankMaterialShell) (DMS) to install both Niri and pre-configured [Quickshell](https://quickshell.org/) to have all the required task/statusbars out of the box with fancy widgets.
It's quite a whack name, but it's an insane project and is quite under the radar. People be glazing Omarchy while being completely unaware of the EndeavourOS+DMS combo. DMS provides the desktop feel with an easy-to-use settings-UI while supporting Niri out of the box.

## Package Management

Using NixOS for a few years got me used to keeping my system clear of random packages everywhere. I don't install development tools globally. I hate cluttering my global paths and system with random packages and libraries. I instead only install development tools specifically into the projects I need them in (more on that in the next section).

For globally installed packages, like editors or other general apps, I use paru (pacman-helper) to install them. I additionally often use paruz [1], a fuzzy package finder, to find the correct package if I'm unsure about the exact name of it.

Sometimes I install something just to quickly test it out, but I don't actually need it lying around. Or sometimes I install certain tools I only require for a course and would like to delete after the course is finished. That is where `aconfmgr` comes in. I already have a blogpost on it [here](https://www.google.com/search?q=/posts/tidy-system-with-aconfmgr). But tl;dr: It keeps track of what packages I install and forces me to organize my installed packages so I know _why_ I installed some package and I can also quickly glance over certain categories to see what packages I could clean up and remove.

For example, this semester I had a `35-course-modules.sh` that included the specific programs two of my courses required:

```sh
# Applied Security Lab
AddPackage --foreign kathara # A lightweight container-based network emulation tool.
# Security of Wireless Networks
AddPackage gnuradio-companion # Signal processing runtime and signal processing software development toolkit (GUI)
```

At the end of the semester I could then simply delete these knowing I don't need them for anything else.

## Development Tools

The workflow of not installing development packages globally is something I picked up and been doing religiously ever since I started using NixOS. I have a [blogpost](https://www.google.com/search?q=/posts/nix-direnv-setup) on how I used to use `shell.nix` for my development tools. I've noticed that Nix is way too overkill for me and has the negatives of being too overcomplicated and pedantic for even the most minor tasks. I want to be developing my projects, not my devtool setup.

The main requirement is that I want to be able to work on multiple projects with maybe overlapping tools (albeit with different versions) and I don't want a change in one project to affect the tooling in another project. For example, I might have an older project that requires nodejs 18 to run but want to use nodejs 24 on a newer project simultaneously. For node, I could use [nvm](https://github.com/nvm-sh/nvm) or [n](https://github.com/tj/n) to select the correct version manually every time I work on a specific project. But I need the same for any type of development tool or language. I don't want to have to install tens of tools for managing the versions of all my devtools. And I also don't want to constantly have to remember to switch around the tool versions when I switch over to another project. So I need a way to set up project-specific environments and I want a single tool to manage all my tools.

A third requirement is that I want to be able to easily clean tools that I don't need anymore.

For quite a few months I've been trying to get devcontainers to work, and while devcontainers are quite a nice way to set up a clean dev environment, it brings with it a lot of negatives like permission problems, making it unnecessarily hard to use my shell and tool configurations, slows down startup of opening the project, as well as actual feature development being quite stale.

After enough annoyances with devcontainers, I finally decided to try out [mise](https://mise.jdx.dev/). I had a few doubts at first, but they were quickly swept away when I started using it and noticed how it can somehow just do everything I need it to. To start using mise, you install it and after entering a project directory, you can execute e.g. `mise use node@24`, which adds `node@24` to the local `mise.toml` configuration as well as also installing node itself with the right version into the current shell. Extremely pain-free and fast. You can also install packages directly from cargo (with `mise use cargo:...`), npm, go, git, or even just install binaries directly. All while not being overly pedantic about EVERY. SINGLE. THING.

If you have mise set up in your shell config with `mise activate`, it will automatically download and add tools your shell environment when you enter a project directory, while removing them again when you leave the directory. My workflow when setting up a new project with mise is to simply enter the directory, type `mise use zig@`, tap tab twice for autocomplete to kick in, up arrow to jump to the most up to date version, and then I simply select whatever new version I want. For languages, I'll fixate the version (to at least the major version), while for less important tools like [just](https://just.systems/) I'll use `@latest`.

![](/content/posts/my-tools/mise-use-zig.png)

PS: Mise also allows for setting environment variables either one-by-one or by directly giving the path to an environment variable file. Extremely nice when you're working on a project that has a local docker-compose setup with docker-network specific environment variables, but you also want to be able to develop locally outside of the docker network, meaning you would need to change some environment variables like `POSTGRES_HOST` from `postgres` to `localhost`. Mise allows you to read a docker environment file (let's call it `.env.docker`), but then simply override a few of those variables without any hassle:

```sh
[env]
_.file = ".env.docker"
POSTGRES_HOST = "localhost"
```

If you're interested, you can view this exact example more concretely on the [VIS Community Solutions repo](https://gitlab.ethz.ch/vseth/sip-com-apps/community-solutions/-/blob/master/mise.toml).

To summarize, I personally find mise one of the best devtools and a must-have if you tend to work on multiple projects with different tech-stacks.

## Dotfiles

Managing dotfiles is something I still can see a lot of improvement in. At first I used a `--bare` git repository. It is very simple to set up and use. But it quickly becomes messy, since I have some dotfiles (like aconfmgr config files) that are vastly different on my KDE desktop and my DMS/Niri laptop. That forced me to maintain two separate branches of my dotfiles, which then makes it extremely annoying to share dotfiles between the two systems.

Because of that, I instead use [dotter](https://github.com/SuperCuber/dotter) to manage my dotfiles. It has two main features I need:

1. By default any config files you manage are symlinked to whatever path you define. This then allows you to store your dotfiles repository wherever you want, while allowing you (or something like a settings app) to edit the original (`~/.config/...`) path as normal.
2. It allows for adding dotfiles conditionally as well as template any files you define. I use this to split what dotfiles are installed on my desktop and laptop. For example, I don't need any Niri or DankMaterialShell config on my desktop running KDE. Or since I use aconfmgr on both devices, I template the aconfmgr files to allow me to conditionally install certain programs, like Steam, solely on my desktop while not installing it on my laptop.

I don't have my dotfiles repository public as of now though, since I have some secrets and slightly confidential files (ssh, kubectl, etc.) that I don't want to publish without first encrypting. I'm still on the lookout for a clean solution for that and open for any suggestions. Or another alternative is to simply keep the repository private as it is now.

## Editor

I use Zed, LazyVim, VSCode (in order of usage) as my editors. Zed is currently my main editor I use for any project. It has some flaws. Like being quite uncustomizable, having bad UX surrounding extensions/LSPs, as well as sometimes having some visual bugs on my laptop. But it is extremely fast, has a minimalistic UI fully controllable with the keyboard whenever desired, git diff viewer, and a debugger (which I frankly haven't used _yet_). I also use it mostly with all default settings, other than a few keybinds like switching tabs or opening/closing the terminal [3].

For smaller things and anything in the terminal, I use Neovim with [LazyVim](https://www.lazyvim.org/). Again going along with the "no-setup" route since I really don't want to be spending my time configuring Neovim anymore. I once forced myself to only use Neovim for half a year for all my coding. I set everything up from scratch and had it be basically everything I wanted, but at the end of the day mouse support is still bad and it being in the terminal drastically restricts it with any UI/UX choices. I now solely use vim for all my quick terminal file edits and lazyvim is a nice way to quickly have a bunch of small must-haves and nice visuals.

The sole settings I've added to lazyvim is the `<leader>y` keymap to copy something into my system clipboard instead of the vim register. Same goes for paste and delete.

```lua
vim.opt.clipboard = ""
vim.keymap.set("n", "<leader>y", '"+y')
vim.keymap.set("v", "<leader>y", '"+y')
vim.keymap.set("n", "<leader>d", '"+d')
vim.keymap.set("v", "<leader>d", '"+d')
vim.keymap.set("n", "<leader>p", '"+p')
vim.keymap.set("v", "<leader>p", '"+p')
```

VSCode is currently my backup editor for when I need some specific feature I can't seem to find on Zed or generally when I need an editor I know just works. I'm a fan of how powerful VSCode extensions are, but when working with a lot of different tools and languages, each requiring a different set of extensions, it quickly slows down VSCode itself. You can choose to only install extensions on a per-project basis, but that feature is still quite clunky and quite an afterthought. I wish I could have an almost empty VSCode with only a handful of global extensions, and then define in a file in the project directory what extensions should be used/installed when working in this directory. Devcontainers allow for that, but as mentioned previously that brings with itself all other sorts of problems. That is why I'm fond of Zed right now, because while the extensions themselves are somewhat limiting, they're written in rust and don't slow down the startup and editor so much.

One thing with the mise workflow is that for the environment to be properly accepted by the three editors, you need to first enter the specific project directory and activate the mise environment, before then launching the editor of your choice. I don't open my editor and afterwards proceed to open the directory from there, since that does not set up the environment correctly. But I also don't want to have to tediously always have to open a terminal, navigate to the project and run `nvim .`/`code .`/`zeditor .`. For that I have a popup that launches with Super+P. I can then select the project directory I want, it `cd`s into the directory, executes `mise x -- zeditor .` and lets me correctly use the the environment in my editor [2]. More details on the popup thing can be found in my [popup blogpost](https://www.google.com/search?q=/posts/popup-tools). I've configured the keybind and popup to work on both my KDE and Niri systems.

## Version Control

I don't use the standard git CLI anymore and have been exclusively using [jj](https://docs.jj-vcs.dev/) (jujutsu) the past weeks. [Steve's Jujutsu Tutorial](https://steveklabnik.github.io/jujutsu-tutorial/introduction/introduction.html) is a great place to learn how to use it. The main thing about jj is that it's a different workflow and way of thinking about git, while being fully compatible with the git backend. This allows you to easily use jj in place of git without everyone else in a project also having to immediately swap. jj makes jumping between commits and branches extremely easy and never requires you to ever explicitly git stash. It also makes rebasing a breeze and if you ever mess something up, you can look at a long log of all operations you executed and revert any change you want. I've always thought that git isn't the perfect version control system and there's room for improvement. jj also isn't perfect, since it also uses the git backend after all, but it is a nice upgrade. jj can handle the same operations (if not more) as git, while being simpler to use. And in the case you need a git-specific feature, like adding tags (doesn't make sense in the jj sense), you can still always just use specific git commands interchangeably with jj.

# Conclusion

The post turned out a bit wordier than I'd hoped. But I hope I could maybe motivate you to try out some new tool or set up your system to allow for a more seamless workflow.

# Footnotes

[1] Paruz used to be available on the [AUR](https://aur.archlinux.org/packages/paruz-git). But the repository has been removed on Github and the version still available on AUR uses a broken paru version. I instead use the following fish function which allows for the same fuzzy matching of packages, but simply keep the name because it's fitting:

```sh
function paruz
    set -l pkgs (
      paru -Sl |\
      sed -e "s: :/:; s/ unknown-version//; /installed/d" |\
      fzf --multi --ansi --preview "paru -Si {1}"
    )
    if test -z "$pkgs"
        return
    end
    echo $pkgs | awk '{print $1}' | xargs -ro paru -S
end

```

[2] The following fish function is what allows me to fuzzy match the directory I want to enter (using fzf) and then instantly open the editor in it with the mise environment tools loaded (if existing, the editor still opens even if no `mise.toml` is located in the given directory).

```sh
function zed_open
    set -l directory (
      fd . --type d ~/Documents |\
        xargs ls -td |\
        sed 's|/$||' |\
        sed 's|^/home/mark/Documents/|~/D/|' |\
        awk -F/ '{OFS="/"; last=$NF; $NF="\033[1;34m"last"\033[0m"; print}' |\
        fzf --ansi --prompt="Zed Recent> "
    )
    if test -n "$directory"
        set -l clean_dir (string replace -ra '\e\[[0-9;]*m' '' -- "$directory")
        set -l formatted (echo "$clean_dir" | sed 's|^~/D/|/home/mark/Documents/|')
        sh -c "cd '$formatted' && mise x -- zeditor ."
    end
end

```

[3] Because of browsers I'm used to `Ctrl+Tab` moving across tabs in the order they're open. Editors like VSC/Zed default to `Ctrl+Tab` switching tabs by recency which I find quite annoying/unintuitive. In Zed I now use these keybinds:

```json
{
  "context": "Pane",
  "bindings": {
    "ctrl-tab": "pane::ActivateNextItem",
  },
},
{
  "context": "Pane",
  "bindings": {
    "ctrl-shift-tab": "pane::ActivatePreviousItem",
  },
},
{
  "context": "Workspace",
  "bindings": {
    "ctrl-tab": null,
  },
},
{
  "context": "Workspace",
  "bindings": {
    "ctrl-shift-tab": null,
  },
}
```
