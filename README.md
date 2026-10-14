# apg

A tiny command-line tool that generates strong random passwords.

By default it produces a 32-character password made of upper- and lowercase
Latin letters, digits and special characters. It can also produce
letters-and-digits passwords or digit-only PIN codes. Every password is
guaranteed to contain at least one character from each enabled group.
Randomness comes from the operating system's CSPRNG (`crypto/rand`).

## Install

### Homebrew

```sh
brew install eslusarenko/tap/apg
```

### Go

```sh
go install github.com/eslusarenko/apg@latest
```

### Binaries

Download an archive for your platform from the
[releases page](https://github.com/eslusarenko/apg/releases).

## Usage

```
apg [options]

  -l, --length N        password length (default 32)
  -s, --strategy NAME   character set to use (default full, see below)
  -a, --alnum           same as --strategy alphanum
  -d, --digits          same as --strategy digits
      --pin             same as --strategy digits
  -n, --no-newline      do not print the trailing newline
  -c, --copy            copy the password to the clipboard
  -v, --version         print version and exit
  -h, --help            show this help
```

Strategies:

| Name       | Characters                             | Flags                     |
|------------|----------------------------------------|---------------------------|
| `full`     | letters, digits and special characters | `-s full`                 |
| `alphanum` | letters and digits                     | `-a`, `--alnum`           |
| `digits`   | digits only, e.g. for PIN codes        | `-d`, `--digits`, `--pin` |

If several strategy flags are given, the last one wins.

Examples:

```sh
$ apg
r9#Vq}Lm2x_TfK@8wZ!pC4n;Hs6.dYb&
$ apg -l 16
aT7!x;Pq2@Mz}k9R
$ apg --length 20 --alnum
Xb4kLp9QzT2mVw7RcN5s
$ apg --pin -l 6    # PIN code
480913
$ apg -n | pbcopy   # no trailing newline in the clipboard
```

With `--copy` the password is also put on the clipboard (`pbcopy` on macOS,
`clip` on Windows, `wl-copy`, `xclip` or `xsel` on Linux). Keep in mind that
clipboard managers may store it in their history.

Special characters used: `` !@#$%^&*()-_=+[]{};:,.<>?/|~ ``

## Configuration

Defaults can be changed in `~/.apg.conf`:

```ini
# password length when -l is not given
default_length=24
# strategy when no strategy flag is given: full, alphanum or digits
default_strategy=alphanum
```

| Key                | Values                          | Built-in default |
|--------------------|---------------------------------|------------------|
| `default_length`   | a number                        | `32`             |
| `default_strategy` | `full`, `alphanum` or `digits`  | `full`           |

One `key=value` per line; blank lines and lines starting with `#` are
ignored. Command-line flags take precedence over the file, e.g. `-s full`
brings special characters back for a single run.

## Development

```sh
go test ./...
go build
```

## Releasing

Releases are built by [GoReleaser](https://goreleaser.com) when a `v*` tag is
pushed:

```sh
git tag v1.0.0
git push origin v1.0.0
```

The release workflow publishes binaries to GitHub Releases and updates the
cask in [eslusarenko/homebrew-tap](https://github.com/eslusarenko/homebrew-tap).
It needs a `HOMEBREW_TAP_GITHUB_TOKEN` repository secret: a token with write
access to the tap repository.

## License

[MIT](LICENSE)
