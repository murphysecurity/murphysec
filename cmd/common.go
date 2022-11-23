package cmd

import (
	"context"
	"github.com/murphysecurity/murphysec/api"
	"github.com/murphysecurity/murphysec/config"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/pkg/errors"
)

func initAPI(ctx context.Context) error {
	if !loggerInitialized {
		panic("logger not initialized")
	}
	token, e := getToken(ctx)
	if errors.Is(e, config.ErrNoToken) {
		displayTokenNotSet(ctx)
		return e
	}
	if e != nil {
		displayGetTokenErr(ctx, e)
		return e
	}
	var cf = &api.Config{
		Logger:             logctx.Use(ctx),
		EnableNetworkDebug: enableNetworkLog,
		Token:              token,
		AllowInsecure:      allowInsecure,
	}
	if env.ServerURLOverride != "" {
		cf.ServerURL = env.ServerURLOverride
	}
	if cliServerAddressOverride != "" {
		cf.ServerURL = cliServerAddressOverride
	}
	e = api.InitDefaultClient(cf)
	if e != nil {
		displayInitializeFailed(ctx, e)
		return e
	}
	return nil
}

func getToken(ctx context.Context) (string, error) {
	if cliTokenOverride != "" {
		return cliTokenOverride, nil
	}
	if env.APITokenOverride != "" {
		return env.APITokenOverride, nil
	}
	return config.ReadTokenFile(ctx)
}
