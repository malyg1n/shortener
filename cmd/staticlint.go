package main

import (
	grit "github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/malyg1n/shortener/pkg/analyze"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"strings"
)

func main() {
	var checks []*analysis.Analyzer
	checks = append(checks, analyze.OsExitCheckCheckAnalyzer)
	checks = append(checks, printf.Analyzer)
	checks = append(checks, shadow.Analyzer)
	checks = append(checks, structtag.Analyzer)
	checks = append(checks, grit.Analyzer)

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			checks = append(checks, v.Analyzer)
		}
		if strings.HasPrefix(v.Analyzer.Name, "ST1") {
			checks = append(checks, v.Analyzer)
		}
	}

	multichecker.Main(
		checks...,
	)
}
