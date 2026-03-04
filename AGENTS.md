# Agent Rules

## Git Commit Signing (MANDATORY)

1. This repository enforces signed commits as a mandatory rule. Every `git commit` must be signed (GPG or SSH signing, following local Git configuration).
2. `--no-gpg-sign`, disabling signing config, or any equivalent bypass is strictly prohibited.
3. In this Windows sandbox environment, commit signing will predictably fail without elevation (for example, keyring/lockfile/pinentry access issues). Do not make repeated non-elevated commit attempts.
4. If signing fails or is expected to fail in sandbox, you must immediately request escalation and retry with signing enabled.
5. If escalation is not granted, stop and ask the user how to proceed. Do not create an unsigned commit.
