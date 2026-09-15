# Repository Security Settings

English | [日本語](SECURITY_SETTINGS.ja.md)

This file records GitHub-side security settings that cannot be enforced only by files in the repository.

The repository owner should keep the following enabled where GitHub provides the feature for this public repository:

- require pull requests for changes to `main`;
- require the `Strict gate` and `Security gate` status checks before merge;
- block force pushes and branch deletion for `main`;
- require conversations to be resolved before merge where practical;
- keep the default `GITHUB_TOKEN` permission read-only and grant write permissions only per job when necessary;
- enable Dependabot alerts and Dependabot security updates;
- enable Secret Scanning and Push Protection;
- enable private vulnerability reporting / repository Security Advisories;
- keep CodeQL code scanning enabled and review new alerts;
- require approval before running workflows from untrusted first-time fork contributors where the repository setting is available;
- review the allowed-actions policy and restrict third-party Actions when practical.

Repository workflows additionally pin all Actions to immutable commit SHAs and use least-privilege job permissions.

Changes to these settings should be treated as security-sensitive configuration changes and reviewed against `SECURITY.md`.
