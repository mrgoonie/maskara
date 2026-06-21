$payload = if ($args.Count -gt 0) { $args -join " " } else { [Console]::In.ReadToEnd() }
$lower = $payload.ToLowerInvariant()

$blocked = @(
  ".env",
  "printenv",
  "authorization:",
  "bearer ",
  "private key",
  "secret"
)

foreach ($term in $blocked) {
  if ($lower.Contains($term)) {
    Write-Error "maskara guardrails: command may expose secrets. Use masked checks or run maskara report."
    exit 2
  }
}

exit 0
