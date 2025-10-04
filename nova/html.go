package nova

import (
	"fmt"
	"html"
	"strings"
)

// HTMLElement represents any HTML element that can be rendered to a string.
// This interface allows for composition of complex HTML structures using
// both predefined elements and custom implementations.
type HTMLElement interface {
	Render() string
}

// Element represents an HTML element with attributes, content, and child elements.
// It supports both self-closing elements (like `<img>`) and container elements (like `<div>`).
// Elements can be chained using fluent API methods for convenient construction.
type Element struct {
	tag        string
	content    string
	attributes map[string]string
	children   []HTMLElement
	selfClose  bool
}

// textNode represents raw text content that should be rendered without HTML tags.
// It implements HTMLElement to allow text to be mixed with HTML elements in compositions.
type textNode struct {
	text string
}

// DocumentConfig provides configuration options for creating a new HTML document.
// Fields left as zero-values (e.g., empty strings) will use sensible defaults
// or be omitted if optional (like Description, Keywords, Author).
type DocumentConfig struct {
	Lang        string        // Lang attribute for `<html>` tag, defaults to "en".
	Title       string        // Content for `<title>` tag, defaults to "Document".
	Charset     string        // Charset for `<meta charset>`, defaults to "utf-8".
	Viewport    string        // Content for `<meta name="viewport">`, defaults to "width=device-width, initial-scale=1".
	Description string        // Content for `<meta name="description">`. If empty, tag is omitted.
	Keywords    string        // Content for `<meta name="keywords">`. If empty, tag is omitted.
	Author      string        // Content for `<meta name="author">`. If empty, tag is omitted.
	HeadExtras  []HTMLElement // Additional HTMLElements to be included in the `<head>` section.
}

// HTMLDocument represents a full HTML document, including the DOCTYPE.
// Its Render method produces the complete HTML string.
type HTMLDocument struct {
	rootElement *Element // The root `<html>` element.
}

// Render converts the HTMLDocument to its full string representation,
// prepending the HTML5 DOCTYPE declaration.
func (d *HTMLDocument) Render() string {
	if d.rootElement == nil {
		// This should ideally not happen if Document() is used correctly.
		return "<!DOCTYPE html>\n"
	}
	return "<!DOCTYPE html>\n" + d.rootElement.Render()
}

// Document creates a complete HTML5 document structure, encapsulated in an HTMLDocument.
// The returned HTMLDocument's Render method will produce the full HTML string,
// including the DOCTYPE.
//
// It uses DocumentConfig to customize the document's head and html attributes,
// and accepts variadic arguments for the body content.
// Sensible defaults are applied for common attributes and meta tags if not specified
// in the config.
func Document(config DocumentConfig, bodyContent ...HTMLElement) *HTMLDocument {
	if config.Lang == "" {
		config.Lang = "en"
	}
	if config.Title == "" {
		config.Title = "Document"
	}
	if config.Charset == "" {
		config.Charset = "utf-8"
	}
	if config.Viewport == "" {
		config.Viewport = "width=device-width, initial-scale=1"
	}

	// Build head elements
	headElements := []HTMLElement{
		MetaCharset(config.Charset),
		MetaNameContent("viewport", config.Viewport),
		TitleEl(config.Title),
	}

	if config.Description != "" {
		headElements = append(headElements, MetaNameContent("description", config.Description))
	}
	if config.Keywords != "" {
		headElements = append(headElements, MetaNameContent("keywords", config.Keywords))
	}
	if config.Author != "" {
		headElements = append(headElements, MetaNameContent("author", config.Author))
	}

	// Add any extra head elements from config
	if len(config.HeadExtras) > 0 {
		headElements = append(headElements, config.HeadExtras...)
	}

	// Construct the root <html> element
	htmlRoot := Html(
		Head(headElements...),
		Body(bodyContent...),
	)
	htmlRoot.Attr("lang", config.Lang)

	return &HTMLDocument{rootElement: htmlRoot}
}

