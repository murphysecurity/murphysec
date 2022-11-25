package auth

import (
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate CLI with murphysec",
	}
	c.AddCommand(authLoginCmd())
	c.AddCommand(authLogoutCmd())
	return c
}
