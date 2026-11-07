// Command publish puts a module, at the version in its VERSION file, into a
// file-based Go module proxy: our stand-in for a private module proxy.
//
//	go run ./tools/publish -out dist modules/cart modules/payments ...
//
// For each module it writes the layout the go command reads with
// GOPROXY=file:///path/to/dist (see `go help goproxy`):
//
//	dist/example.com/shop/cart/@v/list           one version per line
//	dist/example.com/shop/cart/@v/v1.0.0.info    {"Version":"v1.0.0","Time":"..."}
//	dist/example.com/shop/cart/@v/v1.0.0.mod     the module's go.mod
//	dist/example.com/shop/cart/@v/v1.0.0.zip     the module's files, made by golang.org/x/mod/zip
//
// Before publishing it checks that the module's CHANGELOG.md has an entry
// for the version, so nobody ships an undocumented release.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"golang.org/x/mod/zip"
)

func main() {
	out := flag.String("out", "dist", "the file-based module proxy to publish into")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "usage: publish -out DIR MODULE_DIR...")
		os.Exit(2)
	}
	for _, dir := range flag.Args() {
		v, err := publish(dir, *out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "publish %s: %v\n", dir, err)
			os.Exit(1)
		}
		fmt.Printf("%s@%s\n", v.Path, v.Version)
	}
}

func publish(dir, out string) (module.Version, error) {
	goMod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		return module.Version{}, err
	}
	path := modfile.ModulePath(goMod)
	if path == "" {
		return module.Version{}, errors.New("go.mod has no module line")
	}
	raw, err := os.ReadFile(filepath.Join(dir, "VERSION"))
	if err != nil {
		return module.Version{}, err
	}
	version := strings.TrimSpace(string(raw))
	if !semver.IsValid(version) || semver.Canonical(version) != version {
		return module.Version{}, fmt.Errorf("VERSION %q is not a full semantic version such as v1.2.3", version)
	}
	if err := checkChangelog(dir, version); err != nil {
		return module.Version{}, err
	}
	mv := module.Version{Path: path, Version: version}

	escaped, err := module.EscapePath(path)
	if err != nil {
		return mv, err
	}
	vdir := filepath.Join(out, filepath.FromSlash(escaped), "@v")
	if err := os.MkdirAll(vdir, 0o755); err != nil {
		return mv, err
	}
	base := filepath.Join(vdir, version)

	info, _ := json.Marshal(struct {
		Version string
		Time    time.Time
	}{version, time.Now().UTC().Truncate(time.Second)})
	if err := os.WriteFile(base+".info", info, 0o644); err != nil {
		return mv, err
	}
	if err := os.WriteFile(base+".mod", goMod, 0o644); err != nil {
		return mv, err
	}
	if err := writeZip(base+".zip", mv, dir); err != nil {
		return mv, err
	}
	return mv, addToList(filepath.Join(vdir, "list"), version)
}

// writeZip zips the module the way the go command expects: every file
// under dir (minus nested modules and vendor/), with paths prefixed by
// "<module>@<version>/", within the size limits the go command enforces.
func writeZip(file string, mv module.Version, dir string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	if err := zip.CreateFromDir(f, mv, dir); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func addToList(file, version string) error {
	var versions []string
	if raw, err := os.ReadFile(file); err == nil {
		versions = strings.Fields(string(raw))
	}
	if !slices.Contains(versions, version) {
		versions = append(versions, version)
	}
	semver.Sort(versions)
	return os.WriteFile(file, []byte(strings.Join(versions, "\n")+"\n"), 0o644)
}

func checkChangelog(dir, version string) error {
	raw, err := os.ReadFile(filepath.Join(dir, "CHANGELOG.md"))
	if err != nil {
		return err
	}
	heading := "## " + strings.TrimPrefix(version, "v")
	for _, line := range strings.Split(string(raw), "\n") {
		if strings.TrimSpace(line) == heading {
			return nil
		}
	}
	return fmt.Errorf("CHANGELOG.md has no %q entry", heading)
}
