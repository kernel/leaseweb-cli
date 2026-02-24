package cmd

import (
	"fmt"

	"github.com/urfave/cli/v3"
)

var PaginationFlags = []cli.Flag{
	&cli.IntFlag{
		Name:  "limit",
		Usage: "Maximum number of results to return",
		Value: 20,
	},
	&cli.IntFlag{
		Name:  "offset",
		Usage: "Number of results to skip",
		Value: 0,
	},
}

func PaginationQuery(cmd *cli.Command) string {
	limit := cmd.Int("limit")
	offset := cmd.Int("offset")
	return fmt.Sprintf("limit=%d&offset=%d", limit, offset)
}
