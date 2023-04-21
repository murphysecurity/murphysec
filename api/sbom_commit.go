package api

import (
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
)

type SubmitSBOMTreeNode struct {
	Dependencies []SubmitSBOMTreeNode `json:"dependencies,omitempty"`
	CompName     string               `json:"comp_name"`
	CompVersion  string               `json:"comp_version"`
	Ecosystem    string               `json:"ecosystem"`
	Repository   string               `json:"repository"`
}

type SubmitSBOMModule struct {
	ModuleName     string               `json:"module_name"`
	ModuleVersion  string               `json:"module_version"`
	ModulePath     string               `json:"module_path"`
	PackageManager string               `json:"package_manager"`
	Dependencies   []SubmitSBOMTreeNode `json:"dependencies,omitempty"`
}

type SBOMSubmitRequest struct {
	SubtaskId string             `json:"subtask_id"`
	Modules   []SubmitSBOMModule `json:"modules"`
}

func SubmitSBOM(client *Client, subtaskId string, modules []model.Module, codeFragments []model.ComponentCodeFragment) error {
	checkNotNull(client)
	must.NotZero(subtaskId)
	var req = map[string]any{
		"subtask_id":     subtaskId,
		"modules":        utils.NoNilSlice(modules),
		"code_fragments": utils.NoNilSlice(codeFragments),
	}
	return client.DoJson(client.PostJson(joinURL(client.baseUrl, "/platform3/v3/client/upload_data"), req), nil)
}
