---
title: "Getting Started with Hugo"
date: 2021-04-15T18:24:16+02:00
showDate: true
mathjax: true
tags: ["web", "blog"]
---

> _EDIT 21.11.2023: This blog isn't set up using hugo nor hosted on Github pages anymore.
> I've learned a lot in terms of web development since this post. I've remade the full
> blog using Go templating and htmx._

# Welcome

Welcome to the grand opening of my new website. This is the very first post on here.

This website is mainly to test out new stuff in the direction of web development, as that is a direction
I don't currently have a lot of experience with. Additionally, a plan I have in mind is to write quick summaries for topics
I'm currently learning about. This would be useful for myself to look back at when learning for exams and it might even
be interesting for some people to read about the topics. As I'm studying computer science they will mostly be about CS and math.

## Setting up this website

Settings up this website was already a challenge for me, as I'm a real doofus when it's about anything web or internet.

### Hosting

For hosting this website I use [GitHub Pages](https://pages.github.com/). It allows you to host static websites for free.
And I'm not going to say no to free. I already use GitHub for all my other code-related projects, so it's perfect to have this
site also be stored on there. Click the GitHub button at the bottom of the page to visit my GitHub, where you can also see
the source code of this site.

### Domain

The second task on the list was to create a **sick** domain. It needs to be short and meaningful. Nobody wants to type a
long domain with tons of characters every time they visit the site. Unlike this domain with 71 characters in the labels:
http://www.abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk.com/ (this site actually exists).

The domain I chose to go with is made up with my first name and the first three letters of my surname. Additionally, the origin
of the [.su](https://en.wikipedia.org/wiki/.su) top-level domain is pretty interesting. If you want to find a cool top-level domain,
a good place to start is just scrolling through the domains on the [_List of Internet Top-Level Domains_](https://en.wikipedia.org/wiki/List_of_Internet_top-level_domains).
There are tons you surely have never even heard about.

### Layout/Theme

For deciding the theme or basis of the website, there are tons of methods one can go. I decided against creating a website from scratch, as that
takes a lot of time to create a reasonable site, as I have **0** experience. Using frameworks and/or using pre-made themes allows for quickly setting up
a website with the theme you want. It came down to chooing between Jekyll and Hugo for me, as they are both relative easy to set up and take a lot of work
off your shoulders. In the end, I decided to go with [Hugo](https://gohugo.io/) as I preferred the themes and additionally themes are all free. I fell in
love with this theme. It's called [Call me Sam](https://themes.gohugo.io/hugo-theme-sam/). There were only a few things to do to get up and running. For the
background I just use a video downloaded from [here](https://www.pexels.com/video/a-mist-over-water-2534297/). Then I clapped it into my video editor to
make it a seamless loop. I then simply had to place the video file in the `static` folder for it to work.

### 404 Page

Fourth item on the list was to have a good looking 404 page. Everybody knows, that the key to creating a good website, is impressing
users even when you mess up. I picked the 404 page from [freefrontend.com](https://freefrontend.com/html-css-404-page-templates/).
freefrontend.com offers tons of web code snippets you can just clap into your own website to save a lot of effort. And of course, it's free!
Both the html and css part is in the same html code, to avoid any trouble.

### $\LaTeX$ & `Code Blocks`

As I plan to use this blog/website (still a bit unsure what it's really gonna be) for topics related to my studies, making sure that
code blocks and math formulas can be shown is of high importance. By default, all posts/pages are made in markdown when using hugo.
This means all the things that can be made in the usual markdown can also be made here. This allows me to also post code snippets
in any language in the following form:

```java
/*
 * Inefficient way of multiplying a * n
 */
public static int foo(int a, int n) {
    int total = 0;
    for (int i = 0; i < n; i++) {
        total += a;
    }
    return total;
}
```

For math syntax $\LaTeX$ is the best way to go. To make it work on here I use the [MathJax](https://www.mathjax.org/) engine.
With Hugo you have a `layouts` folder inside your `/root` and `themes/<theme name>` folder. The files in your root's
layout folder overwrite the default files from the theme.

We first create the file that then turns the needed parts into LaTeX. We create a new file `mathjax.html` in `./layouts/partials/`
with the following content:

```html
<script>
  MathJax = {
    tex: {
      inlineMath: [
        ["$", "$"],
        ["\\(", "\\)"],
      ],
    },
    svg: {
      fontCache: "global",
    },
  };
</script>
<script
  type="text/javascript"
  id="MathJax-script"
  async
  src="https://cdn.jsdelivr.net/npm/mathjax@3/es5/tex-svg.js"
></script>
```

To make the MathJax script above work in all markdown pages on the whole site, it's best to add it into the `<head>` tag, which is loaded on all pages. We copy the `head.html`
file from `./themes/sam/layouts/partials/` to `./layouts/partials/` and add the following line into the copied `head.html`.

```
{{`{{ if .Params.mathjax}}{{ partial "mathjax.html" . }}{{ end }}`}}
```

It doesn't matter what line exactly, as long as `head.html` gets called in the `<head>` tag.

This is all the setting up that we have to do. To enable MathJax we simply add the parameter `mathjax: true` at the top of every markdown file we want to use LaTeX in. The top of this markdown page looks like this for example:

```
---
title: "Getting Started with Hugo"
date: 2021-04-15T18:24:16+02:00
showDate: true
mathjax: true
tags: ["web", "blog"]
---
```
