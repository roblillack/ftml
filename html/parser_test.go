package html

import (
	"log"
	"strings"
	"testing"

	"github.com/roblillack/ftml"
	"github.com/roblillack/gockl"
	"github.com/stretchr/testify/assert"
)

func TestParsingSimpleParagraphs(t *testing.T) {
	tests := map[string]*ftml.Document{
		`<p>This is a test.</p>`:                              doc(p__("This is a test.")),
		`<p>one</p><p>two</p>`:                                doc(p__("one"), p__("two")),
		`<blockquote><p>one</p></blockquote><p>two</p>`:       doc(quote_(p__("one")), p__("two")),
		`<p><b>Bold</b> text.</p>`:                            doc(p_(b__("Bold"), span(" text."))),
		`<p><b>Bold<br />text.</b></p>`:                       doc(p_(b_(span("Bold"), span("\n"), span("text.")))),
		`<p>Test <b>bold</b></p>`:                             doc(p_(span("Test "), b__("bold"))),
		`<h1> Hello World! </h1>`:                             doc(h1_("Hello World!")),
		`<p>A<br/> B</p>`:                                     doc(p_(span("A"), span("\n"), span("B"))),
		`<ul><li><p>a</p></li><li><p>b</p></li></ul>`:         doc(ul_(li_(p__("a")), li_(p__("b")))),
		`<ul><li><p>a</p></li><li><p>b</p><p>c</p></li></ul>`: doc(ul_(li_(p__("a")), li_(p__("b"), p__("c")))),
		`<ul>
			<li>
				<ul>
					<li>
						<p>a</p>
					</li>
				</ul>
			</li>
			<li>
				<p>b</p>
				<p>c</p>
			</li>
		</ul>`: doc(ul_(li_(ul_(li_(p__("a")))), li_(p__("b"), p__("c")))),
	}
	for input, d := range tests {
		res, err := Parse(strings.NewReader(input))
		if assert.NoError(t, err) {
			assert.Equal(t, d, res)
		}
	}
}

func TestParsingParagraphsWithExtraTags(t *testing.T) {
	tests := map[string]*ftml.Document{
		`<title>bla</title><p>This is a test.</p>`: doc(p__("This is a test.")),
		`<!--[if !((mso)|(IE))]><!-- --><div class="hse-column-container" style="min-width:280px; max-width:600px; Margin-left:auto; Margin-right:auto; border-collapse:collapse; border-spacing:0; background-color:#003740; padding-bottom:25px" bgcolor="#003740"><!--<![endif]-->`: doc(),
		`<p>one</p><div><p>two</p></div>`:                     doc(p__("one"), p__("two")),
		`<blockquote><p>one</p></blockquote><p>two</p>`:       doc(quote_(p__("one")), p__("two")),
		`<p><b>Bold</b> text.</p>`:                            doc(p_(b__("Bold"), span(" text."))),
		`<p><b>Bold<br />text.</b></p>`:                       doc(p_(b_(span("Bold"), span("\n"), span("text.")))),
		`<p>Test <b>bold</b></p>`:                             doc(p_(span("Test "), b__("bold"))),
		`<h1> Hello World! </h1>`:                             doc(h1_("Hello World!")),
		`<p>A<br/> B</p>`:                                     doc(p_(span("A"), span("\n"), span("B"))),
		`<ul><li><p>a</p></li><li><p>b</p></li></ul>`:         doc(ul_(li_(p__("a")), li_(p__("b")))),
		`<ul><li><p>a</p></li><li><p>b</p><p>c</p></li></ul>`: doc(ul_(li_(p__("a")), li_(p__("b"), p__("c")))),
		`<ul>
			<li>
				<ul>
					<li>
						<p>a</p>
					</li>
				</ul>
			</li>
			<li>
				<p>b</p>
				<p>c</p>
			</li>
		</ul>`: doc(ul_(li_(ul_(li_(p__("a")))), li_(p__("b"), p__("c")))),
	}
	for input, d := range tests {
		res, err := Parse(strings.NewReader(input))
		if assert.NoError(t, err) {
			assert.Equal(t, d, res)
		}
	}
}

func TestTrimmingWhiteSpace(t *testing.T) {
	type Scenario struct {
		Before []ftml.Span
		After  []ftml.Span
	}
	for _, scenario := range []Scenario{
		{Before: []ftml.Span{}, After: []ftml.Span{}},
		{Before: []ftml.Span{span("test")}, After: []ftml.Span{span("test")}},
		{Before: []ftml.Span{span(" test")}, After: []ftml.Span{span("test")}},
		{Before: []ftml.Span{span("test ")}, After: []ftml.Span{span("test")}},
		{Before: []ftml.Span{span(" test ")}, After: []ftml.Span{span("test")}},
		{Before: []ftml.Span{span("\n  test, "), span("test.\n")}, After: []ftml.Span{span("test, "), span("test.")}},
	} {
		assert.Equal(t, scenario.After, trimWhiteSpace(scenario.Before))
	}

}

