# vantrue-gpx

Small CLI that reads Vantrue dashcam `GPS/*.dat` files and writes [GPX 1.1](https://www.topografix.com/GPX/1/1/) tracks (including speed and course where present).

## Usage

```text
vantrue-gpx -i <file-or-dir> [-o <output-path-or-dir>] [-speed-column knots|kmh]
```

- **`-i`**: path to a single `.dat` file or a directory (each `*.dat` is converted).
- **`-o`**: optional output file or directory (default: next to each input `.dat`).
- **`-speed-column`**: `knots` (default, converted to km/h internally) or `kmh`.

Build from source with Go 1.22+ (`go build` in the module root).

## Releasing a new version with commit messages

Version bumps follow **[Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)**. go-semantic-release uses those messages (not arbitrary text) to choose **patch**, **minor**, or **major** releases.

| When you want… | Commit message shape | Typical version bump |
|----------------|----------------------|----------------------|
| Bug fixes, small compatible changes | `fix: …` | **Patch** (e.g. 1.0.0 → 1.0.1) |
| New compatible features | `feat: …` | **Minor** (e.g. 1.0.0 → 1.1.0) |
| Breaking API/behavior changes | `feat!: …` or `fix!: …`, or any type with a **BREAKING CHANGE:** footer in the body | **Major** (e.g. 1.0.0 → 2.0.0) |

Other conventional types (`docs:`, `chore:`, `refactor:`, …) usually do **not** trigger a release by themselves unless your analyzer treats them as releasable (the default setup is oriented around `feat` / `fix` / breaking changes).

**Practical flow**

1. Branch from `master`, implement changes, and use conventional commit subjects (and breaking footers when needed).
2. Open a PR, review, merge to **`master`**.
3. After the workflow runs, check **Releases** on GitHub for the new tag, notes, and binaries.

For the full commit grammar (scopes, bodies, footers, breaking change notation), see the [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/). For behavior specific to this toolchain (plugins, config file `.semrelrc`, etc.), see the [go-semantic-release documentation](https://github.com/go-semantic-release/semantic-release/blob/master/README.md).
