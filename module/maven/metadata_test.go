package maven

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parsePomVersionMeta(t *testing.T) {
	var data = `
<?xml version="1.0" encoding="UTF-8"?>
<metadata modelVersion="1.1.0">
  <groupId>space.iseki.envproxyselector</groupId>
  <artifactId>envproxyselector</artifactId>
  <version>0.1.0-SNAPSHOT</version>
  <versioning>
    <snapshot>
      <timestamp>20221222.123814</timestamp>
      <buildNumber>2</buildNumber>
    </snapshot>
    <lastUpdated>20221222123814</lastUpdated>
    <snapshotVersions>
      <snapshotVersion>
        <extension>jar</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
      <snapshotVersion>
        <classifier>sources</classifier>
        <extension>jar</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
      <snapshotVersion>
        <extension>module.asc</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
      <snapshotVersion>
        <extension>pom</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
      <snapshotVersion>
        <classifier>javadoc</classifier>
        <extension>jar</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
      <snapshotVersion>
        <classifier>sources</classifier>
        <extension>jar.asc</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
      <snapshotVersion>
        <extension>module</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
      <snapshotVersion>
        <extension>jar.asc</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
      <snapshotVersion>
        <classifier>javadoc</classifier>
        <extension>jar.asc</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
      <snapshotVersion>
        <extension>pom.asc</extension>
        <value>0.1.0-20221222.123814-2</value>
        <updated>20221222123814</updated>
      </snapshotVersion>
    </snapshotVersions>
  </versioning>
</metadata>
`
	meta, e := parsePomVersionMeta(bytes.NewReader([]byte(data)))
	assert.NoError(t, e)
	t.Log(len(meta.SnapshotVersions))
	assert.Positive(t, len(meta.SnapshotVersions))
	t.Log(meta.getPomVersionSnapshotSuffix())
	assert.Equal(t, "0.1.0-20221222.123814-2", meta.getPomVersionSnapshotSuffix())
}
