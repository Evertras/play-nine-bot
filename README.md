# Play Nine Bot

Plays [Play Nine](https://playnine.com/) with different strategies to figure out
how to win.

[Instructions
here](https://cdn.shopify.com/s/files/1/0503/3010/8062/files/Single_Page_Instructions_english.pdf?v=1695245685)

## Player turn decisions

A player strategy must make the following decisions:

First, draw or take from discard?

If draw, then decide:

- Replace either a face down or face up card on their board
- Discard the drawn card and flip up a facedown card
- OR, if there is only one facedown card, optionally skip the turn

If take from discard:

- Replace either a face down or face up card on their board

## Known tradeoff decisions

There are a bunch of tests in here that are technically flakey because they're
checking for randomness in seeds that change on each run. I'm being a little
lazy about passing custom seeds around and priming tests better that way.
