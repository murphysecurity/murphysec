package maven

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/vifraa/gopom"
	"murphysec-cli-simple/utils/must"
	"net/url"
	"testing"
)

func Test_ResolveProperty(t *testing.T) {
	m := map[string]string{
		"a": "1",
		"b": "2${a}",
		"c": "foo${b}${d}",
	}
	assert.Equal(t, _resolveProperty(m, nil, "c"), "foo21${d}")
	t.Log(_resolveProperty(m, nil, "c"))
}

func TestAaa(t *testing.T) {
	p, _ := gopom.Parse("C:\\Users\\iseki\\Desktop\\新建文件夹\\HundredBai_camellia_master\\camellia-redis\\pom.xml")
	builder := NewPomBuilder(p)
	builder.Path = "C:\\Users\\iseki\\Desktop\\新建文件夹\\HundredBai_camellia_master\\camellia-redis"
	resolver := NewResolver()
	resolver.repos = append(resolver.repos, &LocalRepo{"C:\\Users\\iseki\\.m2\\repository"}, NewHttpRepo(*must.Url(url.Parse("https://mirrors.cloud.tencent.com/nexus/repository/maven-public/"))))
	pf := resolver.Resolve(builder, nil)
	fmt.Println(pf.propertyMap)
	fmt.Println(pf.dependencyManagement)
	fmt.Println(pf.dependencies)
	fmt.Println(pf.path)
	fmt.Println(pf.coordinate)
	fmt.Println(pf.parentPom.path)
	fmt.Println(pf.parentPom.coordinate)
	analyzer := NewDepTreeAnalyzer(resolver)
	fmt.Println(analyzer.Resolve(pf).Tree(pf.coordinate))
}
