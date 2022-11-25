package auth

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/config"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/spf13/cobra"
)

func authLogoutCmd() *cobra.Command {
	c := &cobra.Command{Use: "logout", Run: authLogoutRun}
	return c
}

func authLogoutRun(cmd *cobra.Command, args []string) {
	e := config.RemoveTokenFile(context.TODO())
	if e != nil {
		fmt.Println("Sorry, clear token failed.")
		fmt.Println(e.Error())
		exitcode.Set(-1)
	}
}
