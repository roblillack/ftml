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
		`This is a test.`:            {span("This is a test.")},
		`&emsp14;This is a test.`:    {span(" This is a test.")},
		`This is a test.&emsp14;`:    {span("This is a test. ")},
		`A&emsp14;&emsp14;&emsp14;B`: {span("A   B")},
	}
	indentedTests := map[string][]ftml.Span{
		`This is a <b>test</b>.`:               {span("This is a "), b__("test"), span(".")},
		`This is a <b> test </b>.`:             {span("This is a "), b__(" test "), span(".")},
		`This is a <b><i>second</i> test</b>.`: {span("This is a "), b_(i__("second"), span(" test")), span(".")},
	}

	check := func(input string, doc *ftml.Document, output string) {
		log.Printf("Running tests for “%s” ...\n", strings.TrimSpace(input))
		parsedDoc, err := Parse(strings.NewReader(input))
		assert.NoErrorf(t, err, "unable to parse input string: %s", input)
		assert.Equalf(t, doc, parsedDoc, "string parsed incorrectly: %s", input)
		buf := &strings.Builder{}
		assert.NoErrorf(t, ftml.Write(buf, doc), "unable to write FTML for doc: %s", input)
		assert.Equal(t, output, buf.String(), "written FTML not corrent")
	}

	for input, spans := range simpleTests {
		// expectedOutput := "<p>" + strings.TrimSpace(input) + "</p>\n"
		expectedOutput := "<p>" + input + "</p>\n"
		check("<p>"+input+"</p>\n", doc(&ftml.Paragraph{Type: ftml.TextParagraph, Content: spans}), expectedOutput)
	}

	for input, spans := range indentedTests {
		// expectedOutput := "<p>\n  " + strings.TrimSpace(input) + "\n</p>\n"
		expectedOutput := "<p>" + input + "</p>\n"
		check("<p>"+input+"</p>\n", doc(&ftml.Paragraph{Type: ftml.TextParagraph, Content: spans}), expectedOutput)
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
		spans, err := readContent(z, "END")
		if assert.NoError(t, err) {
			assert.Equal(t, spans, expected)
		}
	}
}

func TestParsingInlineErrors(t *testing.T) {
	tests := map[string]string{
		`This is a <b>test</i>.</END>`:               "unexpected token",
		`This is a <b> test.`:                        "unexpected eof",
		`This is a <b> test.</END>`:                  "unexpected token",
		`This is a <b><i>second</b> test</i>.</END>`: "unexpected token",
		`This is a <b> test.<hr><b></END>`:           "unexpected token",
	}

	for input, expectedErrMsg := range tests {
		z := gockl.New(input)
		_, err := readContent(z, "END")
		if assert.Error(t, err) {
			assert.Contains(t, strings.ToLower(err.Error()), expectedErrMsg)
		}
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
