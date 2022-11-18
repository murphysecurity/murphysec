package cmd

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/murphysecurity/murphysec/config"
	"github.com/spf13/cobra"
	"strings"
)

func authCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate CLI with murphysec",
	}
	c.AddCommand(authLoginCmd())
	c.AddCommand(authLogoutCmd())
	return c
}

var forceTokenOverwrite bool

func authLoginCmd() *cobra.Command {
	c := &cobra.Command{Use: "login [token]", Run: authLoginRun}
	c.Flags().BoolVarP(&forceTokenOverwrite, "force", "f", false, "force overwrite current token")
	c.Args = cobra.MaximumNArgs(1)
	return c
}

func authLogoutCmd() *cobra.Command {
	c := &cobra.Command{Use: "logout", Run: authLogoutRun}
	return c
}

func authLoginRun(cmd *cobra.Command, args []string) {
	var token string
	if len(args) == 1 {
		SLOG.Debug("Read token from args")
		token = args[0]
	} else {
		fmt.Println("Tips: You can generate a Personal Access Token here https://www.murphysec.com/control/set")
		var rs string
		e := survey.AskOne(&survey.Input{
			Message: "Input your token",
			Help:    "Tips: You can generate a Personal Access Token here https://www.murphysec.com/control/set",
		}, &rs, survey.WithValidator(survey.Required))
		if e != nil {
			fmt.Println("Token setup failed")
			SetGlobalExitCode(-1)
			return
		}
		token = rs
	}
	if strings.TrimSpace(token) == "" {
		fmt.Println("Token is invalid")
		SetGlobalExitCode(-1)
		return
	}
	if _, e := config.ReadTokenFile(context.TODO()); e == nil {
		var rs bool
		e := survey.AskOne(&survey.Confirm{Message: "Warning: You have a token, continue will overwrite it. That's OK?", Default: false}, &rs)
		if e != nil {
			fmt.Println("Cancelled.")
			SetGlobalExitCode(-1)
			return
		}
		if !rs {
			return
		}
	}
	e := config.WriteLocalTokenFile(context.TODO(), token)
	if e != nil {
		SLOG.Error("token setup failed: %s", e.Error())
		fmt.Println("Sorry, token setup failed")
		fmt.Println(e.Error())
		SetGlobalExitCode(-1)
	}
}

func authLogoutRun(cmd *cobra.Command, args []string) {
	e := config.RemoveTokenFile(context.TODO())
	if e != nil {
		SLOG.Errorf("auth logout failed: %s", e.Error())
		fmt.Println("Sorry, clear token failed.")
		fmt.Println(e.Error())
		SetGlobalExitCode(-1)
	}
}
