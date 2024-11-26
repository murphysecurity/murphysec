package common

import (
	"fmt"
	"github.com/spf13/pflag"
	"strconv"
	"strings"
)

type WebhookModeFlag struct {
	simple bool
}

func (w *WebhookModeFlag) IsSimple() bool {
	return w.simple
}

func (w *WebhookModeFlag) String() string {
	if w.simple {
		return "simple"
	}
	return "full"
}

func (w *WebhookModeFlag) Set(s string) error {
	s = strings.ToLower(s)
	if s == "simple" {
		w.simple = true
	} else if s == "" || s == "full" {
		w.simple = false
	} else {
		return fmt.Errorf("invalid webhook mode: %s", strconv.Quote(s))
	}
	return nil
}

func (w *WebhookModeFlag) Type() string {
	return "WebhookModeFlag"
}

var _ pflag.Value = (*WebhookModeFlag)(nil)
