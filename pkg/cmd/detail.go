package cmd

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/tidwall/gjson"
)

// camelToTitle converts camelCase or PascalCase to Title Case.
// e.g. "serialNumber" → "Serial Number", "networkInterfaces" → "Network Interfaces"
// Single well-known abbreviations are kept uppercase.
func camelToTitle(s string) string {
	if s == "" {
		return s
	}
	if upper, ok := knownAbbreviations[strings.ToLower(s)]; ok {
		return upper
	}
	var b strings.Builder
	runes := []rune(s)
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			if unicode.IsLower(prev) {
				b.WriteRune(' ')
			} else if i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				b.WriteRune(' ')
			}
		}
		if i == 0 {
			b.WriteRune(unicode.ToUpper(r))
		} else {
			b.WriteRune(r)
		}
	}
	result := b.String()
	words := strings.Split(result, " ")
	for i, w := range words {
		if repl, ok := wordAbbreviations[w]; ok {
			words[i] = repl
		}
	}
	return strings.Join(words, " ")
}

var knownAbbreviations = map[string]string{
	"id": "ID", "ip": "IP", "cpu": "CPU", "ram": "RAM",
	"mac": "MAC", "sla": "SLA", "os": "OS", "hdd": "HDD",
	"ssd": "SSD", "url": "URL", "dns": "DNS", "ssl": "SSL",
	"uuid": "UUID", "ipmi": "IPMI",
}

var wordAbbreviations = map[string]string{
	"Id": "ID", "Ip": "IP", "Cpu": "CPU", "Ram": "RAM",
	"Mac": "MAC", "Sla": "SLA", "Os": "OS", "Hdd": "HDD",
	"Ssd": "SSD", "Url": "URL", "Dns": "DNS", "Ssl": "SSL",
	"Uuid": "UUID", "Ipmi": "IPMI", "Ipv4": "IPv4", "Ipv6": "IPv6",
	"Ddos": "DDoS", "Pci": "PCI",
}

func formatValue(v gjson.Result) string {
	switch v.Type {
	case gjson.Null:
		return "—"
	case gjson.True:
		return "yes"
	case gjson.False:
		return "no"
	case gjson.String:
		return v.Str
	case gjson.Number:
		f := v.Float()
		if f == float64(int64(f)) {
			return fmt.Sprintf("%d", int64(f))
		}
		return fmt.Sprintf("%.2f", f)
	default:
		return v.Raw
	}
}

// ShowDetail renders a JSON object as a grouped text display.
func ShowDetail(w io.Writer, data gjson.Result) {
	colors := shouldUseColors(w)
	if !data.IsObject() {
		fmt.Fprintln(w, data.Raw)
		return
	}

	var scalars []field
	var sections []field

	data.ForEach(func(key, val gjson.Result) bool {
		k := key.String()
		if val.IsObject() || val.IsArray() {
			sections = append(sections, field{k, val})
		} else {
			scalars = append(scalars, field{k, val})
		}
		return true
	})

	if len(scalars) > 0 {
		renderKVBlock(w, "", scalars, colors)
	}

	for _, sec := range sections {
		renderSection(w, sec.key, sec.val, 0, colors)
	}
}

func renderKVBlock(w io.Writer, header string, fields []field, colors bool) {
	indent := "  "
	if header != "" {
		fmt.Fprintf(w, "\n%s\n", styledHeader(header, colors))
	} else {
		indent = ""
	}

	maxKeyLen := 0
	for _, f := range fields {
		label := camelToTitle(f.key)
		if len(label) > maxKeyLen {
			maxKeyLen = len(label)
		}
	}

	for _, f := range fields {
		label := camelToTitle(f.key)
		fmt.Fprintf(w, "%s%-*s  %s\n", indent, maxKeyLen, label, formatValue(f.val))
	}
}

type field struct {
	key string
	val gjson.Result
}

