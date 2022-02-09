package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/logger"
	"strings"
)

func authCmd() *cobra.Command {
	c := &cobra.Command{Use: "auth"}
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
		logger.Debug.Println("Read token from args")
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
	if _, e := conf.ReadTokenFile(); e == nil {
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
	e := conf.StoreToken(token)
	if e != nil {
		logger.Err.Println("token setup failed")
		logger.Err.Println(e.Error())
		fmt.Println("Sorry, token setup failed")
		fmt.Println(e.Error())
		SetGlobalExitCode(-1)
	}
}

func authLogoutRun(cmd *cobra.Command, args []string) {
	e := conf.RemoveToken()
	if e == conf.TokenFileNotFound {
		logger.Warn.Println("Token file is not exists")
		SetGlobalExitCode(0)
		return
	}
	if e != nil {
		logger.Err.Println("auth logout failed.", e.Error())
		fmt.Println("Sorry, clear token failed.")
		fmt.Println(e.Error())
		SetGlobalExitCode(-1)
	}
}
