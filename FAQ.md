### How can I display FTML content without having FTML-capable software available?

- `w3m -T text/html -cols 73 MYFTMLDOCUMENT.FTML | less -r`
- `lynx -force_html -dump examples/rfc1896.ftml | less`
