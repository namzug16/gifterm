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
