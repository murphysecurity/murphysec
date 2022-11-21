package api

import "time"

type CreateSubTaskRequest struct {
	AccessType  AccessType `json:"access_type"` // 接入方式
	Addr        *string    `json:"addr,omitempty"`
	Author      *string    `json:"author,omitempty"` // 作者/提交者Email
	Branch      *string    `json:"branch,omitempty"`
	Commit      *string    `json:"commit,omitempty"`
	Message     *string    `json:"message,omitempty"`   // 提交信息
	PushTime    *time.Time `json:"push_time,omitempty"` // 提交时间
	SubtaskName string     `json:"subtask_name"`
	TaskID      *string    `json:"task_id,omitempty"`
}

type CreateSubTaskResponse struct {
	ProjectsName string `json:"projects_name"` // 项目名称
	TaskID       string `json:"task_id"`       // 任务ID
	TaskName     string `json:"task_name"`     // 任务名称
}

// AccessType 接入方式
type AccessType string

const (
	AccessTypeCLI    AccessType = "cli"
	AccessTypeIdea   AccessType = "idea"
	AccessTypeBinary AccessType = "binary"
	AccessTypeIOT    AccessType = "iot"
)

func CreateSubTask(client *Client, request *CreateSubTaskRequest) (*CreateSubTaskResponse, error) {
	checkNotNull(client)
	checkNotNull(request)
	var resp CreateSubTaskResponse
	if e := client.DoJson(client.PostJson("/v3/client/create_subtask", request), &resp); e != nil {
		return nil, e
	}
	return &resp, nil
}
