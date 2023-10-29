---
title: "Nix -  The best package manager you should be using"
date: 2023-05-23T22:21:19+02:00
showDate: true
tags: ["nix", "nixos", "package manager"]
---

This is an article I wrote for the upcoming
[Visionen](https://visionen.vis.ethz.ch/) issue "Nachwuchs". The issue will be
released in the coming weeks and this serves as a little sneak-peak.

---

Have you ever tried to install a program, only for it to require some specific
Python version, GCC to be exactly version 11, or some other random program you
never even heard of? Additionally, you have to figure out where to install what
and how, because everybody has their standard on how programs should be
installed. A solution for this is a package manager which handles the complete
install for you. From checking for dependencies down to what steps have to be
executed for the program to be installed correctly.

There are a lot of package managers around. Some of which you might’ve even used
already. A few popular examples include Homebrew for OSX, Chocolatey for
Windows, and APT for Debian and Ubuntu. They make installing programs so simple.
You can simply look up an application in their online glossary and then execute
a quick command for the application to be automatically installed for you. But
have you ever wondered how they work?

The package managers all operate slightly differently, but in general, they all
share the same elements. They contain a huge list of installation instructions.
Each instruction contains all the details required to install the package on a
system: how to properly install a package, what dependencies are required, what
version the package is, where to download it from, etc.

Over the last few months, I've been playing around on NixOS which is an
interesting take on a Linux distro, but the details are for another Visionen
article (or blogpost). The important part is that NixOS comes with a package
manager which is simply called Nix. Nix can be quite quickly installed on any OS
(even OSX) and comes with a lot of nice goodies. The main selling point that
made me fall in love with Nix is the very organized structure for installing
programs. Installing a program doesn’t require you to handle and install
dependencies manually. Additionally Nix makes it extremely easy to try out
programs. The garbage collector will handle removing packages you don’t need
anymore. That way you can try out whatever you want. Anything you don’t need
will then be cleaned up for you without ever having to think about it again.
When was the last time you cleaned out programs and packages you don’t need
anymore? I barely remember the times I did before using Nix. Now I do it daily
in seconds without having to ever worry about breaking my installs. Nix does
this all without lacking in the amount of packages it supports as shown on
repology.org.

So how does installing differ from other package managers? Installing a package
with dependencies does not depend on what is already on the system. Every
installation is handled in a separate environment. Let's say for example we want
to install a program which requires a specific Python version and some Python
library to be built on a clean system. Often what happens with standard package
managers is that they try to install a package and then either one of the two
following two actions takes place:

1. The package manager simply fails, stating that the dependencies are missing.
   Well great, now we have to install all the dependencies manually.

2. But more often, the package manager then notices the dependencies, asks if
   you want to install them, and then installs the dependencies as well. Now
   that is cool because it works. But you now have some random dependency
   installs that you didn't even intend to install and they were solely
   installed for the building step. You might've already had a Python install
   before and now you have yet another one with some other version. While this
   method works and is very often utilized, do we want our system to be filled
   with random programs that we never even intended to install?

Nix specially combats this. Nix installs all dependencies in a read-only path in
a special Nix directory (/nix/store). This ensures that when Nix installs Python
3.10, it will always be Python 3.10 and not accidentally have been modified by
the user or some non-standard program. Programs that are installed system-wide
will then be correctly linked up to the PATH to be able to be called from
anywhere.

Nix keeps a close watch on which of the installed packages are required by other
programs or are just lying around. With a simple command, Nix is then able to
automatically garbage collect all packages that are not required anywhere
anymore. This allows Nix to delete packages that were only required to build
another package but are not required afterwards or packages that were only
installed to quickly try out something.

So how does this trying out of programs work? Say you are looking to test out
some photo editors. People say GIMP works nicely. I want to try that out. You
can open the terminal, run “nix-shell -p gimp” and bam you are in a terminal
environment with GIMP in the PATH. In simple terms, this means I can now run
GIMP from this terminal session.

Okay, I tried it out and I don’t like it. I can now just close GIMP and my
terminal, then when I run my Nix garbage collector it will see that GIMP isn’t
properly installed anywhere and delete it for me. I’ve gotten so used to this
workflow that I barely install anything globally anymore. I install whatever I
need depending on the directory or project I’m currently in. This can then be
automated so you never even have to execute anything manually and always have
the one single correct Python version in whatever Python project you’re in.

Now think back to the last programs you downloaded. Possibly even over a package
manager. Now think of the last time you went through your installed programs and
deleted programs you didn’t need anymore. Are you sure you didn’t miss one of
those hundred Python installs? Join the Nix gang and you’ll never have to worry
about manually cleaning everything up anymore.
