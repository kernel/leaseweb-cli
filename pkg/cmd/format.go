package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/itchyny/json2yaml"
	"github.com/tidwall/gjson"
	"github.com/tidwall/pretty"
	"golang.org/x/term"
)

var OutputFormats = []string{"auto", "json", "jsonline", "pretty", "raw", "yaml"}

func isTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func shouldUseColors(w io.Writer) bool {
	force, ok := os.LookupEnv("FORCE_COLOR")
	if ok {
		if force == "1" {
			return true
		}
		if force == "0" {
			return false
		}
	}
	return isTerminal(w)
}

// ShowJSON displays a JSON string in the requested format.
func ShowJSON(out *os.File, raw string, format string, transform string) error {
	res := gjson.Parse(raw)
	if transform != "" {
		transformed := res.Get(transform)
		if transformed.Exists() {
			res = transformed
		}
	}
	switch strings.ToLower(format) {
	case "auto":
		ShowDetail(out, res)
		return nil
	case "json":
		prettyJSON := pretty.Pretty([]byte(res.Raw))
		if shouldUseColors(out) {
			_, err := out.Write(pretty.Color(prettyJSON, pretty.TerminalStyle))
			return err
		}
		_, err := out.Write(prettyJSON)
		return err
	case "pretty":
		prettyJSON := pretty.Pretty([]byte(res.Raw))
		_, err := out.Write(prettyJSON)
		return err
	case "jsonline":
		_, err := out.Write([]byte(res.Raw + "\n"))
		return err
	case "raw":
		_, err := out.Write([]byte(res.Raw + "\n"))
		return err
	case "yaml":
		input := strings.NewReader(res.Raw)
		var yamlOut strings.Builder
		if err := json2yaml.Convert(&yamlOut, input); err != nil {
			return err
		}
		_, err := out.Write([]byte(yamlOut.String()))
		return err
	default:
		return fmt.Errorf("invalid format: %s, valid formats are: %s", format, strings.Join(OutputFormats, ", "))
	}
}

// ShowResult displays a gjson.Result in the requested format.
func ShowResult(out *os.File, res gjson.Result, format string, transform string) error {
	return ShowJSON(out, res.Raw, format, transform)
}
