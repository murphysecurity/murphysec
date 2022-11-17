package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"path/filepath"
)

func internalCmd() *cobra.Command {
	var c = &cobra.Command{Use: "internal", Hidden: true}
	c.AddCommand(cppFileHashCmd())
	return c
}

func cppFileHashCmd() *cobra.Command {
	c := &cobra.Command{
		Use:  "cpphash",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := must.M1(zap.NewProduction())
			ctx := logctx.With(context.TODO(), logger)
			dir := must.M1(filepath.Abs(args[0]))
			hashes, e := inspector.MD5HashingCppFiles(ctx, dir)
			if e != nil {
				logger.Error("error", zap.Error(e))
			}
			var s = make([]string, 0, len(hashes))
			for _, hash := range hashes {
				s = append(s, hex.EncodeToString(hash[:]))
			}
			fmt.Println(string(must.M1(json.Marshal(s))))
		},
	}
	return c
}
