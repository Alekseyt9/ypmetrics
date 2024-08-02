// Package staticlint provides a multichecker.
//
// The multichecker is created to help developers maintain high code quality
// by performing static analysis on their Go codebases. The analyzers included
// in the multichecker cover a wide range of potential issues, including
// performance, correctness, and code style.
//
// The multichecker can be run using the following command:
//
//	go run cmd/staticlint
//
// Each included analyzer is documented below.
//
// * Standard Static Analyzers
//
// * Staticcheck Analyzers:
//   - All SA class analyzers: Detect various issues in code.
//
// * Additional Staticcheck Analyzers:
//   - One analyzer from other classes.
//
// * Public Analyzers:
//   - sqlrows: Detects issues with SQL rows usage.
//   - emptycase: Detects empty case blocks in switch statements.
//
// * Custom Analyzer:
//   - mainexit: Prohibits the direct use of os.Exit in the main function of the main package.
package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"

	"github.com/gostaticanalysis/emptycase"
	"github.com/gostaticanalysis/sqlrows/passes/sqlrows"
	"honnef.co/go/tools/staticcheck"

	"github.com/Alekseyt9/ypmetrics/cmd/staticlint/mainexitanalyzer"
)

func main() {
	var analyzers []*analysis.Analyzer
	addStaticAnalyzers(&analyzers)
	addStaticcheckAnalyzers(&analyzers)
	addPublicAnalyzers(&analyzers)
	addOwnAnalyzers(&analyzers)
	multichecker.Main(analyzers...)
}

// собственный анализатор
func addOwnAnalyzers(analyzers *[]*analysis.Analyzer) {
	*analyzers = append(*analyzers, mainexitanalyzer.Analyzer)
}

// два публичных анализатора
func addPublicAnalyzers(analyzers *[]*analysis.Analyzer) {
	*analyzers = append(*analyzers,
		sqlrows.Analyzer,
		emptycase.Analyzer)
}

func addStaticcheckAnalyzers(analyzers *[]*analysis.Analyzer) {
	// все анализаторы класса SA пакета staticcheck
	for _, v := range staticcheck.Analyzers {
		if v.Analyzer.Name[:2] == "SA" {
			*analyzers = append(*analyzers, v.Analyzer)
		}
	}

	// анализаторы остальных классов пакета staticcheck
	seenClasses := make(map[string]bool, 0)
	for _, v := range staticcheck.Analyzers {
		class := v.Analyzer.Name[:2]
		if class != "SA" && !seenClasses[class] {
			*analyzers = append(*analyzers, v.Analyzer)
			seenClasses[class] = true
		}
	}
}

// стандартные статические анализаторы
func addStaticAnalyzers(analyzers *[]*analysis.Analyzer) {
	*analyzers = append(*analyzers,
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stdversion.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
	)
}
