package internalcmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/murphysecurity/murphysec/cpphasher"
	"github.com/murphysecurity/murphysec/infra/exitcode"
	"github.com/murphysecurity/murphysec/infra/logctx"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"path/filepath"
)

func cppFileHashCmd() *cobra.Command {
	c := &cobra.Command{
		Use:  "cpphash",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			logger := must.M1(zap.NewProduction())
			ctx := logctx.With(context.TODO(), logger)
			dir := must.M1(filepath.Abs(args[0]))
			hashes, e := cpphasher.MD5HashingCppFiles(ctx, dir)
			if e != nil {
				logger.Error("error", zap.Error(e))
				exitcode.Set(1)
				return
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
