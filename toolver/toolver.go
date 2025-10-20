package toolver

type Config struct {
	Maven Maven `json:"maven,omitempty"`
}

type Maven struct {
	JavaHome string `json:"java_home,omitempty"`

	JdkVersion string `json:"jdk_version,omitempty"`

	MavenCommand string `json:"maven_command,omitempty"`

	MavenVersion string `json:"maven_version,omitempty"`

	MavenSettingPath string `json:"maven_setting_path,omitempty"`

	AdditionalArgs []string `json:"additional_args,omitempty"`

	AdditionalPrependArgs []string `json:"additional_prepend_args,omitempty"`
}
