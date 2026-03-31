package conan

import (
	"context"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/pkg/errors"
	"os/exec"
	"strings"
	"testing"
)

func TestGetConanVersion(t *testing.T) {
	_, e := exec.LookPath("conan")
	if errors.Is(e, exec.ErrNotFound) {
		t.Skip("Conan not found in test environment, test skipped.")
		return
	}
	t.Log(GetConanVersion(context.TODO(), must.A(LocateConan(context.TODO()))))
}

func TestConanArgsV1DoesNotForceCoreNonInteractive(t *testing.T) {
	args := conanArgs(1, "info", ".")
	got := strings.Join(args, " ")
	if strings.Contains(got, "core:non_interactive=True") {
		t.Fatalf("unexpected v1 core:non_interactive flag in args: %v", args)
	}
}

func TestConanArgsV2ForcesCoreNonInteractive(t *testing.T) {
	args := conanArgs(2, "graph", "info", ".")
	got := strings.Join(args, " ")
	if !strings.Contains(got, "core:non_interactive=True") {
		t.Fatalf("missing v2 core:non_interactive flag in args: %v", args)
	}
}

func TestConanRemoteEnvVarSuffix(t *testing.T) {
	if got := conanRemoteEnvVarSuffix("conan-center.prod"); got != "CONAN_CENTER_PROD" {
		t.Fatalf("unexpected env suffix: %q", got)
	}
}

func TestGetEnvForConanAddsRemoteCredentialEnv(t *testing.T) {
	env := getEnvForConan(2, []conanRemoteCredential{{
		Name:     "conan-center",
		Username: "alice",
		Password: "secret",
	}})
	joined := strings.Join(env, "\n")
	if !strings.Contains(joined, "CONAN_LOGIN_USERNAME_CONAN_CENTER=alice") {
		t.Fatalf("missing username env, env=%v", env)
	}
	if !strings.Contains(joined, "CONAN_PASSWORD_CONAN_CENTER=secret") {
		t.Fatalf("missing password env, env=%v", env)
	}
}
