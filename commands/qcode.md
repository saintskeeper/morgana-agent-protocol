# QCODE Command

1. Run all pre-commit hooks first
2. Implement your plan following all DEF-* rules
3. Ensure new tests pass AND existing tests still pass
4. For Go: run `go test -race ./...` to check for race conditions
5. For Next.js: run `npm run typecheck` before any build
6. Run containerized integration tests if modifying API contracts
7. Run `prettier` on newly created files
8. Run `turbo typecheck lint`