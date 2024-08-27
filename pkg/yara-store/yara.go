package yarastore

import (
	"os"

	"github.com/hillu/go-yara/v4"
)

type CompilerState struct {
	compiler      *yara.Compiler
	ruleStorePath string
}

func NewCompilerState(ruleStorePath string) (*CompilerState, error) {
	compiler, err := yara.NewCompiler()
	return &CompilerState{compiler, ruleStorePath}, err
}

func (c *CompilerState) ReadString(ruleString string) error {
	return c.compiler.AddString(ruleString, "input/text")
}

func (c *CompilerState) ReadFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	return c.compiler.AddFile(file, file.Name())
}

func (c *CompilerState) ReadDir(dirpath string) error {
	files, err := ListDirWithExt(dirpath, ".yara")
	for _, file := range files {
		if err := c.ReadFile(file); err != nil {
			return err
		}
	}
	return err
}

func (c *CompilerState) MatchString(testString string) ([]yara.MatchRule, error) {
	rules, err := c.compiler.GetRules()
	if err != nil {
		return nil, err
	}

	testBytes := []byte(testString)
	var m yara.MatchRules
	err = rules.ScanMem(testBytes, 0, 0, &m)
	return m, err
}

func (c *CompilerState) MatchFile(filepath string) ([]yara.MatchRule, error) {
	rules, err := c.compiler.GetRules()
	if err != nil {
		return nil, err
	}

	var m yara.MatchRules
	err = rules.ScanFile(filepath, 0, 0, &m)
	return m, err
}

func (c *CompilerState) MatchDir(dirpath string) ([]yara.MatchRule, error) {
	rules, err := c.compiler.GetRules()
	if err != nil {
		return nil, err
	}

	var m yara.MatchRules
	files, err := ListDir(dirpath)
	for _, filepath := range files {
		if err = rules.ScanFile(filepath, 0, 0, &m); err != nil {
			return nil, err
		}
	}
	return m, err
}