// Render converts the element to its HTML string representation.
// It handles both self-closing and container elements, attributes, content, and children.
// The output is properly formatted HTML that can be sent to browsers.
// Content and attribute values are HTML-escaped to prevent XSS, except for
// specific tags like `<script>` and `<style>` whose content must be raw.
func (e *Element) Render() string {
	var sb strings.Builder

	sb.WriteString("<")
	sb.WriteString(e.tag)

	for key, value := range e.attributes {
		sb.WriteString(fmt.Sprintf(` %s="%s"`, html.EscapeString(key), html.EscapeString(value)))
	}

	if e.selfClose {
		sb.WriteString(" />")
		return sb.String()
	}

	sb.WriteString(">")

	// Tags whose content should be rendered as raw text (not HTML-escaped)
	rawTextTags := map[string]bool{
		"style":  true,
		"script": true,
	}

	if e.content != "" {
		if _, isRaw := rawTextTags[e.tag]; isRaw {
			sb.WriteString(e.content) // Write raw content for `<style>`, `<script>`
		} else {
			sb.WriteString(html.EscapeString(e.content)) // Escape for other tags
		}
	}

	for _, child := range e.children {
		sb.WriteString(child.Render())
	}

	sb.WriteString("</")
	sb.WriteString(e.tag)
	sb.WriteString(">")

	return sb.String()
}

// Attr sets an attribute on the element and returns the element for method chaining.
func (e *Element) Attr(key, value string) *Element {
	if e.attributes == nil {
		e.attributes = make(map[string]string)
	}
	e.attributes[key] = value
	return e
}

// BoolAttr sets or removes a boolean attribute on the element.
// If present is true, the attribute is added (e.g., `<input disabled>`).
// If present is false, the attribute is removed if it exists.
func (e *Element) BoolAttr(key string, present bool) *Element {
	if e.attributes == nil && present {
		e.attributes = make(map[string]string)
	}
	if present {
		e.attributes[key] = key
	} else {
		if e.attributes != nil {
			delete(e.attributes, key)
		}
	}
	return e
}

// Class sets the class attribute on the element.
func (e *Element) Class(class string) *Element {
	return e.Attr("class", class)
}

// ID sets the id attribute on the element.
func (e *Element) ID(id string) *Element {
	return e.Attr("id", id)
}

// Style sets the style attribute on the element.
func (e *Element) Style(style string) *Element {
	return e.Attr("style", style)
}

// Text sets the text content of the element. This content is HTML-escaped during rendering.
func (e *Element) Text(text string) *Element {
	e.content = text
	return e
}

// Add appends child elements to this element.
func (e *Element) Add(children ...HTMLElement) *Element {
	e.children = append(e.children, children...)
	return e
}

// Render returns the raw text content, HTML-escaped.
func (t *textNode) Render() string {
	return html.EscapeString(t.text)
}

// Html creates an `<html>` element.
func Html(content ...HTMLElement) *Element {
	return &Element{tag: "html", children: content}
}

// Head creates a `<head>` element.
func Head(content ...HTMLElement) *Element {
	return &Element{tag: "head", children: content}
}

// Body creates a `<body>` element.
func Body(content ...HTMLElement) *Element {
	return &Element{tag: "body", children: content}
}

// TitleEl creates a `<title>` element with the specified text.
// Renamed from Title to TitleEl to avoid conflict with (*Element).Title method if it existed.
func TitleEl(titleText string) *Element {
	return &Element{tag: "title", content: titleText}
}

// Meta creates a generic `<meta>` element. It's self-closing.
func Meta() *Element {
	return &Element{tag: "meta", selfClose: true}
}

// MetaCharset creates a `<meta charset="...">` element.
func MetaCharset(charset string) *Element {
	return Meta().Attr("charset", charset)
}

// MetaNameContent creates a `<meta name="..." content="...">` element.
func MetaNameContent(name, contentVal string) *Element {
	return Meta().Attr("name", name).Attr("content", contentVal)
}

// MetaPropertyContent creates a `<meta property="..." content="...">` element.
func MetaPropertyContent(property, contentVal string) *Element {
	return Meta().Attr("property", property).Attr("content", contentVal)
}

// MetaViewport creates a `<meta name="viewport" content="...">` element.
func MetaViewport(contentVal string) *Element {
	return MetaNameContent("viewport", contentVal)
}

// Base creates a `<base>` element.
func Base(href string) *Element {
	return &Element{tag: "base", selfClose: true, attributes: map[string]string{"href": href}}
}

// LinkTag creates a generic `<link>` element. It's self-closing.
func LinkTag() *Element {
	return &Element{tag: "link", selfClose: true}
}

