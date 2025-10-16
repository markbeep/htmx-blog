---
title: "Popup Tools"
date: 2025-10-16T13:13:43+02:00
showDate: true
tags: ["Endeavour OS", "programming"]
draft: false
---

I'm a huge fan of using fuzzy matching for search. It's such a great solution for so many types of searches that is extremely underutilized. And for the terminal, [fzf](https://github.com/junegunn/fzf) is always my goto. I've been using fzf to allow me to quickly jump to any directory in the terminal for years already. But only recently have I started creating small popup tools that I can activate using some keybind. Here I wanna quickly explain how I have it set up.

I use [Niri](https://github.com/YaLTeR/niri) as my window manager. It's an amazing window manager I've really taken a liking to the past months. For some reason I have not released my blog post about it yet. So stay tuned. For this blog post, all you need to know is that Niri is a tiling window manager with some additional benefits I'll not go into here. In basics, that means that the window manager automatically splits up the screen so that each open window has a space somewhere. But you can also add a filter for certain types of programs to float around and not automatically get tiled. We'll use this to create our popup.

What I found was that it's extremely easy to create a quick CLI tool that pops up and allows you to select something. This saves me having to open the terminal and using a command there for things I just want to quickly execute and then have gone again. Let me give you an example.

### Install packages quickly

If I want to quickly install a package, what I can do, is simply press `Super+P`, which then opens my fuzzy finding package manager terminal window as a popup in the middle of my screen. I can now search for any packages I want with fuzzy matching. Additionally I also get all the relevant information like URL, author, last updated at, etc. to make sure I'm downloading the right package. Once the download is complete the popup nicely disappears and I can go on with my work.

![](/content/posts/popup-tools/paruz-popup.png)

You can create most of this without any third-party dependencies, using just pacman/yay/paru. I'll also mention how to below. In the screenshot above, and because of straight-up laziness, I myself use [paruz](https://github.com/joehillen/paruz).

### Launch VSC Workspaces

But why stay at only using it for installing packages? I'm a VSCode user, and on this system I use devcontainers for basically every project. I have more than enough RAM to handle Docker and it allows me to neatly manage all the dependencies without cluttering my whole system or having to write nix files.

Now one of the problems is, that when I want to open a specific devcontainer workspace, I first have to open VSCode and wait for it to load into the project. Probably I had some other devcontainer open, so I have to first wait until VSCode finishes initializing the dev docker container, before I can properly switch over to the workspace (using `Ctrl+R`) I was intending to go to.

So my task was to create some terminal script that allowed me to launch VSCode in a workspace. But now that always requires me to first open the terminal just to then open VSCode. So why not expand that into a popup just like I did with the package manager above.

![](/content/posts/popup-tools/vsc-popup.png)

Now I can simply run `Ctrl+T` (tbd, might change) and open a nice popup that allows me to pick the VSCode workspace I want to resume.

### Quick How-to

The config files I'll show are for my setup, so might have to be setup and translated (with gpt) into whatever you use:

- WM: Niri
- Terminal: Kitty
- Shell: Fish

#### WM

Add a window rule so that any apps opened with the ID `floating-popup` are not automatically tiled, but instead initially opened as a floating window:

```kdl
// .config/niri/config.kdl
window-rule {
    // floating terminal popup (like paruz and vscode workspace finder)
    match app-id="floating-popup"
    open-floating true
    default-window-height { fixed 400; }
    default-column-width { fixed 700; }
}
```

Then we add the keybinds to launch our scripts:

```kdl
// .config/niri/config.kdl
Mod+P hotkey-overlay-title="Pacman Install" { spawn "bash" "/home/user/.../paruz-popup.sh"; }
Mod+T hotkey-overlay-title="VSCode Workspace Launcher" { spawn "bash" "/home/user/.../vscode-popup.sh"; }
```

#### Fish

Now we come to the scripting part. Make sure to swap out the `kitty` terminal for whatever you use.

In the `paruz-popup.sh` I simply have the following:

```sh
kitty --class floating-popup -e bash -c 'paruz' &
```

Alternatively, if you want to use pacman/yay/paru without using paruz:

```sh
kitty --class floating-popup -e bash -c '
  pacman -Slq | fzf --multi --preview "pacman -Si {1}" \
  | xargs -ro sudo pacman -S;
' &
```

And now we're already finished. With `Super+P` we can easily install what we want.

For the VSCode workspaces, you can find all the recent workspace paths in your `.config` directory:

```sh
find ~/.config/Code/User/workspaceStorage -type f \
	-name "workspace.json" \
	-exec jq -r .folder {} +
```

Then using a vibe coded fish function, the base directory is extracted and shown in the fzf menu:

```fish
# .config/fish/functions/vscode_recent.fish
function vscode_recent --description "Open recent VSCode workspaces"
    #!/usr/bin/env fish

    # Collect all workspace folder URIs
    set entries (find ~/.config/Code/User/workspaceStorage -type f -name "workspace.json" -exec jq -r '.folder' {} + | sort -u)

    if test (count $entries) -eq 0
        echo "No VSCode workspaces found."
        exit 1
    end

    # Lists
    set basenames
    set uris
    set display_list
    set seen_hashes

    for uri in $entries
        set uri (string trim -r -c '/' $uri)
        set base (basename $uri)

        # Identify whether this is Remote or Local
        if string match -q "vscode-remote://*" $uri
            set label 'Remote '
        else if string match -q "file://*" $uri
            set label 'Local  '
        else
            set label Unknown
        end

        set basenames $basenames $base
        set uris $uris $uri
        set display_list $display_list "$label | $base"
    end

    # Adds the fzf popup/selection
    set selected (printf '%s\n' $display_list | sort | fzf --prompt="VSCode Workspaces > ")

    if test -n "$selected"
        set chosen_base (string split -f2 '|' $selected | string trim)
        for i in (seq (count $basenames))
            if test $basenames[$i] = $chosen_base
                code --folder-uri $uris[$i]
                break
            end
        end
    end
end
```

Placing this function into `.config/fish/functions/vscode_recent.fish` will allow you to execute it in the terminal and test it. Lastly, we just add the sh file that will launch it as a popup.

```sh
kitty --class floating-popup -e fish -c 'vscode_recent' &
```

### Looking back

Now having written this blogpost, I just realized that both of the shell scripts are so short, I could just put them directly into the niri config without additional files. But oh well, that's for you to then do.

But also remember that fzf basically just takes in a list and allows you to fuzzy match over it. You can use this for anything you can think of that uses text.
