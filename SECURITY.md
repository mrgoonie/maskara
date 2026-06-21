# Security Policy

## Reporting

Please report vulnerabilities privately through GitHub Security Advisories when
available. If that is unavailable, contact the maintainers without including
real credentials in the report body.

## Secret Handling

Maskara intentionally scans sensitive local files. Bug reports and test cases
must use generated dummy values. Do not paste live credentials into issues,
pull requests, discussions, logs, or screenshots.

## Runtime Model

Maskara is designed to run offline. It does not verify whether a detected token
is active. Treat every finding as potentially exposed and rotate it with the
owning provider.
