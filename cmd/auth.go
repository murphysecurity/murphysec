package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"
	"murphysec-cli-simple/api"
	"murphysec-cli-simple/conf"
	"murphysec-cli-simple/util/output"
	"os"
	"strings"
)

func authCmd() *cobra.Command {
	login := &cobra.Command{
		Use:   "login",
		Short: "setup API token",
		Run:   setupToken,
	}
	logout := &cobra.Command{
		Use:   "logout",
		Short: "clear API token",
		Run:   clearToken,
	}
	checkToken := &cobra.Command{
		Use:   "check",
		Short: "check API token",
		Run:   checkToken,
	}
	c := &cobra.Command{
		Use:   "auth",
		Short: "manage the API token",
	}
	c.AddCommand(login)
	c.AddCommand(logout)
	c.AddCommand(checkToken)
	return c
}

func checkToken(cmd *cobra.Command, args []string) {
	valid, err := api.CheckAPIToken(conf.APIToken())
	if err != nil {
		output.Error(err.Error())
		os.Exit(1)
		return
	}
	if valid {
		output.Info("OK, the token is valid!")
	} else {
		output.Error("Sorry, the token is invalid!")
		os.Exit(1)
	}
}

func clearToken(cmd *cobra.Command, args []string) {
	if e := conf.RemoveToken(); e != nil {
		if e == conf.TokenFileNotFound {
			output.Info("Token not set.")
			return
		}
		output.Error(e.Error())
		os.Exit(1)
		return
	}
	output.Info("Token cleared!")
}
func setupToken(cmd *cobra.Command, args []string) {
	// if token is set, alert user
	if len(strings.TrimSpace(conf.APIToken())) != 0 {
		var rs bool
		err := survey.AskOne(&survey.Confirm{
			Message: "You're logged in, Do you want to re-authenticate?",
			Default: false,
		}, &rs)
		if err != nil {
			output.Error(err.Error())
			os.Exit(1)
			return
		}
		if !rs {
			return
		}
	}
	fmt.Println(heredoc.Doc(`
		Tip:  you can generate a Personal Access Token here https://www.murphysec.com/control/set
	`))
	rs := ""
	err := survey.AskOne(&survey.Input{Message: "Paste an authentication token"}, &rs, survey.WithValidator(survey.Required))
	if err != nil {
		output.Error(err.Error())
		os.Exit(1)
		return
	}
	if e := conf.StoreToken(rs); e != nil {
		output.Error(e.Error())
		os.Exit(1)
		return
	}
	output.Info("Token setup succeed!")
}
