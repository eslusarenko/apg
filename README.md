# apg — Advanced Password Generator

A command-line password generator inspired by
[apg, the Automated Password Generator](https://github.com/wilx/apg) by
Adel I. Mirzazhanov — a long-time Linux favourite that is missing from
Homebrew.

This apg keeps the classic command-line options, including FIPS-181
pronounceable passwords, and adds what makes everyday use simpler:
strong defaults, PIN codes, copying to the clipboard and a config file.
It is a single binary with no dependencies for macOS, Linux and Windows.

It is a new implementation in Go, not affiliated with the original project.

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

Options (compatible with the classic apg):
  -a 0|1         algorithm: 0 pronounceable, 1 random (default 1)
  -M mode        symbol classes: S special symbols, N numbers, C capital
                 letters, L small letters; upper case means every password
                 must contain the class, lower case means it may (default SNCL)
  -E chars       exclude these characters
  -n N           number of passwords to generate (default 1)
  -m N           minimum password length (default 32)
  -x N           maximum password length (default 32)
  -r file        reject passwords found in a dictionary file
  -t             print the pronunciation of pronounceable passwords
  -l             spell out generated passwords
  -d             do not print a newline after each password
  -q             quiet mode: do not print warnings
  -s, -c seed    random seed options; accepted for compatibility and ignored
  -v, --version  print version and exit
  -h, --help     show this help

Extensions:
  --length N        exact password length, same as -m N -x N
  --strategy NAME   full (-M SNCL), alphanum (-M NCL) or digits (-M N)
  --alnum           same as --strategy alphanum
  --digits, --pin   same as --strategy digits
  --no-newline      same as -d
  -C, --copy        copy the generated passwords to the clipboard
```

Examples:

```sh
$ apg                                  # one strong 32-character password
xVEgPf%HgA3sRI)a_@{eQvX5^p}}&WQh
$ apg -a 0 -n 3 -m 8 -x 10 -M sncl     # pronounceable, like the classic apg
Hecdicirr
wehycsyap
nerpecjic
$ apg -a 0 -m 10 -x 10 -M ncl -t       # with pronunciation
itVoijOpuv (it-Voij-Op-uv)
$ apg -m 8 -x 8 -l                     # spelled out
A4c[up6d Alfa-FOUR-charlie-LEFT_BRACKET-uniform-papa-SIX-delta
$ apg -E 'O0Il1|' --length 16          # no look-alike characters
2#S._4$,SzGv2^tp
$ apg --pin --length 6                 # PIN code
118588
$ apg --alnum --length 20              # letters and digits only
1JwuVJf545B7jIXvXsZl
$ apg -d | pbcopy                      # no trailing newline in the clipboard
```

### Strategies

Strategies are named presets for `-M`:

| Name       | Mode   | Characters                             | Flags               |
|------------|--------|----------------------------------------|---------------------|
| `full`     | `SNCL` | letters, digits and special characters | `--strategy full`   |
| `alphanum` | `NCL`  | letters and digits                     | `--alnum`           |
| `digits`   | `N`    | digits only, e.g. for PIN codes        | `--digits`, `--pin` |

If several strategy or `-M` options are given, the last one wins.

Special characters are the 32 ASCII punctuation characters, as in the
classic apg: `` !"#$%&'()*+,-./:;<=>?@[\]^_`{|}~ ``. Use `-E` to exclude
the ones that get in your way.

With `--copy` the passwords are also put on the clipboard (`pbcopy` on
macOS, `clip` on Windows, `wl-copy`, `xclip` or `xsel` on Linux). Keep in
mind that clipboard managers may store them in their history.

## Coming from the classic apg

The options mean the same, with these differences:

- **Defaults.** Without options apg prints one random 32-character password
  (`-a 1 -n 1 -m 32 -x 32 -M SNCL`). The classic apg prints six
  pronounceable 8–10 character passwords; to get that, run
  `apg -a 0 -n 6 -m 8 -x 10 -M sncl`.
- **Randomness.** Passwords always come from the operating system's
  cryptographic random generator, so the seed options `-s` and `-c` are
  accepted and ignored.
- **Not supported:** Bloom filter checks (`-b`, `-p` and the `apgbfm`
  utility), crypt(3) output (`-y`) and cracklib checks (`-k`). apg exits
  with an error instead of silently skipping a check.
- **Short options** cannot be combined: write `-d -l`, not `-dl`.

## Configuration

Defaults can be changed in `~/.apg.conf`:

```ini
# password length when no length option is given
default_length=24
# strategy when no strategy or -M option is given: full, alphanum or digits
default_strategy=alphanum
```

| Key                | Values                          | Built-in default |
|--------------------|---------------------------------|------------------|
| `default_length`   | a number                        | `32`             |
| `default_strategy` | `full`, `alphanum` or `digits`  | `full`           |

One `key=value` per line; blank lines and lines starting with `#` are
ignored. Command-line options take precedence over the file.

## How it works

- Every random choice uses Go's `crypto/rand` and is uniform.
- Random passwords (`-a 1`) pick each character from the allowed set.
  Passwords that miss a required class (upper case letters in `-M`) or
  appear in the `-r` dictionary are discarded and generated again.
- Pronounceable passwords (`-a 0`) use the FIPS-181 algorithm ported from
  the classic apg: words are built from syllables that follow English
  pronunciation rules, then single-letter syllables may become digits or
  special characters and syllables may be capitalized. They are easier to
  remember and type, but have noticeably less entropy than random passwords
  of the same length.

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

[MIT](LICENSE). The pronounceable password generator is ported from the
classic apg and is distributed under its BSD license; see [NOTICE](NOTICE).
