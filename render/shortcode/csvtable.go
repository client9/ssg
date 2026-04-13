package shortcode

import (
	"strings"

	"github.com/olekukonko/tablewriter"
)

// CSVTable is a HandlerFunc that renders the body as a CSV-formatted table.
// The first row of the CSV is used as the header. Macro expressions inside
// cells are parsed before CSV splitting, so macro bodies may safely contain
// commas and are expanded when each cell is rendered.
//
// Example:
//
//	$csvtable{
//	Name,Score,Note
//	Alice,100,$b{perfect}
//	Bob,95,$link[https://example.com]{details}
//	}
func CSVTable(ctx *Context, args []string, body string) string {
	records := parseCSVWithMacros(strings.TrimSpace(body))
	if len(records) == 0 {
		return body
	}

	var sb strings.Builder
	table := tablewriter.NewTable(&sb)

	header := make([]any, len(records[0]))
	for i, cell := range records[0] {
		header[i] = ctx.Render(cell)
	}
	table.Header(header...)

	for _, row := range records[1:] {
		rendered := make([]string, len(row))
		for i, cell := range row {
			rendered[i] = ctx.Render(cell)
		}
		table.Append(rendered)
	}
	table.Render()

	return sb.String()
}

// parseQuotedField reads a CSV quoted field starting just after the opening
// quote. Handles "" escapes and macro tokens. Returns (content, endPos) where
// endPos is the index after the closing quote (or end of string if unterminated).
func parseQuotedField(s string, i int) (string, int) {
	var sb strings.Builder
	for i < len(s) {
		switch {
		case s[i] == '"':
			if i+1 < len(s) && s[i+1] == '"' {
				sb.WriteByte('"') // escaped quote
				i += 2
			} else {
				return sb.String(), i + 1 // skip closing quote
			}
		case s[i] == '$':
			_, _, _, end, ok := parseMacro(s, i)
			if ok {
				sb.WriteString(s[i:end])
				i = end
			} else {
				sb.WriteByte('$')
				i++
			}
		default:
			sb.WriteByte(s[i])
			i++
		}
	}
	return sb.String(), i // unterminated quote: consume to end
}

// parseCSVWithMacros splits s into rows and fields like a CSV parser, but
// treats macro expressions ($cmd, $cmd[...], $cmd{...}) as atomic tokens so
// that commas and braces inside them are never interpreted as CSV delimiters.
//
// Quoting follows standard CSV rules ("" inside a quoted field is a literal
// quote). Macros are also recognised inside quoted fields.
func parseCSVWithMacros(s string) [][]string {
	var rows [][]string
	var row []string
	var field strings.Builder
	i := 0

	commitField := func() {
		row = append(row, field.String())
		field.Reset()
	}
	commitRow := func() {
		commitField()
		rows = append(rows, row)
		row = nil
	}

	for i < len(s) {
		switch s[i] {
		case '$':
			// Consume the entire macro token as a single unit so that
			// commas/braces inside args or body don't act as CSV delimiters.
			_, _, _, end, ok := parseMacro(s, i)
			if ok {
				field.WriteString(s[i:end])
				i = end
			} else {
				field.WriteByte('$')
				i++
			}

		case '"':
			// Standard CSV quoted field; also honour macros inside it.
			content, end := parseQuotedField(s, i+1)
			field.WriteString(content)
			i = end

		case ',':
			commitField()
			i++

		case '\n':
			commitRow()
			i++

		case '\r':
			i++ // skip; \n will commitRow

		default:
			field.WriteByte(s[i])
			i++
		}
	}

	// Flush any trailing content.
	if field.Len() > 0 || len(row) > 0 {
		commitRow()
	}

	return rows
}
