---
title: "Keeping the system tidy with aconfmgr"
date: 2025-09-29T10:23:29+02:00
showDate: true
tags: [ linux nix arch ]
draft: false
---

This is just a quick blogpost to share an amazing tool I just stumbled upon, called aconfmgr.

## Wait what, don't I use nix?

As some might know, the past years I was an avid NixOS user. I was attracted by the whole idea of configuring your system using some config files. The Nix language allows you to split up and organize your files however you wanted, which meant I had a separate `modules.nix` file just listing all my installed packages. Anytime I wanted to install something system-wide*, I had to add it to my `modules.nix`. This forced me to think about anything I wanted to install on my system with maybe a comment on why I installed it (i.e. required for X to work), so when I do my occasional pruning I knew if I could remove it.

This made me quite conservative on what I installed and hence mostly relied on installing packages only for a shell session instead of system-wide. Read my more than two year old blogpost about this [here](/posts/nix-direnv-setup). But this whole ordeal with a config file OS also forced me to be tidy about what resides on my system. My system was basically in a permanent clean state with no garbage packages anywhere.

I have recently started using a new laptop and with that I felt it was time to change things up again. I was also starting to get annoyed at simple things like binaries with shared libraries not working out of the box on NixOS (only took me almost 3 years). I initially went for an immutable fedora, but that wasn't able to fulfill the itch I was looking for. Upon stumbling upon the aforementioned tool, I knew exactly what distro I have to run.

## The itch fulfilled

I was looking for nix alternatives in terms of package management, curious to see what was out there. I knew about ansible and was even close to writing my own specific package management that automatically installed and exported packages in distroboxes using ansible, but luckily didn't pull through with it.

I then stumbled upon [aconfmgr](https://github.com/CyberShadow/aconfmgr). The basics are, it creates a diff of expected and installed packages and files which is all outputted into a new "unsorted" config file.

## The Workflow

The idea with aconfmgr is "install first, worry and sort later", which is a nice change to NixOS's approach of "worry and sort first, install later".

So I install what I want as normal. Oh, I need X for this course. I can quickly install it and start using it. But let's say it didn't quite work and i also had to install the dependency D for X to work.

Afterwards I can then run `aconfmgr save` in the terminal and it will generate a "unsorted" config file which will then list X and D. The config file looks like this:

```bash
# 99-unsorted.sh
AddPackage X
AddPackage D
AddPackage SomeGarbagePackage
```

I can then take X and D and sort them into one of my other config files together with a quick description of why I need it:

```bash
# 40-course-modules.sh

# Required for course Y HS2025
AddPackage X
AddPackage D
```

I can also see that I seem to have installed some `SomeGarbagePackage`. I don't need it so I can either just delete it manually using my package manager or simply not sort it into any of my module files.

After sorting my packages I would then delete the unsorted config file (`99-unsorted.sh`) and then upon running `aconfmgr apply`, aconfmgr will automatically delete any packages I don't have in my config.

I then also add my aconfmgr config files to my gitified dotfiles so I can also see when I committed which line.

## And there's more

`aconfmgr apply` doesn't only remove unlisted packages, it also adds uninstalled but listed packages. It also doesn't track packages, but it also keeps track of files (like the config files in /etc).

So in theory you can use it to backup your system state and if you ever decide to reset your system, you can hop back with the correct packages and files.

But I don't really use it for that and mainly use it to nicely keep track of my packages.

Another con is that it only works for arch-based distros. I was open to distro-hopping anyway, so this fit quite well. So yes, I'm on Endeavour OS now.  

Have a go, try out aconfmgr and see if it also clicks for you like it did for me. A must-have if you want to keep your system clean of random packages.
