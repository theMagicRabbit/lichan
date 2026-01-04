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

## Quick Start

Lichan should run on any modern Linux distro.

Requires a Lichess account and API access token. After create a Lichess account, a
(token can be created here.)[https://lichess.org/account/oauth/token/create?]

Requires stockfish to be installed to the user path.
(Stockfish can be downloaded here.)[https://stockfishchess.org/download/]

To compile Lichan, (you will need Go installed.)[https://go.dev/doc/install]

Once the requrirements are completed, you can start install Lichan using the go package
manager. Run this command in your terminal:
`go install theMagicRabbit/lichan`

Copy the contents of (the sample config file)[sample_config.toml] to `~/.config/lichan/config.toml`
on your system. Enter your token as the `pat` variable and list the usernames you wish to analize
in the usernames variable list.

## Usage

Run Lichan by typing `lichan` into a command prompt.

Lichan is intended to be run from a cron or systemd timer. This allows automated processing of any recent
games from the accounts that are being tracked with Lichan.

## Contributing

Contributions to Lichan are welcome. If you'd like to contribute, please fork the repository and open a 
pull request to the `main` branch.

## License

Lichan is licenced under the GPL-3.0 or later to respect the licensing of Stockfish and Lichess.

