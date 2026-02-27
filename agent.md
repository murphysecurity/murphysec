# Agent Rules

## Git Commit Signing

1. All `git commit` operations in this repository must be signed (GPG/SSH signing is acceptable, following local Git configuration).
2. If signing fails due to permission issues (for example keyring, pinentry, or lockfile access errors), you must request elevated permissions and retry.
3. Do not bypass signing requirements under any circumstances, including `--no-gpg-sign`, disabling signing config, or any equivalent workaround.