// StyleSheet creates a `<link rel="stylesheet">` element.
func StyleSheet(href string) *Element {
	return LinkTag().Attr("rel", "stylesheet").Attr("href", href)
}

// Favicon creates a `<link>` element for a favicon.
func Favicon(href string, rel ...string) *Element {
	link := LinkTag().Attr("href", href)
	if len(rel) > 0 && rel[0] != "" {
		link.Attr("rel", rel[0])
	} else {
		link.Attr("rel", "icon")
	}
	return link
}

// Preload creates a `<link rel="preload">` element.
func Preload(href string, asType string) *Element {
	return LinkTag().Attr("rel", "preload").Attr("href", href).Attr("as", asType)
}

// Script creates a `<script>` element for external JavaScript files.
func Script(src string) *Element {
	return &Element{tag: "script", attributes: map[string]string{"src": src}}
}

// InlineScript creates a `<script>` element with inline JavaScript content.
func InlineScript(scriptContent string) *Element {
	return &Element{tag: "script", content: scriptContent}
}

// StyleTag creates a `<style>` element for embedding CSS.
func StyleTag(cssContent string) *Element {
	return &Element{tag: "style", content: cssContent}
}

// NoScript creates a `<noscript>` element.
func NoScript(content ...HTMLElement) *Element {
	return &Element{tag: "noscript", children: content}
}

// Text creates a raw text node.
func Text(text string) HTMLElement {
	return &textNode{text: text}
}

// Div creates a `<div>` element.
func Div(content ...HTMLElement) *Element {
	return &Element{tag: "div", children: content}
}

// P creates a `<p>` paragraph element.
func P(content ...HTMLElement) *Element {
	return &Element{tag: "p", children: content}
}

// Span creates a `<span>` inline element.
func Span(content ...HTMLElement) *Element {
	return &Element{tag: "span", children: content}
}

// H1 creates an `<h1>` heading element.
func H1(content ...HTMLElement) *Element {
	return &Element{tag: "h1", children: content}
}

// H2 creates an `<h2>` heading element.
func H2(content ...HTMLElement) *Element {
	return &Element{tag: "h2", children: content}
}

// H3 creates an `<h3>` heading element.
func H3(content ...HTMLElement) *Element {
	return &Element{tag: "h3", children: content}
}

// H4 creates an `<h4>` heading element.
func H4(content ...HTMLElement) *Element {
	return &Element{tag: "h4", children: content}
}

// H5 creates an `<h5>` heading element.
func H5(content ...HTMLElement) *Element {
	return &Element{tag: "h5", children: content}
}

// H6 creates an `<h6>` heading element.
func H6(content ...HTMLElement) *Element {
	return &Element{tag: "h6", children: content}
}

// Br creates a self-closing `<br>` line break element.
func Br() *Element {
	return &Element{tag: "br", selfClose: true}
}

// Hr creates a self-closing `<hr>` horizontal rule element.
func Hr() *Element {
	return &Element{tag: "hr", selfClose: true}
}

// Header creates a `<header>` semantic element.
func Header(content ...HTMLElement) *Element {
	return &Element{tag: "header", children: content}
}

// Footer creates a `<footer>` semantic element.
func Footer(content ...HTMLElement) *Element {
	return &Element{tag: "footer", children: content}
}

// Main creates a `<main>` semantic element.
func Main(content ...HTMLElement) *Element {
	return &Element{tag: "main", children: content}
}

// Nav creates a `<nav>` semantic element.
func Nav(content ...HTMLElement) *Element {
	return &Element{tag: "nav", children: content}
}

// Section creates a `<section>` semantic element.
func Section(content ...HTMLElement) *Element {
	return &Element{tag: "section", children: content}
}

// Article creates an `<article>` semantic element.
func Article(content ...HTMLElement) *Element {
	return &Element{tag: "article", children: content}
}

// Aside creates an `<aside>` semantic element.
func Aside(content ...HTMLElement) *Element {
	return &Element{tag: "aside", children: content}
}

// Address creates an `<address>` semantic element.
func Address(content ...HTMLElement) *Element {
	return &Element{tag: "address", children: content}
}

// Figure creates a `<figure>` element.
func Figure(content ...HTMLElement) *Element {
	return &Element{tag: "figure", children: content}
}

// Figcaption creates a `<figcaption>` element.
func Figcaption(content ...HTMLElement) *Element {
	return &Element{tag: "figcaption", children: content}
}

