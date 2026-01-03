FROM ubuntu:24.04

COPY stockfish/stockfish-* /bin/stockfish
COPY lichan /bin/lichan

CMD ["/bin/lichan"]

