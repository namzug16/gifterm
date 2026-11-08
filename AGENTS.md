## Global Coding Standards
1. **YAGNI**: Don't build it until you need it
2. **DRY**: Extract patterns after second duplication, not before
3. **FAIL FAST**: Explicit errors beat silent failures
4. **SIMPLE FIRST**: Write the obvious solution, optimize only if needed
5. **DELETE AGGRESSIVELY**: Less code = fewer bugs
6. **SEMANTIC NAMING**: Always name variables, parameters, and API endpoints with verbose, self-documenting names that optimize for comprehension by both human and LLMs, not brevity (e.g. `wait_until_obs_is_saved=true` vs `wait=true`)

## Commits
Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) when creating commits.

**Format:** `<type>(<optional scope>): <description>`

- **Types:** `feat`, `fix`, `refactor`, `perf`, `docs`, `style`, `test`, `build`, `ci`, `chore`, `revert`.
- Use imperative, lowercase descriptions without a trailing period (e.g., `feat(auth): add password reset`).
- Keep the subject under 72 characters, ideally 50.
- Make commits atomic: one logical change per commit.
- Include a body when additional context or rationale is needed.
- Mark breaking changes with `!` after the type/scope and include a `BREAKING CHANGE:` footer.
- Reference related issues in footers when applicable (e.g., `Fixes #123`).
- Never commit unrelated changes or broken code.

## Releases
Pushing a `vX.Y.Z` tag runs `.github/workflows/release.yml`: it checks the tag against `version` in `main.go`, runs vet and tests, builds `make release-artifacts`, publishes a GitHub release with the tag's `CHANGELOG.md` section as notes, and updates `Formula/gifterm.rb` in `namzug16/homebrew-tap`.

To prepare a release:
1. Pick the version with [SemVer](https://semver.org): breaking flag/behavior changes bump minor while < 1.0 (major after), `feat` bumps minor, only `fix`/`perf` bumps patch.
2. Read `git log <last tag>..HEAD` and add a `## X.Y.Z - YYYY-MM-DD` section at the top of `CHANGELOG.md`. Write user-facing bullets (what changed for someone running `gifterm`), not commit subjects; skip `test`, `refactor`, `chore`, `ci` and `docs` unless users notice them. The heading must start with `## X.Y.Z ` or the workflow fails.
3. Set `const version = "X.Y.Z"` in `main.go`.
4. Run `go vet ./... && go test ./... && make release-artifacts`.
5. Commit as `chore(release): vX.Y.Z` with only those changes.
6. Ask before tagging and pushing, since it publishes: `git tag vX.Y.Z && git push origin master vX.Y.Z`.

Never move or reuse a published tag; fix forward with a new patch version.
