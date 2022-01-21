package inspector

import "testing"

func Test_reportIdeaStatus(t *testing.T) {
	reportIdeaStatus(IdeaUnknownErr, "")
	reportIdeaStatus(IdeaSucceed, "")
}
