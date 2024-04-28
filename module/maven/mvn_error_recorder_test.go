package maven

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMavenErrorRecorder(t *testing.T) {
	var w = launchERecorder()
	var a = `C:\Users\Babar>mvn archetype:generate -DgroupId=com.mycompany.app -DartifactId=m
y-app -DarchetypeArtifactId=maven-archetype-quickstart -DinteractiveMode=false
[INFO] Scanning for projects...
[INFO]
[INFO] ------------------------------------------------------------------------
[INFO] Building Maven Stub Project (No POM) 1
[INFO] ------------------------------------------------------------------------
[INFO]
[INFO] >>> maven-archetype-plugin:2.1:generate (default-cli) @ standalone-pom >>
>
[INFO]
[INFO] <<< maven-archetype-plugin:2.1:generate (default-cli) @ standalone-pom <<
<
[INFO]
[INFO] --- maven-archetype-plugin:2.1:generate (default-cli) @ standalone-pom --
-
[INFO] Generating project in Batch mode
[INFO] -------------------------------------------------------------------------
---
[INFO] Using following parameters for creating project from Old (1.x) Archetype:
 maven-archetype-quickstart:1.0
[INFO] -------------------------------------------------------------------------
---
[INFO] Parameter: groupId, Value: com.mycompany.app
[INFO] Parameter: packageName, Value: com.mycompany.app
[INFO] Parameter: package, Value: com.mycompany.app
[INFO] Parameter: artifactId, Value: my-app
[INFO] Parameter: basedir, Value: C:\Users\Babar
[INFO] Parameter: version, Value: 1.0-SNAPSHOT
[INFO] ------------------------------------------------------------------------
[INFO] BUILD FAILURE
[INFO] ------------------------------------------------------------------------
[INFO] Total time: 22.971s
[INFO] Finished at: Fri Nov 18 00:07:12 EET 2011
[INFO] Final Memory: 6M/11M
[INFO] ------------------------------------------------------------------------
[ERROR] Failed to execute goal org.apache.maven.plugins:maven-archetype-plugin:2.1:generate (default-cli) on project standalone-pom: Directory my-app already exists - please run from a clean directory -> [Help 1]
[ERROR]
[ERROR] To see the full stack trace of the errors, re-run Maven with the -e switch.
[ERROR] Re-run Maven using the -X switch to enable full debug logging.
[ERROR]
[ERROR] For more information about the errors and possible solutions, please read the following articles:
[ERROR] [Help 1] http://cwiki.apache.org/confluence/display/MAVEN/MojoFailureException
'cmd' is not recognized as an internal or external command, operable program or batch file.
`
	_, e := w.Write([]byte(a))
	assert.NoError(t, e)
	assert.NoError(t, w.Close())
	t.Log(w.String())
}
