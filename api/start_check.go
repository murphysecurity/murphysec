package api

import (
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/toolver"
)

func StartCheck(client *Client, task *model.ScanTask) error {
	checkNotNull(client)
	// 后端一定要我把 maven 参数再传一遍
	var data = map[string]any{
		"subtask_id":           task.SubtaskId,
		"package_private_name": task.MavenSourceName,
		"skip_skill_scan":      task.SkipSkillScan,
	}
	if task.MavenSourceId != "" {
		data["package_private_id"] = task.MavenSourceId
	}
	var buildOptions = make(map[string]any)
	data["build_options"] = buildOptions
	var maven = make(map[string]any)
	buildOptions["maven"] = maven
	if toolver.Default.Maven.MavenVersion != "" {
		maven["maven_version"] = "jdk" + toolver.Default.Maven.MavenVersion
	}
	if toolver.Default.Maven.JdkVersion != "" {
		maven["jdk_version"] = "maven" + toolver.Default.Maven.JdkVersion
	}
	if len(toolver.Default.Maven.AdditionalArgs) > 0 {
		maven["arguments"] = toolver.Default.Maven.AdditionalArgs
	}
	if len(toolver.Default.Maven.AdditionalPrependArgs) > 0 {
		maven["prepend_arguments"] = toolver.Default.Maven.AdditionalPrependArgs
	}
	return client.DoJson(client.PostJson(joinURL(client.baseUrl, "/platform3/v3/client/start_check"), data), nil)
}