func renderSection(w io.Writer, name string, val gjson.Result, depth int, colors bool) {
	indent := strings.Repeat("  ", depth)

	if val.IsArray() {
		arr := val.Array()
		if len(arr) == 0 {
			fmt.Fprintf(w, "\n%s%s\n", indent, styledHeader(name, colors))
			fmt.Fprintf(w, "%s  (none)\n", indent)
			return
		}
		if arr[0].IsObject() {
			renderArrayOfObjects(w, name, arr, depth, colors)
		} else {
			fmt.Fprintf(w, "\n%s%s\n", indent, styledHeader(name, colors))
			for _, item := range arr {
				fmt.Fprintf(w, "%s  • %s\n", indent, formatValue(item))
			}
		}
		return
	}

	if !val.IsObject() {
		fmt.Fprintf(w, "\n%s%s\n", indent, styledHeader(name, colors))
		fmt.Fprintf(w, "%s  %s\n", indent, formatValue(val))
		return
	}

	var scalars []field
	var nested []field

	val.ForEach(func(key, v gjson.Result) bool {
		k := key.String()
		if v.IsObject() || v.IsArray() {
			nested = append(nested, field{k, v})
		} else {
			scalars = append(scalars, field{k, v})
		}
		return true
	})

	fmt.Fprintf(w, "\n%s%s\n", indent, styledHeader(name, colors))

	if len(scalars) > 0 {
		maxKeyLen := 0
		for _, f := range scalars {
			label := camelToTitle(f.key)
			if len(label) > maxKeyLen {
				maxKeyLen = len(label)
			}
		}
		for _, f := range scalars {
			label := camelToTitle(f.key)
			fmt.Fprintf(w, "%s  %-*s  %s\n", indent, maxKeyLen, label, formatValue(f.val))
		}
	}

	for _, n := range nested {
		renderSection(w, n.key, n.val, depth+1, colors)
	}
}

func renderArrayOfObjects(w io.Writer, name string, arr []gjson.Result, depth int, colors bool) {
	indent := strings.Repeat("  ", depth)
	fmt.Fprintf(w, "\n%s%s\n", indent, styledHeader(name, colors))

	// Collect all keys to determine if this is a simple table
	allKeys := make([]string, 0)
	keySet := make(map[string]bool)
	hasNested := false
	for _, item := range arr {
		item.ForEach(func(key, val gjson.Result) bool {
			k := key.String()
			if !keySet[k] {
				keySet[k] = true
				allKeys = append(allKeys, k)
			}
			if val.IsObject() || val.IsArray() {
				hasNested = true
			}
			return true
		})
	}

	// If items are flat and few columns, render as inline table
	if !hasNested && len(allKeys) <= 6 {
		headers := make([]string, len(allKeys))
		widths := make([]int, len(allKeys))
		for i, k := range allKeys {
			headers[i] = camelToTitle(k)
			widths[i] = len(headers[i])
		}
		rows := make([][]string, len(arr))
		for i, item := range arr {
			row := make([]string, len(allKeys))
			for j, k := range allKeys {
				row[j] = formatValue(item.Get(k))
				if len(row[j]) > widths[j] {
					widths[j] = len(row[j])
				}
			}
			rows[i] = row
		}
		// Print header
		for i, h := range headers {
			if i < len(headers)-1 {
				fmt.Fprintf(w, "%s  %-*s", indent, widths[i]+1, h)
			} else {
				fmt.Fprintf(w, "%s\n", h)
			}
		}
		for _, row := range rows {
			for i, cell := range row {
				if i < len(row)-1 {
					fmt.Fprintf(w, "%s  %-*s", indent, widths[i]+1, cell)
				} else {
					fmt.Fprintf(w, "%s\n", cell)
				}
			}
		}
		return
	}

	// Complex objects: render each as a sub-block
	for i, item := range arr {
		if i > 0 {
			fmt.Fprintf(w, "%s  ---\n", indent)
		}
		var fields []field
		var nested []field
		item.ForEach(func(key, val gjson.Result) bool {
			k := key.String()
			if val.IsObject() || val.IsArray() {
				nested = append(nested, field{k, val})
			} else {
				fields = append(fields, field{k, val})
			}
			return true
		})
		maxKeyLen := 0
		for _, f := range fields {
			label := camelToTitle(f.key)
			if len(label) > maxKeyLen {
				maxKeyLen = len(label)
			}
		}
		for _, f := range fields {
			label := camelToTitle(f.key)
			fmt.Fprintf(w, "%s  %-*s  %s\n", indent, maxKeyLen, label, formatValue(f.val))
		}
		for _, n := range nested {
			renderSection(w, n.key, n.val, depth+1, colors)
		}
	}
}

func styledHeader(name string, colors bool) string {
	title := camelToTitle(name)
	if colors {
		return fmt.Sprintf("\033[1m%s\033[0m", title)
	}
	return title
}