// Details creates a `<details>` element.
func Details(content ...HTMLElement) *Element {
	return &Element{tag: "details", children: content}
}

// Summary creates a `<summary>` element.
func Summary(content ...HTMLElement) *Element {
	return &Element{tag: "summary", children: content}
}

// Blockquote creates a `<blockquote>` element.
func Blockquote(content ...HTMLElement) *Element {
	return &Element{tag: "blockquote", children: content}
}

// Q creates a `<q>` inline quotation element.
func Q(content ...HTMLElement) *Element {
	return &Element{tag: "q", children: content}
}

// Cite creates a `<cite>` element.
func Cite(content ...HTMLElement) *Element {
	return &Element{tag: "cite", children: content}
}

// Dfn creates a `<dfn>` definition element.
func Dfn(content ...HTMLElement) *Element {
	return &Element{tag: "dfn", children: content}
}

// Abbr creates an `<abbr>` abbreviation element.
func Abbr(content ...HTMLElement) *Element {
	return &Element{tag: "abbr", children: content}
}

// Mark creates a `<mark>` element.
func Mark(content ...HTMLElement) *Element {
	return &Element{tag: "mark", children: content}
}

// Small creates a `<small>` element.
func Small(content ...HTMLElement) *Element {
	return &Element{tag: "small", children: content}
}

// TimeEl creates a `<time>` element. Renamed to avoid potential conflicts.
func TimeEl(content ...HTMLElement) *Element {
	return &Element{tag: "time", children: content}
}

// Pre creates a `<pre>` element.
func Pre(content ...HTMLElement) *Element {
	return &Element{tag: "pre", children: content}
}

// Code creates a `<code>` element.
func Code(content ...HTMLElement) *Element {
	return &Element{tag: "code", children: content}
}

// Em creates an `<em>` emphasis element.
func Em(content ...HTMLElement) *Element {
	return &Element{tag: "em", children: content}
}

// Strong creates a `<strong>` element.
func Strong(content ...HTMLElement) *Element {
	return &Element{tag: "strong", children: content}
}

// Sub creates a `<sub>` subscript element.
func Sub(content ...HTMLElement) *Element {
	return &Element{tag: "sub", children: content}
}

// Sup creates a `<sup>` superscript element.
func Sup(content ...HTMLElement) *Element {
	return &Element{tag: "sup", children: content}
}

// I creates an `<i>` idiomatic text element.
func I(content ...HTMLElement) *Element {
	return &Element{tag: "i", children: content}
}

// B creates a `<b>` element for stylistically offset text.
func B(content ...HTMLElement) *Element {
	return &Element{tag: "b", children: content}
}

// U creates a `<u>` unarticulated annotation element.
func U(content ...HTMLElement) *Element {
	return &Element{tag: "u", children: content}
}

// VarEl creates a `<var>` variable element. Renamed to avoid keyword conflict.
func VarEl(content ...HTMLElement) *Element {
	return &Element{tag: "var", children: content}
}

// Samp creates a `<samp>` sample output element.
func Samp(content ...HTMLElement) *Element {
	return &Element{tag: "samp", children: content}
}

// Kbd creates a `<kbd>` keyboard input element.
func Kbd(content ...HTMLElement) *Element {
	return &Element{tag: "kbd", children: content}
}

// Wbr creates a `<wbr>` word break opportunity element. It's self-closing.
func Wbr() *Element {
	return &Element{tag: "wbr", selfClose: true}
}

// A creates an `<a>` anchor element.
func A(href string, content ...HTMLElement) *Element {
	return &Element{tag: "a", attributes: map[string]string{"href": href}, children: content}
}

// Link creates an `<a>` anchor element with href and text content.
func Link(href, textContent string) *Element {
	return A(href, Text(textContent))
}

// Img creates a self-closing `<img>` element.
func Img(src, alt string) *Element {
	return &Element{tag: "img", selfClose: true, attributes: map[string]string{"src": src, "alt": alt}}
}

// Image creates an `<img>` element (alias for Img).
func Image(src, alt string) *Element {
	return Img(src, alt)
}

// Ul creates a `<ul>` unordered list element.
func Ul(content ...HTMLElement) *Element {
	return &Element{tag: "ul", children: content}
}

