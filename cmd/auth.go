package cmd

import (
	"github.com/spf13/cobra"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/logger"
)

func authCmd() *cobra.Command {
	c := &cobra.Command{Use: "auth"}
	c.AddCommand(authLoginCmd())
	c.AddCommand(authLogoutCmd())
	c.AddCommand(authCheckCmd())
	return c
}

func authCheckCmd() *cobra.Command {
	c := &cobra.Command{Use: "check", RunE: authCheckRun}
	return c
}

func authLoginCmd() *cobra.Command {
	c := &cobra.Command{Use: "login", RunE: authLogoutRun}
	return c
}

func authLogoutCmd() *cobra.Command {
	c := &cobra.Command{Use: "logout", RunE: authLogoutRun}
	return c
}

func authCheckRun(cmd *cobra.Command, args []string) error {
	return nil
}

func authLogoutRun(cmd *cobra.Command, args []string) error {
	e := conf.RemoveToken()
	if e == conf.TokenFileNotFound {
		logger.Warn.Println("Token file is not exists")
		return nil
	}
	return e
}
