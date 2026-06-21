# CLAUDE.md - Version Management Lessons

## Key Takeaway: ldflags Injection Over Hardcoding

The documentation warns against hardcoding version numbers in CLI tools. Instead, use build-time injection:

### Anti-Pattern: Hardcoded Version

In cmd/go-code/main.go, avoid:
```go
var version = "0.1.0"  // ❌ Forgot to update = stale releases
```

### Best Practice: ldflags Injection

- Set default to "dev" in code
- Inject actual version via GitHub Actions
- Ensures every release has correct version number

### Release Checklist

1. Verify cmd/go-code/main.go version defaults to "dev"
2. Confirm .github/workflows/release.yml uses ldflags injection
3. Test version locally before tagging
4. Verify downloaded binary with `--version`
5. Update README/changelog

### Release Process

```bash
# 1. Pull latest
git pull origin main

# 2. Create tag with v prefix
git tag v0.1.0

# 3. Push tag
git push origin v0.1.0

# 4. GitHub Actions will build and release
```

## Why This Matters

Manual version updates get forgotten. Automated injection ensures build process always produces accurate version metadata without human intervention.