// Ol creates an `<ol>` ordered list element.
func Ol(content ...HTMLElement) *Element {
	return &Element{tag: "ol", children: content}
}

// Li creates a `<li>` list item element.
func Li(content ...HTMLElement) *Element {
	return &Element{tag: "li", children: content}
}

// Table creates a `<table>` element.
func Table(content ...HTMLElement) *Element {
	return &Element{tag: "table", children: content}
}

// Thead creates a `<thead>` table header group element.
func Thead(content ...HTMLElement) *Element {
	return &Element{tag: "thead", children: content}
}

// Tbody creates a `<tbody>` table body group element.
func Tbody(content ...HTMLElement) *Element {
	return &Element{tag: "tbody", children: content}
}

// Tr creates a `<tr>` table row element.
func Tr(content ...HTMLElement) *Element {
	return &Element{tag: "tr", children: content}
}

// Td creates a `<td>` table data cell element.
func Td(content ...HTMLElement) *Element {
	return &Element{tag: "td", children: content}
}

// Th creates a `<th>` table header cell element.
func Th(content ...HTMLElement) *Element {
	return &Element{tag: "th", children: content}
}

// Caption creates a `<caption>` element for a table.
func Caption(content ...HTMLElement) *Element {
	return &Element{tag: "caption", children: content}
}

// Colgroup creates a `<colgroup>` element.
func Colgroup(content ...HTMLElement) *Element {
	return &Element{tag: "colgroup", children: content}
}

// Col creates a `<col>` element. It's self-closing.
func Col() *Element {
	return &Element{tag: "col", selfClose: true}
}

// Form creates a `<form>` element.
func Form(content ...HTMLElement) *Element {
	return &Element{tag: "form", children: content}
}

// Input creates a self-closing `<input>` element.
func Input(inputType string) *Element {
	return &Element{tag: "input", selfClose: true, attributes: map[string]string{"type": inputType}}
}

// Button creates a `<button>` element.
func Button(content ...HTMLElement) *Element {
	return &Element{tag: "button", children: content}
}

// Label creates a `<label>` element.
func Label(content ...HTMLElement) *Element {
	return &Element{tag: "label", children: content}
}

// Select creates a `<select>` dropdown element.
func Select(content ...HTMLElement) *Element {
	return &Element{tag: "select", children: content}
}

// Option creates an `<option>` element.
func Option(value string, content ...HTMLElement) *Element {
	e := &Element{tag: "option", children: content}
	if value != "" {
		e.Attr("value", value)
	}
	return e
}

// Textarea creates a `<textarea>` element.
func Textarea(content ...HTMLElement) *Element {
	return &Element{tag: "textarea", children: content}
}

// Fieldset creates a `<fieldset>` element.
func Fieldset(content ...HTMLElement) *Element {
	return &Element{tag: "fieldset", children: content}
}

// Legend creates a `<legend>` element.
func Legend(content ...HTMLElement) *Element {
	return &Element{tag: "legend", children: content}
}

// Optgroup creates an `<optgroup>` element.
func Optgroup(label string, content ...HTMLElement) *Element {
	return &Element{tag: "optgroup", attributes: map[string]string{"label": label}, children: content}
}

// Datalist creates a `<datalist>` element.
func Datalist(id string, content ...HTMLElement) *Element {
	element := &Element{
		tag:      "datalist",
		children: content,
	}

	return element.ID(id)
}

// OutputEl creates an `<output>` element. Renamed to avoid potential conflicts.
func OutputEl(content ...HTMLElement) *Element {
	return &Element{tag: "output", children: content}
}

// ProgressEl creates a `<progress>` element. Renamed to avoid potential conflicts.
func ProgressEl(content ...HTMLElement) *Element {
	return &Element{tag: "progress", children: content}
}

// MeterEl creates a `<meter>` element. Renamed to avoid potential conflicts.
func MeterEl(content ...HTMLElement) *Element {
	return &Element{tag: "meter", children: content}
}

// TextInput creates an `<input type="text">` field.
func TextInput(name string) *Element {
	return Input("text").Attr("name", name)
}

// EmailInput creates an `<input type="email">` field.
func EmailInput(name string) *Element {
	return Input("email").Attr("name", name)
}

// PasswordInput creates an `<input type="password">` field.
func PasswordInput(name string) *Element {
	return Input("password").Attr("name", name)
}

