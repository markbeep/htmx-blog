---
title: "Mark's System Setups"
description: "The search for my best system setup is coming to an end with these tools"
date: 2026-02-19T09:10:50+01:00
showDate: true
tags: ["programming", "development"]
draft: true
---

Over the past years, during my studies, I've tried out a lot of tools and programs to find "the best setup" for myself. I've recently been asked a bit about my setup and thought some might be able to discover something new if I shared my setup.

Over time my mentality on what is important has changed a lot. I used to enjoy procrastinating studies by spending absurd amounts of time configuring neovim or my previous NixOS setup. But I now don't enjoy that anymore as it also distracts from the actual work I want to do (on my projects). I now only want things that *"just work"* with mostly the default settings. I want to be able to get to work and not have to configure random things all the time.

## Systems

I use two systems: A [Tuxedo](https://www.tuxedocomputers.com/) laptop and a full AMD CPU+GPU self-built desktop gaming PC.

On both systems I have 64GB ram. It's personally the minimum amount of RAM that I now get for my systems, not because I _need_ the RAM, but because it was a relative cheap (before the RAM price spike) upgrade you could do your PC and then it's a resource that you will never struggle with or be bottlenecked by. Granted I also don't do any ML.

With the desktop I am fully satisfied with. I haven't done a single bit of ML since I have it so I can't tell you how good cuda(-replacements) are nowadays. But as for daily use and gaming the only issues I've had which were due to it being AMD were that the GPU sometimes crashed for three specific games. Underclocking the GPU to 95% fixed any crashes.

The Tuxedo laptop is a bit of a mixed bag. Performance is great, but the touchpad is actually ass and sometimes starts becoming a bit unresponsive on fine movements (reloading the touchpad kernel module fixes it but I haven't been able to pinpoint where this issue comes from exactly and how to fix it. If you have an idea please hit me up.). The battery is not the best, but tbh I can just charge it wherever I am so not really something I mind.

## OS / Distro

Both on my laptop and desktop I daily-drive linux. On both I installed [EndeavourOS](https://endeavouros.com/). Here, I'll say it, EndeavourOS is the best linux distro out there if you just want a working system. On my desktop I used to run windows. I've now replaced it with EndeavourOS running KDE Plasma and it all works without having had to endlessly configuring anything. Everything you need is preinstalled and the feel is extremely familiar if you come from windows. Everything you expect to _just work_ from an OS just works; printers, drivers, usable settings UI, etc. I don't need any tiling window managers on my desktop and rather like to flick-shot my windows to the right places with my gaming mouse.

As a note, I mainly use my systems for programming and gaming, while most general work I do is in browsers (photopea, google docs, etc.) and therefore OS-agnostic. I don't work with excel/docx files locally and use google sheets and docs for those types of documents if required. Installing Steam and playing some games all just worked with no hiccups on my desktop.

EndeavourOS is arch-based, so you have the powerful AUR (package repository) at your fingertips, which is superior to homebrew or Nix (nix theoretically has more packages if you only look at the numbers, but if you've used Nix you'll know that it can often happen that some random library or tool simply doesn't exist on nix yet, while on arch you practically never face the problem of a package not existing.).

I'm a big fan of tiling window managers on laptops though since I don't like having to organize all my windows using the touchpad. Out of the box EndeavourOS can install Sway/i3, but I have a few gripes with the standard tiling window managers, mainly because I very dislike all my windows and browsers constantly resizing when a browser or something like Anki opens a new temporary window. But I also don't want to have to spend a lot of time configuring anything.

The solution is [DankMaterialShell](https://github.com/AvengeMedia/DankMaterialShell) (DMS); which is quite a whack name, but it's an insane project. People be glazing Omarchy while being completely unaware of the EndeavourOS+DMS combo. DMS provides the desktop feel with an easy-to-use settings-UI and customizable task/status bars with widgets. 

I therefore have been using [Niri](https://github.com/YaLTeR/niri) (easily installed with DMS) the last bunch of months as my tiling window manager. It is basically a tiling winow manager that have it stand out by providing infinite horizontal scrolling workspaces, with the mentality that opening any new window never results in an already open window from being resized.

## Package Management

Using NixOS for a few years got me used to keeping my system clear of random packages everywhere. I don't install development tools globally. I hate cluttering my global paths and system with random packages and libraries. I instead only install development tools specifically into the projects I need them in (more on that in the next section).

For globally installed packages, like editors or other general apps, I use paru (pacman-helper) to install them. I additionally often use paruz [1], a fuzzy package finder, to find the correct package if I'm unsure about the exact name of it.

Sometimes I install something just to quickly test it out, but I don't actually need it lying around. Or sometimes I install certain tools I only require for a course and would like to delete after the course is finished. That is where `aconfmgr` comes in. I already have a blogpost on it [here](/posts/tidy-system-with-aconfmgr). But tl;dr: It keeps track of what packages I install and forces my to organize my installed packages so I primarly know *why* I installed some package and *secondly* I can also quickly glance over certain categories to see what packages I could clean up and remove.

For example, this semester I had a `35-course-modules.sh` that included the specific programs two of my courses required:

```sh
# Applied Security Lab
AddPackage --foreign kathara # A lightweight container-based network emulation tool.
# Security of Wireless Networks
AddPackage gnuradio-companion # Signal processing runtime and signal processing software development toolkit (GUI)
```

At the end of the semester I could then simply delete these knowing I don't need them for anything else.

## Development Tools

The workflow of not installing development packages globally is something I picked up and been doing religiously ever since I started using NixOS. I have a [blogpost](/posts/nix-direnv-setup) on how I used to use `shell.nix` for my development tools. I've noticed that Nix is way too overkill for me and has the negatives of being too overcomplicated and pedantic for even the most minor tasks. I want to be developing my projects, not my devtool setup.

The main requirement is that I want to be able to work on multiple projects with the same tools, albeit with different versions. For example, I might have an older project that requires nodejs 18 to run but want to use nodejs 24 on a newer project simultanously. And the same goes for any type of development tool or language. I don't want to have to install tens of tools for managing the versions of all my devtools. So I need a way to set up project-specific environments and I want a single tool to manage all my tools.

The second requirement is that I want to be able to easily clean tools that I don't need anymore.

For quite a few months I've been trying to get devcontainers to work, and while devcontainers are quite a nice way to set up a clean dev environment, it brings with it a lot of negatives like permission problems, making it unnecessarily hard to use my shell and tool configurations, slows down startup of opening the project, as well as generally just not really being developed on a lot.

After enough annoyances with devcontainers, I finally decided to try out [mise](https://mise.jdx.dev/). I had a few doubts at first, but they were quickly swept away when I started using it and noticed how it can somehow just do everything I need it to. To start using mise, you install it and after entering a project directory, you can execute i.e. `mise use node@24`, which adds `node@24` to the local `mise.toml` configuration as well as also installing node itself with the right version. Extremely pain-free and fast. You can also install packages directly from cargo (with `mise use cargo:...`), npm, go, git, or even just install binaries directly. All while not being overly pedantic about EVERY. SINGLE. THING.

Mise is one of the best devtools I feel and a must-have if you tend to work on multiple projects with different tech-stacks often.

## Dotfiles

Managing dotfiles is something I still can see a lot of improvement in. At first I used a `--bare` git repository, but that quickly turns into a huge annoyance when I have to start diverging the branches because some things (aconfmgr config files for example) can be quite different on my laptop versus on my desktop.

Because of that, I instead use [dotter](https://github.com/SuperCuber/dotter) to manage my dotfiles. It has two main features I need:
1. By default any config files you manage are symlinked to whatever path you define. This then allows you to store your dotfiles repository wherever you want, while allowing you (or something like a settings app) to edit the original (`~/.config/...`) path as normal.
2. It allows for adding dotfiles conditionally as well as template any files I want. I use this to split what dotfiles are installed on my desktop and laptop. For example, I don't need any Niri or DankMaterialShell config on my desktop running KDE. Or since I use aconfmgr on both devices, I template the aconfmgr files to allow me to conditionally install certain programs like steam solely on my laptop while not installing it on my laptop.

I don't have my dotfiles repository public as of now though, since I have some secrets and slightly confidential files (ssh, kubectl, etc.) that I don't want to publish without first encrypting. I'm still on the lookout for a clean solution for that (or I simply keep my dotfiles private).

## Editor

I use Zed, LazyVim, VSCode (in order of usage) as my editors. Zed is currently my main editor I use for any project. It has some flaws. Like being quite uncustomizable, having bad UX surrounding extensions/LSPs, as well as sometimes having some visual bugs on my laptop. But it is extremely fast, has a minimalistic UI fully controllable with the keyboard whenever desired, git diff viewer, and a debugger (which I frankly haven't used *yet*).

Otherwise I use Neovim with [LazyVim](https://www.lazyvim.org/), again going along with the "no-setup" route since I really don't want to be spending my time configuring Neovim anymore. I once forced myself to only use Neovim for half a year for all my coding. I set everything up from scratch and had it be basically everything I wanted, but at the end of the day mouse support is still bad and it being in the terminal drastically limits it with any UI choices. I now solely use vim for all my quick terminal file edits and lazyvim is a nice way to quickly have a bunch of small must-haves and nice visuals.

The sole settings I've added to lazyvim is the `<leader>y` keymap to copy something into my system clipboard instead of the vim register. Same goes for paste and delete.

```lua
-- Only copy to clipboard on <leader>y/p/d
vim.opt.clipboard = ""
vim.keymap.set("n", "<leader>y", '"+y')
vim.keymap.set("v", "<leader>y", '"+y')
vim.keymap.set("n", "<leader>d", '"+d')
vim.keymap.set("v", "<leader>d", '"+d')
vim.keymap.set("n", "<leader>p", '"+p')
vim.keymap.set("v", "<leader>p", '"+p')
```

VSCode is currently my backup editor for when I need some specific feature I can't seem to find on Zed or generally when I need an editor I know just works.

One thing with the mise workflow is that for the environment to be properly accepted by the three editors, you need to first enter the specific project directory and activate the mise environment, before then launching the editor of your choice. I don't open my editor and afterwards proceed to open the directory from there, since that does not set up the environment correctly. But I also don't want to have to tediously always have to open a terminal, navigate to the project and run `nvim .`/`code .`/`zeditor .`. For that I have a popup that launches with Super+P. I can then select the project directory I want, it `cd`s into the directory, executes `mise x -- zeditor .` and lets me correctly use the the environment in my editor [2]. More details on the popup thing can be found in my [popup blogpost](/posts/popup-tools). I've configured the keybind and popup to work on both my KDE and Niri systems.

## Git

I don't use the standard `git` CLI anymore and have been exclusively using [jj](https://docs.jj-vcs.dev/) (jujutsu) the past weeks. I blame Lu (in a positive way) for getting me hooked on yet another great tool. [Steve's Jujutsu Tutorial](https://steveklabnik.github.io/jujutsu-tutorial/introduction/introduction.html) is a great place to learn how to use jj. The main thing about jj is that it's a different workflow and way of thinking about git, while being fully compatible with the git backend. This allows you to easily use jj in place of git without everyone else in a project also having to immediately swap. jj makes jumping between commmits and branches extremely easy and never requires you to ever explicitly git stash. It also makes rebasing a breeze and if you ever mess something up can look at a long log of all operations you executed and revert any change you want. If you ever had the thought of git not being the perfect tool, then you should definitely check out jj.

## Footnotes

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
    # Shows all directories in ~/Documents, enters directory and loads
    # the existing environment.
    # - Initial fzf is sorted most recent edited. Most recent at the bottom.
    # - sed to shorten the path
    # - awk colors the final directory blue and bold
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