func TestParsingAndWritingStyles(t *testing.T) {
	simpleTests := map[string][]ftml.Span{
		// TODO: Unclear if we want to keep this behavior or just convert
		// &emsp14; to a space for HTML import, too (to stay consistent
		// with FTML reading).
		`This is a test.`:            {span("This is a test.")},
		`&emsp14;This is a test.`:    {span("\u2005This is a test.")},
		`This is a test.&emsp14;`:    {span("This is a test.\u2005")},
		`A&emsp14;&emsp14;&emsp14;B`: {span("A\u2005\u2005\u2005B")},
	}
	indentedTests := map[string][]ftml.Span{
		`This is a <b>test</b>.`:               {span("This is a "), b__("test"), span(".")},
		`This is a <b> test </b>.`:             {span("This is a "), b__(" test "), span(".")},
		`This is a <b><i>second</i> test</b>.`: {span("This is a "), b_(i__("second"), span(" test")), span(".")},
	}

	check := func(input string, doc *ftml.Document) {
		log.Printf("Running tests for “%s” ...\n", strings.TrimSpace(input))
		parsedDoc, err := Parse(strings.NewReader(input))
		assert.NoErrorf(t, err, "unable to parse input string: %s", input)
		assert.Equalf(t, doc, parsedDoc, "string parsed incorrectly: %s", input)
	}

	for input, spans := range simpleTests {
		check("<p>"+input+"</p>\n", doc(&ftml.Paragraph{Type: ftml.TextParagraph, Content: spans}))
	}

	for input, spans := range indentedTests {
		check("<p>"+input+"</p>\n", doc(&ftml.Paragraph{Type: ftml.TextParagraph, Content: spans}))
	}
}

func TestParsingInlineStyles(t *testing.T) {
	tests := map[string][]ftml.Span{
		`This is a <b>test</b>.`:               {span("This is a "), b__("test"), span(".")},
		`This is a <b> test </b>.`:             {span("This is a "), b__(" test "), span(".")},
		`This is a <b><i>second</i> test</b>.`: {span("This is a "), b_(i__("second"), span(" test")), span(".")},
	}

	for input, expected := range tests {
		z := gockl.New(input + "</END>")
		spans, nextPara, err := readContent(z, "END", ftml.TextParagraph)
		assert.Empty(t, nextPara)
		if assert.NoError(t, err) {
			assert.Equal(t, spans, expected)
		}
	}
}

func TestParsingInlineErrors(t *testing.T) {
	tests := map[string]string{
		// Inline errors
		`<p>This is a <b>test</i>.</p>`:                                "<p>This is a <b>test.</b></p>",
		`<p>This is a <b> test.`:                                       `<p>This is a <b> test.</b></p>`,
		`<p>This is a <b> test.</p>`:                                   `<p>This is a <b> test.</b></p>`,
		`<p>This is a <b><i>second</b> test</i>.</p>`:                  `<p>This is a <b><i>second test</i>.</b></p>`,
		`<p>This is a <b> test.<hr><b></p>`:                            `<p>This is a <b> test.</b></p>`,
		`<p>This is a <b><a href="asdasdasd">second</a> test</b>.</p>`: `<p>This is a <b>second test</b>.</p>`,
		`<p>This is a &copy;.</p>`:                                     `<p>This is a ©.</p>`,
		`<footer><div>This is a &copy;.</div></footer>`:                `<p>This is a ©.</p>`,
	}
	for input, expected := range tests {
		res, err := Parse(strings.NewReader(input))
		if assert.NoError(t, err) {
			buf := &strings.Builder{}
			assert.NoError(t, ftml.Write(buf, res))
			result := strings.TrimSpace(buf.String())
			assert.Equal(t, expected, result, "input:    %s\nexpected: %s\nresult:   %s\n", input, expected, result)
		}
	}
}

func TestUnclosedBlockElements(t *testing.T) {
	input := `<p>Hello<h1>World</h1>`
	expected := "<p>Hello</p>\n\n<h1>World</h1>"
	res, err := Parse(strings.NewReader(input))
	if assert.NoError(t, err) {
		buf := &strings.Builder{}
		assert.NoError(t, ftml.Write(buf, res))
		result := strings.TrimSpace(buf.String())
		assert.Equal(t, expected, result, "input:    %s\nexpected: %s\nresult:   %s\n", input, expected, result)

	}
}

func TestParsingErrors(t *testing.T) {
	tests := map[string]string{
		`This is a test.`:              "unexpected text content",
		`<p>one<p>two</p></p>`:         "non-inline token",
		`<blockquote>one</blockquote>`: "unexpected text content",
		`<p>one</blockquote>`:          "unexpected token",
		`<h1><p>one</p></h1>`:          "non-inline token",
		`<ul><p>boo</p></ul>`:          "content for list without list item",
		`<li>boo</li>`:                 "unexpected list item",
		`<ul></li>`:                    "unexpected closing tag for list item",
	}
	for input, expectedErrMsg := range tests {
		_, err := Parse(strings.NewReader(input))
		if assert.Error(t, err) {
			assert.Contains(t, strings.ToLower(err.Error()), expectedErrMsg)
		}
	}
}
