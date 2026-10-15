# Contributing

Bug reports, ideas and pull requests are welcome.

- **Bugs and ideas:** open an [issue](https://github.com/eslusarenko/apg/issues).
  For security problems see [SECURITY.md](SECURITY.md).
- **Pull requests:** keep them focused, add tests, and run `go vet ./...` and
  `go test ./...` before pushing.

A few principles the project sticks to:

- Options of the classic apg keep their meaning. New features are added as
  extensions, preferably as long options.
- Standard library only, no dependencies.
- Every random choice uses `crypto/rand` and is uniform.

If you change options, update the help text in `main.go`, the man page
(`man/apg.1`), the shell completions (`completions/`) and the README.
