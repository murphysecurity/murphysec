package api

import (
	"github.com/murphysecurity/murphysec/utils/must"
	"net/url"
	"testing"
)

func TestJoinURL(t *testing.T) {
	t.Log(joinURL(must.A(url.Parse("https://iseki.space/")), "/aa/b"))
}
