package api

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

func SBOMSubmit(client *Client, req *SBOMSubmitRequest) error {
	checkNotNull(client)
	checkNotNull(req)
	return client.DoJson(client.PostJson("/v3/client/upload_data", req), nil)
}
