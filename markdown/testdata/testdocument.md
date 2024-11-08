# This is an example FTML document

The purpose of this document is to demonstrate the following use-cases:

- Describe the features supported by FTML.

- Showcase the FTML standard formatting enforced by `fmtftml`.

## Supported features

### Inline styles

The following inline styles are supported in FTML:

- **Bold** text.

- _Italic_ text.

- <mark>Highlighted</mark> text.

- <u>Underlined</u> text.

- ~~Striked~~ text.

- Text formatted as `code`.

### Text paragraphs

This is a rather long paragraph to show the text-wrapping supported by the ASCII and Markdown exporters. Applications integrating FTML support should “do the right thing” regarding breaking lines in paragraphs. This could mean that the full available width of the window or screen is used, or that the lines are broken to limit the number of characters per line to a specific number. FTML itself _does not_ encorce any specific rules here, the included Markdown exporter does wrap at 80 characters, though.

### Quoted paragraphs

> Like any other blocks, quotes may contain any of the different paragraph types as children. First, let's look into wrapping of quotes, though. As you can see, this is a very long line again.
> 
> This is the second paragraph of the very same quote. This paragraph, too, will be broken into multiple lines, if necessary.
> 
> ### Quoted Headers
> 
> As expected, you'll be able to quote _all_ kinds of paragraphs—including section headers.
> 
> > Now, finally, let's look at a second level quote. Like expected, this paragraph will also wrap at the same width, even if if it indented further then the paragraphs above.

### Bullet points

TBD

**Example:**

- This is a bullet point.

- This is another bullet point.
  
  Contrasting to the first bullet point, this one contains multiple paragraps, with this one being the second one.

- This is a bullet point with a bunch of hard line breaks.

- > This bulletpoint contains a quote.

### Paragraph nesting

FTML supports nesting of paragraphs, so this is entirely possible:

> Please see, how the following list is part of a quote and contains nested paragraphs.
> 
> - This is a paragraph inside of a quoted paragraph
> 
> - This bullet points contains another quote:
>   
>   > You can never have enough nesting of paragraphs.\
   —Robert Lillack
> 
> - 1.  One
>   
>   2.  …
>   
>   3.  …
>   
>   4.  …
>   
>   5.  …
>   
>   6.  …
>   
>   7.  …
>   
>   8.  …
>   
>   9.  …
>   
>   10. aaaaand
>       
>       Ten!

## Whitespace and line-break handling

This is a line with a hard line-break:\
First line\
Second line.

This is a line with multiple line-breaks:\
First line\
\
Second line.

This is a line with multiple spaces:\
A   B

> This is a paragraph that contains a very long line of <mark>highlighted text to force the formatter to break\
the\
line\
in the middle.</mark> But afterwards, of course, things should continue normally.

## Unicode support

The following three tests check for the word wrapping to break at the right spot.

72 ASCII characters:

######## ######## ######## ######## ######## ######## ######## ########.

=======================================================================^

8 Unicode multi-byte characters, 64 ASCII chars:

cafébabe cafébabe cafébabe cafébabe cafébabe cafébabe cafébabe cafébabe.

=======================================================================^

32 double-with emoji, 7 spaces, 1 dot:

😎😎😎😎 😎😎😎😎 😎😎😎😎 😎😎😎😎 😎😎😎😎 😎😎😎😎 😎😎😎😎 😎😎😎😎.

=======================================================================^
