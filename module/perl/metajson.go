package perl

import (
	"fmt"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils/simplejson"
	"sort"
)

type meta struct {
	Name        string
	Version     string
	RuntimeReqs reqItemList
}

func (m *meta) deps() (r []model.DependencyItem) {

	for _, it := range m.RuntimeReqs {
		r = append(r, model.DependencyItem{
			Component: model.Component{
				CompName:    it.Name,
				CompVersion: it.Version,
				EcoRepo:     EcoRepo,
			},
		})
	}
	return
}

type reqItem struct {
	Name    string
	Version string
}

type reqItemList []reqItem

func (r reqItemList) Len() int {
	return len(r)
}

func (r reqItemList) Less(i, j int) bool {
	if r[i].Name != r[j].Name {
		return r[i].Name < r[j].Name
	}
	return r[i].Version < r[j].Version
}

func (r reqItemList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func parseMeta(data []byte) (r *meta, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("parseMeta: %w", err)
		}
	}()
	doc, e := simplejson.NewJSON(data)
	if e != nil {
		return nil, e
	}
	m := &meta{
		Name:        doc.Get("name").String(),
		Version:     doc.Get("version").String(),
		RuntimeReqs: reqItemList{},
	}
	for k, v := range doc.Get("prereqs", "runtime", "requires").JSONMap() {
		if v == nil {
			continue
		}
		m.RuntimeReqs = append(m.RuntimeReqs, reqItem{
			Name:    k,
			Version: v.String(),
		})
	}
	sort.Sort(m.RuntimeReqs)
	return m, nil
}