// CheckboxInput creates an `<input type="checkbox">`.
func CheckboxInput(name string) *Element {
	return Input("checkbox").Attr("name", name)
}

// RadioInput creates an `<input type="radio">`.
func RadioInput(name, value string) *Element {
	return Input("radio").Attr("name", name).Attr("value", value)
}

// NumberInput creates an `<input type="number">` field.
func NumberInput(name string) *Element {
	return Input("number").Attr("name", name)
}

// DateInput creates an `<input type="date">` field.
func DateInput(name string) *Element {
	return Input("date").Attr("name", name)
}

// FileInput creates an `<input type="file">` field.
func FileInput(name string) *Element {
	return Input("file").Attr("name", name)
}

// HiddenInput creates an `<input type="hidden">` field.
func HiddenInput(name string, value string) *Element {
	return Input("hidden").Attr("name", name).Attr("value", value)
}

// RangeInput creates an `<input type="range">` field.
func RangeInput(name string) *Element {
	return Input("range").Attr("name", name)
}

// SearchInput creates an `<input type="search">` field.
func SearchInput(name string) *Element {
	return Input("search").Attr("name", name)
}

// TelInput creates an `<input type="tel">` field.
func TelInput(name string) *Element {
	return Input("tel").Attr("name", name)
}

// UrlInput creates an `<input type="url">` field.
func UrlInput(name string) *Element {
	return Input("url").Attr("name", name)
}

// ColorInput creates an `<input type="color">` field.
func ColorInput(name string) *Element {
	return Input("color").Attr("name", name)
}

// DateTimeLocalInput creates an `<input type="datetime-local">` field.
func DateTimeLocalInput(name string) *Element {
	return Input("datetime-local").Attr("name", name)
}

// MonthInput creates an `<input type="month">` field.
func MonthInput(name string) *Element {
	return Input("month").Attr("name", name)
}

// WeekInput creates an `<input type="week">` field.
func WeekInput(name string) *Element {
	return Input("week").Attr("name", name)
}

// TimeInput creates an `<input type="time">` field.
func TimeInput(name string) *Element {
	return Input("time").Attr("name", name)
}

// SubmitButton creates a `<button type="submit">`.
func SubmitButton(text string) *Element {
	return Button(Text(text)).Attr("type", "submit")
}

// ResetButton creates a `<button type="reset">`.
func ResetButton(text string) *Element {
	return Button(Text(text)).Attr("type", "reset")
}

// ButtonInput creates an `<input type="button">`.
func ButtonInput(valueText string) *Element {
	return Input("button").Attr("value", valueText)
}

// Audio creates an `<audio>` element.
func Audio(content ...HTMLElement) *Element {
	return &Element{tag: "audio", children: content}
}

// Video creates a `<video>` element.
func Video(content ...HTMLElement) *Element {
	return &Element{tag: "video", children: content}
}

// Source creates a `<source>` element. It's self-closing.
func Source(src string, mediaType string) *Element {
	return &Element{tag: "source", selfClose: true, attributes: map[string]string{"src": src, "type": mediaType}}
}

// Track creates a `<track>` element. It's self-closing.
func Track(kind, src, srclang string) *Element {
	return &Element{tag: "track", selfClose: true, attributes: map[string]string{"kind": kind, "src": src, "srclang": srclang}}
}

// Iframe creates an `<iframe>` element.
func Iframe(src string) *Element {
	return &Element{tag: "iframe", attributes: map[string]string{"src": src}}
}

// EmbedEl creates an `<embed>` element. It's self-closing. Renamed to avoid keyword conflict.
func EmbedEl(src string, embedType string) *Element {
	return &Element{tag: "embed", selfClose: true, attributes: map[string]string{"src": src, "type": embedType}}
}

// ObjectEl creates an `<object>` element. Renamed to avoid keyword conflict.
func ObjectEl(content ...HTMLElement) *Element {
	return &Element{tag: "object", children: content}
}

// Param creates a `<param>` element. It's self-closing.
func Param(name, value string) *Element {
	return &Element{tag: "param", selfClose: true, attributes: map[string]string{"name": name, "value": value}}
}

// DialogEl creates a `<dialog>` element. Renamed to avoid potential conflicts.
func DialogEl(content ...HTMLElement) *Element {
	return &Element{tag: "dialog", children: content}
}
