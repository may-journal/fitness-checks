package main

import (
	"github.com/may-journal/fitness-checks/internal/describe"
	"github.com/may-journal/fitness-checks/internal/lang/java/version"
	"github.com/may-journal/fitness-checks/internal/program"
	"github.com/may-journal/fitness-runner/go/pkg/checkkit"
)

func main() {
	program.Main(docSource, checkkit.Check{
		Describe: describe.JavaBuild("java-version", ""),
		Run:      version.Run,
	})
}
