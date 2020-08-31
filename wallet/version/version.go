package version

import (
	"fmt"
	"runtime"
)

var (
	// application's name
	Name = ""
	// application's version string
	Version = ""
	// cosmos's version string
	VersionCosmos = ""
	// ethereum's version string
	VersionEthereum = ""
	// commit
	Commit = ""
	// build tags
	BuildTags = ""
	// go version
	GoVersion = fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
)

func ShowVersion() {
	fmt.Printf(`%s		: %s
with		:
  %s
  %s
git commit	: %s
build tags	: %s
go		: %s
`,
		Name, Version,
		VersionCosmos, VersionEthereum,
		Commit,
		BuildTags,
		GoVersion,
	)
}
