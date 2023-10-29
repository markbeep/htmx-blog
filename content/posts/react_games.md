---
title: "React.js Games"
date: 2022-04-07T09:14:40+02:00
showDate: true
draft: false
justify: true
tags: ["programming","react.js", "games"]
---

# React.js Games
Recently I was at a hackathon (StartHack 22) and there I did a lot of frontend for the app we produced. I told myself there I'll not write another line of Javascript or CSS for the weeks to come, because I had enough frontend for the time being. Then came the weekend and some people were playing Minesweeper on my Discord bot [Lecturfier](https://github.com/markbeep/Lecturfier). The problem with Minesweeper on Discord is though, that the playing field is simply one big field of spoiler-tagged characters, meaning there is no way to actually flag bombs and there's no way for the bot to know if you won or lost. That is why I decided to create Minesweeper on React.js, because I thought it will be a quick little project. After all I already made Minesweeper once. Somehow I wasn't able to stay from React.js for even a week. Darn.


## Minesweeper
Minesweeper was a nice little train ride project, that I at first wanted to only do in my train rides, but soon I just started to work on it throughout the whole day. It was a nice way to procrastinate. The programming also went quite smoothly. The source code is available over [here](https://github.com/markbeep/minesweeper). Most of the things worked, the CSS was also quite cooperative this time around and in the end I got what is now playable here on my blog under [https://markc.su/minesweeper](https://markc.su/minesweeper). There's currently no button going there. So it's kind of a hidden place on my blog. Just like [https://markc.su/video](https://markc.su/video) in a way. This is a little project I find turned out really well and it works surprisingly well, considering I made it.

## Tetris
After this little project I was feeling ready to take on another quick challenge. I got the suggestion to create Tetris. I thought simple enough, I should be able to do that in a short amount of time. Oh boy was I mistaken. Tetris is a bit harder than Minesweeper and there are a lot more things to consider. At first I tried to create a non-object oriented way of handling the falling of pieces where I simply had a grid with the pieces in it and what pieces are still falling, but that turned out to get me from behind, because when it came to rotating the piece it would've been an immense task to figure out where the current falling piece is, then figure out how it's located and how I need to move all the pieces in the grid.

I then moved to an object-oriented setup, which luckily was quite easy to change over to. This also made everything a lot easier and less buggy. I then simply had to define how the different rotations look exactly for each piece using a 4x4 boolean array and then this allowed me to render the pieces with their correct rotation. The project started to drag out and I was spending way too much time on it for my taste, so I left the design part a bit short and simply went for some blueish color theme. The source code for this is available over [here](https://github.com/markbeep/tetris-react). The final product is also playable on my blog under [https://markc.su/tetris](https://markc.su/tetris), but it is also hidden and there are no buttons leading to it from my blog.

There are a few bugs that I noticed later, but I didn't want to spend more time on it than I already did, so they're staying there for now. One of the major bugs you'll only notice if you're good anyway hehe. I'm curious what your highscores are for this Tetris though. Comment them below. As of writing this post my highscore is `18'040`. Let's see if you can beat it :)