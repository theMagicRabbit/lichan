# Lichess Analyzer

Downloads games from lichess and analizes them locally with stockfish.

## Motivation

Lichess is a free website for playing chess. As they describe themselves:

> lichess.org is a free/libre, open-source chess server powered by volunteers and donations.

Analyzing games with stockfish is one of the features of lichess, but this is a manual process
and costs server compute time. This application allows a user of Lichess to download and analize
their games locally. This reduces server load and cost for Lichess and while also providing detailed
game analysis for study.

The analyzed games are in standard PGN (portable game notation) format which can be read by the user
or viewed through a number of chess programs such as (pychess)[https://pychess.github.io/]

