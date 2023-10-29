---
title: "Setting up project specific tools using Nix"
date: 2023-06-09T09:47:41+02:00
showDate: true
tags: ["programming", "nix", "direnv"]
---

Over my past blog posts, I've made it clear that I enjoy how much I like
setting up all my projects in the declarative way of Nix. In this article, I
will go more in detail on how exactly I have `direnv` set up to automatically
load my tools when I hop into a project.

In my [last article](/posts/nix-direnv-setup) I talked about how I don't have
most of my programs installed globally, but instead, I have it set up so that
whenever I jump into one of my project directories it automatically loads the
correct packages and tools into my terminal environment. Since I discovered
`direnv` 4 months ago I've been using it for all my projects.

# The Idea

The basic idea of how I have everything set up is that globally I barely have
anything installed. Only tools I use a lot or don't want to install anew every time I wanna use them. This includes applications like Firefox,
Neovim, Kitty, etc. Everything else related to programming or tools I only need
for a one-time use I install differently.

Tools that I only need for one-time use, say I want to quickly run `dig` to
figure out what IP a domain points to or I quickly try something out in GIMP, I
open up my terminal and run `nix-shell -p dig` (or `gimp`). Then I run the use tool and
when I'm done I'll simply close the terminal and let the Nix garbage collector
clean up anything I installed. The garbage collector runs automatically on a daily basis.

# The Execution

To get everything to work I use [lorri](https://github.com/nix-community/lorri). It extends [direnv](https://direnv.net/) to add some features specific to Nix.

## direnv

[direnv](https://direnv.net/) is a tool that allows you to automatically load and unload environment variables depending on the directory you're in.
Simply add `direnv` to your configuration.nix, home manager or nix-env packages:

On NixOS:

```
environment.systemPackages = [
    pkgs.direnv
];
```

On non-NixOS:

```bash
nix-env -iA nixpkgs.direnv
```

Now enter a directory, create a file called `.envrc` and add some export.

```bash
# .envrc
export HELLO=123
```

direnv initially blocks a `.envrc` file upon its creation or upon being changed. You first need to allow it:

```bash
direnv allow
```

direnv now automatically loads the export (which can be confirmed with `echo $HELLO`). How cool is that?
When you exit the directory the export will automatically be unloaded again.

## shell.nix

It can quickly become tedious to always install tools using the above `nix shell` if need to install a lot of things.
To make it easier you can create a `shell.nix` file and specify what packages should be installed and what commands should be executed in the shell.

A barebones `shell.nix` is as follows:

```nix
{ pkgs ? import <nixpkgs> {} }:
with pkgs;
mkShell {
  buildInputs = [
    python311
    poetry
    postgresql_15
  ];
}
```

To get into this shell simply execute `nix-shell`. If the file is in another directory you can also specify the file with the `-f` flag. This will install all the listed programs and set you into a new environment.

But we want this automated; by adding `use nix` to the `.envrc` file direnv will now automatically load this `shell.nix` file whenever you enter the directory, automatically giving you access to Python 3.11, Poetry and Postgresql. When you exit the directory they will all be unloaded from your shell environment.

## lorri

Now you might notice that jumping into the directory can start to get really slow once your `shell.nix` file grows to a lot of packages. Especially after you run the Nix garbage collector, because then everything has to be downloaded anew.

That is what [lorri](https://github.com/nix-community/lorri) solves. Lorri adds two game-changing changes:

1. shell.nix environments are cached. That makes loading into big projects very quick. You don't even notice that something is happening anymore (ignoring the big message appearing in the terminal).
2. All installed packages are marked to not be garbage collected. Now you can garbage collect without having to worry about deleting all the packages you need for a project.

This is all handled completely automatically. All you need to do is initialize your `.envrc` file by running `lorri init`, so that lorri is correctly ran upon entering the directory.

Additionally, when you delete a project or shell.nix, remember to also run `lorri gc` to remove all the links so that the Nix garbage collector will delete the unused packages.

# The Truth

What motivated me to write this blog post was that I actually started to finally get into what Nix flakes are and how they work. Explained very roughly, they are a new thing from Nix that upgrades the reproducibility for everything. It takes inspiration from other project-based package managers like npm or [Poetry](https://python-poetry.org/) which generate a lock file to make everything you do extremely reproducible.
Have a look at this [introduction to flakes](https://woile.dev/posts/nix-journey-part-1-creating-a-flake/) blog post or a nice [article](https://blog.ysndr.de/posts/guides/2021-12-01-nix-shells/) that talks about some of the new commands.

As of now I'm "flaking" my projects and moving them from `shell.nix` projects to proper flake projects. After 4 months of using lorri I've now moved on to another tool because lorri doesn't support flakes. The exact benefits, why I moved over and how you can set it up will be a future blog post once I got some more experience with flakes.

A small teaser; flakes allow you to simply use `nix develop github:markbeep/markbeep.github.io` to get all the tools I use for working on this blog (basically just [hugo](https://gohugo.io/)) without even having to clone the project.
