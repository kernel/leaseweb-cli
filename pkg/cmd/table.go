package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"golang.org/x/term"
)

const columnGap = 2

type TableWriter struct {
	w       io.Writer
	headers []string
	widths  []int
	rows    [][]string

	// TruncOrder specifies column indices in truncation priority order.
	// The first index is truncated first when the table is too wide.
	// Columns not listed are never truncated.
	TruncOrder []int
}

func NewTableWriter(w io.Writer, headers ...string) *TableWriter {
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	return &TableWriter{
		w:       w,
		headers: headers,
		widths:  widths,
	}
}

func (t *TableWriter) AddRow(cells ...string) {
	row := make([]string, len(t.headers))
	for i := range row {
		if i < len(cells) {
			row[i] = cells[i]
		}
		if len(row[i]) > t.widths[i] {
			t.widths[i] = len(row[i])
		}
	}
	t.rows = append(t.rows, row)
}

func getTerminalWidth() int {
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
		return w
	}
	if w, _, err := term.GetSize(int(os.Stderr.Fd())); err == nil && w > 0 {
		return w
	}
	if cols := os.Getenv("COLUMNS"); cols != "" {
		if w, err := strconv.Atoi(cols); err == nil && w > 0 {
			return w
		}
	}
	return 120
}

func (t *TableWriter) renderWidths() []int {
	n := len(t.headers)
	widths := make([]int, n)
	copy(widths, t.widths)

	termWidth := getTerminalWidth()

	total := func() int {
		s := 0
		for _, w := range widths {
			s += w
		}
		s += columnGap * (n - 1)
		return s
	}

	if total() <= termWidth {
		return widths
	}

	for _, col := range t.TruncOrder {
		if col < 0 || col >= n {
			continue
		}
		excess := total() - termWidth
		if excess <= 0 {
			break
		}
		minW := len(t.headers[col])
		if minW < 5 {
			minW = 5
		}
		canShrink := widths[col] - minW
		if canShrink <= 0 {
			continue
		}
		shrink := excess
		if shrink > canShrink {
			shrink = canShrink
		}
		widths[col] -= shrink
	}

	return widths
}

func (t *TableWriter) Render() {
	widths := t.renderWidths()
	last := len(t.headers) - 1

	for i, h := range t.headers {
		cell := truncateCell(h, widths[i])
		if i < last {
			fmt.Fprintf(t.w, "%-*s", widths[i]+columnGap, cell)
		} else {
			fmt.Fprint(t.w, cell)
		}
	}
	fmt.Fprintln(t.w)

	for _, row := range t.rows {
		for i, cell := range row {
			cell = truncateCell(cell, widths[i])
			if i < last {
				fmt.Fprintf(t.w, "%-*s", widths[i]+columnGap, cell)
			} else {
				fmt.Fprint(t.w, cell)
			}
		}
		fmt.Fprintln(t.w)
	}
}

func truncateCell(s string, maxWidth int) string {
	if len(s) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return s[:maxWidth]
	}
	return s[:maxWidth-3] + "..."
}

func FormatTimeAgo(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		mins := int(d.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", mins)
	case d < 24*time.Hour:
		hours := int(d.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}
}
