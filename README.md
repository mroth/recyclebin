# :recycle: recyclebin
> A tool to measure usage of a particular term on Twitter.

Why "recyclebin"?
-----------------
Because my primary motivation to write it was figuring out what's going on with
the :recycle: symbol on Twitter. (But you can use it to measure anything via the
`-term=FOO` option.)

See the article I used it for here:
https://medium.com/@mroth/why-the-emoji-recycling-symbol-is-taking-over-twitter-65ad4b18b04b

Screenshot
----------
<img width="696" height="743" src="https://user-images.githubusercontent.com/40650/32742490-99bd1e02-c877-11e7-8ee1-657ed40e942a.png" alt="screenshot">


Here be dragons
---------------
This was written as a quick one-off program to just grab some data for an article.
Please treat the code with that understanding. :-)

Using this yourself
-------------------
Compile, set the appropriate Twitter API keys as environment variables and run!
See `--help` for available options. If you need more assistance, ping me or open
an issue and I'll make this more obvious.

Please note that I have full Partner level access to the Twitter Streaming API on
Emojitracker’s dev account — if you don’t, when you try to use this script to
track anything super high volume, tweets will be dropped from your results, and
your numbers will reflect a sampling rather than "everything". Plan accordingly.
