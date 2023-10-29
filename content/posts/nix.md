---
title: "I use NixOS btw"
date: 2023-03-06T13:36:25+01:00
showDate: true
tags: ["programming", "nix", "linux"]
---

# Table of Contents

- [Table of Contents](#table-of-contents)
- [Introduction](#introduction)
- [Fun Study Phase Project](#fun-study-phase-project)
- [Purely functional? What?](#purely-functional-what)
- [How do I install a program?](#how-do-i-install-a-program)
- [Global Install](#global-install)
- [So this is the perfect OS, why doesn't everyone use it?](#so-this-is-the-perfect-os-why-doesnt-everyone-use-it)

# Introduction

In the world of Linux distros, there are quite a few unique distros. Just to mention two amazing ones of the bunch:
[Hannah Montana Linux](https://hannahmontana.sourceforge.net/) and [Biebian](https://biebian.sourceforge.net/)
can improve your Linux experience with awesome pre-installed theming. [NixOS](https://nixos.org/) is another unique distro
that goes a few steps further and completely changes how the underlying system works. How?
By making the whole system configuration purely functional.

---

# Fun Study Phase Project

After spending almost all my five ETH semesters using PopOS, I decided it was time for a change. The latest update broke a lot of things for me and
there were numerous things I installed for one-time use with a big dependency mess. And because I thought studying on a working distro was too easy,
I decided to start with a fresh NixOS install right as the study phase began. My first
few days were spent setting up the basics, such as how to connect a second monitor or figuring out what the best PDF readers are.
I spent countless hours learning how to do what, but luckily I was able to justify it by saying I was studying for the Computer Systems exam.

---

# Purely functional? What?

NixOS is a fully fletched Linux distribution that is completely built on top of Nix (published in
2004 in a [research paper](https://edolstra.github.io/pubs/nspfssd-lisa2004-final.pdf)). Nix is a fully functional and declarative
programming language. If you have some
experience with [Haskell](https://www.haskell.org/), then it'll be quite easy to understand Nix's syntax.

For example, we can sum up a list of strings by mapping them to integers:

```nix
foldl add 0 (map toInt ["1" "2" "3" "4"]) # returns 10
```

A big part of Nix's language is the object type. Objects are like dictionaries in Python or Objects in JavaScript. We can create an object with
multiple values as follows:

```nix
{
    name = "Foo";
    file = fetchurl "www.random.org/download";
    path = ./imports/image.png
}
```

This will create an object where the `name` attribute is simply the string `"Foo"`. The `file` and `path` attributes are where things get interesting.
The `fetchurl` function is a built-in function in Nix and it simply downloads the given URL, stores the result in a specialized path (/nix/store/) and
then returns the path to the downloaded file. That means our `file` attribute simply contains the path to the downloaded file. `path` here allows us to use a path
and automatically have Nix check if it's valid. These objects are what NixOS builds upon.

---

# How do I install a program?

This at first sounds like a stupid question. Just throw programs into /bin or /usr/bin and be done with it.
But that's where you would be wrong on NixOS. Instead of installing everything globally into the /usr/bin
directory like on every other distro, NixOS doesn't use the /bin or /usr/bin directories at all:

```bash
~ ls /bin
sh
```

Instead, programs are installed into /nix/store/ and then properly symlinked to work as normal.

So if I'm not supposed to clap my programs into /bin or /usr/bin, how am I supposed to install programs? Well, that is where the unique parts of NixOS start to show.

If I want to install a program, there are a few ways to do it. If for example, I want to install a program for a quick one-time use, say for example I want to quickly
run `dig`, but I probably won't use it in the future. I can open my terminal and simply run `nix-shell -p dig` to start up a custom terminal environment in which I have
`dig` installed. Or if I'm in a project that requires a specific Python version, I can just run `nix-shell -p python311` to jump into a terminal environment where I have
Python 3.11 installed. Now when I exit and reopen the terminal, I won't have Python in my path anymore.

```bash
~ python -V # no python at the start
The program 'python' is not in your PATH.
~ nix-shell -p python311
~ python -V
Python 3.11.1
~ exit # exit the shell again
~ python -V
The program 'python' is not in your PATH.
```

Now you might already see how this can be extremely useful when working on many different projects that all have different version requirements. If on one project
I need Python 3.8 while on another one I need Python 3.11 I can just install the correct Python version when in a directory. On my NixOS setup, I don't have any globally installed
Python version. By doing this with all my programs I completely declutter my global namespace and I can simply run the NixOS garbage collector to get rid of all
the programs I don't use anymore in a clean way. When was the last time you used that outdated Python install on your system? By never installing it globally, I
never have to go through the hassle of uninstalling it and worrying about some weird dependencies hanging around.

Of course, typing `nix-shell -p <package>` in every project I enter can get annoying quite fast. Especially if there are multiple packages I need to have installed. For this,
you can write a `shell.nix` file in which you can describe how and what should be installed. For example, it could look like this:

```nix
with import <nixpkgs> {};
with pkgs;
mkShell {
    buildInputs = [
        python311
        poetry
    ];
}
```

I can now simply run `nix-shell` in the terminal and it will automatically look at the shell.nix file, see that I want the two inputs Python 3.11 and
[Poetry](https://python-poetry.org/), then install them for the current environment. But I'm lazy, running `nix-shell` every time is already too much effort.
That is why I also use [direnv](https://direnv.net/) which automatically loads the nix-shell environment when I open the directory and automatically
unloads it when I exit again. Currently, I've also been keeping my eyes on [lorri](https://github.com/nix-community/lorri) which is a more optimized version
of direnv.

```bash
~ cd Python311Project/
~ python -V
Python 3.11.1
~ cd ..
~ python -V
The program 'python' is not in your PATH.
```

This also means that if I install something that would somehow break the environment, I can simply exit the terminal and because of how Nix installs packages,
I can be safe that the program will eventually be deleted by the garbage collector.

Okay, that's all cool and dandy, but what about programs I want globally installed? Like vim, firefox or VSCode?

---

# Global Install

Now how do I globally install a program? I don't want to constantly have to write `nix-shell -p vscode` when I want to run VSCode. This is where the famous
NixOS configuration files come into play. Often with other distros, you have some configuration here, some configuration in .config and system settings stored
yet again somewhere else. In NixOS you have all the configurations in /etc/nixos/configuration.nix. The configuration.nix file is the holy file of NixOS. It has
everything defined on what should be installed down to what drivers and services to run. Now you might think this is some thousand-long nightmare configuration file,
but this is where the power of using a functional language as a configuration comes into action. It is possible to split up the configuration file into a lot of
files to easily split up all the configurations into bite-sized chunks.

```nix
{
    imports = [
        ./nixos/modules
    ]
}
```

This quick snippet checks for the `default.nix` file in the ./nixos/modules/ directory and then copies over all the attributes from there. With this, the configuration
file can then nicely be split up into small modules.

In the configuration.nix file is also where all the system-wide applications are installed. This includes programs like the i3 window manager and also my general
programs like alacritty or vim.

```nix
{
    environment.systemPackages = with pkgs; [
        i3
        alacritty
        vim
        git
    ]
}
```

If I want to add a new system-level program, I simply look it up on the [NixOS package repository](https://search.nixos.org), add it to the list of packages and
rebuild the system with `sudo nixos-rebuild switch`. This completely evaluates the configuration file and creates a new generation. If anything
breaks, I can always roll back to a previous generation without any issues. This means I can test out any window manager or driver I want and if anything breaks,
I can select an older generation. NixOS will then evaluate the old configuration file and because everything is declarative it can set up the OS
just the way it was before.

For user-specific packages and to create my non-nix config files (like NVim's init.vim), I use [home-manager](https://github.com/nix-community/home-manager).
It is an add-on that handles user packages. It allows me to store all my configurations for my programs (like init.vim, alacritty.yaml, etc.) in my
NixOS Github repository and have home-manager automatically install them in the correct place.

If I were to set up a new device, I could simply clone my Github repo, run about three commands and the new device will be set up just like my current setup.

---

# So this is the perfect OS, why doesn't everyone use it?

Because this blog post is already turning out longer than intended and I don't want to write a book over here, I decided to not go into the full
details of how the package manager and packages work. I will leave that for a future blog post because how that all works is extremely interesting.

But in case this blog post piqued your interest, good. It's for sure one of the Linux distros with the highest learning curve because you have to
learn a whole new language to even start doing things. But once it works, it is extremely satisfying. I hope to see some fellow Nixers :)

If you want to read more, have a look at [How Nix Works](https://nixos.org/guides/how-nix-works.html) by NixOS themselves.
