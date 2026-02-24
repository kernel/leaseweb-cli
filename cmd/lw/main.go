package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kernel/leaseweb-cli/pkg/cmd"
)

func main() {
	app := cmd.Command
	if err := app.Run(context.Background(), os.Args); err != nil {
		var apiErr *cmd.APIError
		if ok := isAPIError(err, &apiErr); ok {
			fmt.Fprintf(os.Stderr, "%s %q: %d %s\n", apiErr.Method, apiErr.URL, apiErr.StatusCode, apiErr.Status)
			if apiErr.Body != "" {
				showErr := cmd.ShowJSON(os.Stdout, apiErr.Body, app.String("format"), app.String("transform"))
				if showErr != nil {
					fmt.Fprintf(os.Stderr, "%s\n", apiErr.Body)
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
		os.Exit(1)
	}
}

func isAPIError(err error, target **cmd.APIError) bool {
	if apiErr, ok := err.(*cmd.APIError); ok {
		*target = apiErr
		return true
	}
	return false
}
