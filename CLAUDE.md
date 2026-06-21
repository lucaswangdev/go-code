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
5. **Update README.md version number** (find and replace old version with new tag version)
6. Commit README.md change

### Release Process

```bash
# 1. Pull latest
git pull origin main

# 2. Update README.md version (replace old version with new, e.g., v0.1.0 -> v0.2.0)
# sed -i 's/v0.1.0/v0.2.0/g' README.md

# 3. Commit README change
git add README.md && git commit -m "Bump version to v0.2.0"

# 4. Create tag with v prefix
git tag v0.2.0

# 5. Push tag
git push origin v0.2.0

# 6. GitHub Actions will build and release
```

## Why This Matters

Manual version updates get forgotten. Automated injection ensures build process always produces accurate version metadata without human intervention.
