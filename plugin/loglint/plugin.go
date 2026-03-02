package loglintplugin

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"linter.com/loglint/analyzer"
)

func init() {
	register.Plugin("loglint", New)
}

type Settings struct {
}

type Plugin struct {
	settings Settings
}

func New(raw any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](raw)
	if err != nil {
		return nil, err
	}
	return &Plugin{settings: s}, nil
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
