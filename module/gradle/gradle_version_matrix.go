package gradle

var versionMatrix = [][3]string{
	{"5.0", "gradle-5.0", "jdk-11.0.22+7"},
	{"5.1.1", "gradle-5.1.1", "jdk-11.0.22+7"},
	{"5.2.1", "gradle-5.2.1", "jdk-11.0.22+7"},
	{"5.3.1", "gradle-5.3.1", "jdk-11.0.22+7"},
	{"5.4.1", "gradle-5.4.1", "jdk-11.0.22+7"},
	{"5.5.1", "gradle-5.5.1", "jdk-11.0.22+7"},
	{"5.6.4", "gradle-5.6.4", "jdk-11.0.22+7"},
	{"6.0.1", "gradle-6.0.1", "jdk-11.0.22+7"},
	{"6.1.1", "gradle-6.1.1", "jdk-11.0.22+7"},
	{"6.2.2", "gradle-6.2.2", "jdk-11.0.22+7"},
	{"6.3", "gradle-6.3", "jdk-11.0.22+7"},
	{"6.4.1", "gradle-6.4.1", "jdk-11.0.22+7"},
	{"6.5.1", "gradle-6.5.1", "jdk-11.0.22+7"},
	{"6.6.1", "gradle-6.6.1", "jdk-11.0.22+7"},
	{"6.7.1", "gradle-6.7.1", "jdk-11.0.22+7"},
	{"6.8.3", "gradle-6.8.3", "jdk-11.0.22+7"},
	{"6.9.4", "gradle-6.9.4", "jdk-11.0.22+7"},
	{"7.0.2", "gradle-7.0.2", "jdk-11.0.22+7"},
	{"7.1.1", "gradle-7.1.1", "jdk-11.0.22+7"},
	{"7.2", "gradle-7.2", "jdk-11.0.22+7"},
	{"7.3.3", "gradle-7.3.3", "jdk-17.0.10+7"},
	{"7.4.2", "gradle-7.4.2", "jdk-17.0.10+7"},
	{"7.5.1", "gradle-7.5.1", "jdk-17.0.10+7"},
	{"7.6.4", "gradle-7.6.4", "jdk-17.0.10+7"},
	{"8.0.2", "gradle-8.0.2", "jdk-17.0.10+7"},
	{"8.1.1", "gradle-8.1.1", "jdk-17.0.10+7"},
	{"8.2.1", "gradle-8.2.1", "jdk-17.0.10+7"},
	{"8.3", "gradle-8.3", "jdk-17.0.10+7"},
	{"8.4", "gradle-8.4", "jdk-17.0.10+7"},
	{"8.5", "gradle-8.5", "jdk-17.0.10+7"},
	{"8.6", "gradle-8.6", "jdk-17.0.10+7"},
}

func selectGradleAndJavaVersion(input string) (string, string) {
	var a = -1
	var b = -1
	for i, it := range versionMatrix {
		var m = lenStringSamePrefix(input, it[0])
		if m > a {
			a = m
			b = i
		}
	}
	if b != -1 {
		return "/opt/gradle/" + versionMatrix[b][1] + "/bin/gradle", "/opt/openjdk/" + versionMatrix[b][2]
	}
	return "", ""
}

func lenStringSamePrefix(a, b string) (i int) {
	var m = len(a)
	if b < a {
		m = len(b)
	}
	for i := 0; i < m; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return m
}
