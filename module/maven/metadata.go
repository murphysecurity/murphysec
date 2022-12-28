package maven

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/ianaindex"
	"io"
)

type Metadata struct {
	ModelVersion     string            `xml:"model_version,attr"`
	GroupId          string            `xml:"group_id"`
	ArtifactId       string            `xml:"artifact_id"`
	Version          string            `xml:"version"`
	SnapshotVersions []SnapshotVersion `xml:"versioning>snapshotVersions>snapshotVersion"`
}

func (m *Metadata) getPomVersionSnapshotSuffix() string {
	for _, it := range m.SnapshotVersions {
		if it.Extension == "pom" {
			return it.Value
		}
	}
	return ""
}

type SnapshotVersion struct {
	XMLName    xml.Name `xml:"snapshotVersion"`
	Classifier string   `xml:"classifier"`
	Extension  string   `xml:"extension"`
	Value      string   `xml:"value"`
	Updated    string   `xml:"updated"`
}

func parsePomVersionMeta(reader io.Reader) (*Metadata, error) {
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		enc, err := ianaindex.IANA.Encoding(charset)
		if err != nil {
			return nil, fmt.Errorf("charset %s: %s", charset, err.Error())
		}
		if enc == nil {
			// Assume it's compatible with (a subset of) UTF-8 encoding
			// Bug: https://github.com/golang/go/issues/19421
			return reader, nil
		}
		return enc.NewDecoder().Reader(reader), nil
	}
	var metadata Metadata
	if e := decoder.Decode(&metadata); e != nil {
		return nil, fmt.Errorf("parse pom metadata: %w", e)
	}
	return &metadata, nil
}
