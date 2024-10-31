_Warning: This is a work in progress. The specification is not final and this repository really only is public because I'm using the library as a dependency in some other projects._

# FTML – Formatted Text Markup Language

A WIP specification (and Go library) of a HTML-compatible text markup language which only contains the most basic formatting options.

## Design goals

1. File format needs to be easily diffable—even for tools that are not aware of the file format.
2. Try to reduce ambiguoties: There's usually exactly _one_ way of expressing something in FTML.

## Design constraints

1. Must be upwards compatible to HTML 5 (not as full documents, but the FTML needs to be embeddable into a HTML5 file as is)

2. TBD

### Non-goals

- Support for:
  - Tables
  - Embedded images
  - Checklists
  - Colors
  - Font families
- Automatic replacing of emoji, quotes, arrow symbols, dashes, etc.
- Ease of editing for humans
- Support for other chartsets than UTF-8

### Potential goals which are not decided, yet

- Code blocks
- Links
- Footnotes

### TODO

- How to set up Prettier
- How to set up Visual Studio Code/Vim/Emacs
  - Live Preview: https://marketplace.visualstudio.com/items?itemName=ms-vscode.live-server

### Prior art

- [RFC 1896](https://datatracker.ietf.org/doc/html/rfc1896): The text/enriched MIME Content-type
- [Enriched Text Format -- A Primer](http://users.starpower.net/ksimler/eudora/etf.html): Some more information about `text/enriched`; including a list of compatible clients

### References

- [Unicode Spaces](https://jkorpela.fi/chars/spaces.html): What lead to using EMSP14 as non-collapsing white-space entity in FTML files.
