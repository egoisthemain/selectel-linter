package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"linter.com/loglint/analyzer"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
