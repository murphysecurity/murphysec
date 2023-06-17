package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/murphysecurity/murphysec/config"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/spf13/cobra"
)

func authLoginCmd() *cobra.Command {
	c := &cobra.Command{Use: "login [token]", Run: authLoginRun}
	var forceTokenOverwrite bool
	c.Flags().BoolVarP(&forceTokenOverwrite, "force", "f", false, "force overwrite current token")
	c.Args = cobra.MaximumNArgs(1)
	return c
}

func authLoginRun(cmd *cobra.Command, args []string) {
	var token string
	if len(args) == 1 {
		token = args[0]
	} else {
		var rs string
		e := survey.AskOne(&survey.Input{
			Message: "Input your token",
		}, &rs, survey.WithValidator(survey.Required))
		if e != nil {
			fmt.Println("Token setup failed")
			exitcode.Set(-1)
			return
		}
		token = rs
	}
	if strings.TrimSpace(token) == "" {
		fmt.Println("Token is invalid")
		exitcode.Set(-1)
		return
	}
	if _, e := config.ReadTokenFile(context.TODO()); e == nil {
		var rs bool
		e := survey.AskOne(&survey.Confirm{Message: "Warning: You have a token, continue will overwrite it. That's OK?", Default: false}, &rs)
		if e != nil {
			fmt.Println("Cancelled.")
			exitcode.Set(-1)
			return
		}
		if !rs {
			return
		}
	}
	e := config.WriteLocalTokenFile(context.TODO(), token)
	if e != nil {
		fmt.Println("Sorry, token setup failed")
		fmt.Println(e.Error())
		exitcode.Set(-1)
	}
}
