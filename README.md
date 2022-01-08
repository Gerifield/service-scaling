# Scaling examples

This is the repo for my (Hungarian) [Twitch streams](https://www.twitch.tv/gerifield) where I speak a bit about scaling an application. (There's also a Youtube backup for the videos here: [https://youtube.com/playlist?list=PLOIpbzPG9JERPYLml4sUiQ8sjt2awKu1J](https://youtube.com/playlist?list=PLOIpbzPG9JERPYLml4sUiQ8sjt2awKu1J) )

Each scale# folder more or less the same basic message board application, but I'll try to introduce new scaling methods and summarize the positive and negative sides of them.

Apps:

- [`scale0/`](scale0/) - Basic application where you can save a message and then retrive the message list. Nothing special, just a simple MySQL database save an fetch.
- [`scale1/`](scale1/) - First step for the actual scaling with some Redis based cache in the application.

