package yarastore

import (
	"fmt"
	"os"

	"github.com/hillu/go-yara/v4"

    "yarastore/pkg/utils"
)

type CompilerState struct {
	compiler      *yara.Compiler
	ruleset       *yara.Rules
	ruleStorePath string
}

func NewCompilerState(ruleStorePath string) (*CompilerState, error) {
	compiler, err := yara.NewCompiler()
	return &CompilerState{compiler, nil, ruleStorePath}, err
}

func (c *CompilerState) Compile() error {
	ruleset, err := c.compiler.GetRules()
	c.ruleset = ruleset
	return err
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
	files, err := utils.ListDirWithExt(dirpath, ".yara")
	for _, file := range files {
		if err := c.ReadFile(file); err != nil {
			return err
		}
	}
	return err
}

func (c *CompilerState) ReadCompiled(filename string) error {
	ruleset, err := yara.LoadRules(filename)
	c.ruleset = ruleset
	return err
}

func (c *CompilerState) MatchString(testString string) ([]yara.MatchRule, error) {
	if c.ruleset == nil {
		return nil, fmt.Errorf("Ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	testBytes := []byte(testString)
	var m yara.MatchRules
	err := c.ruleset.ScanMem(testBytes, 0, 0, &m)
	return m, err
}

func (c *CompilerState) MatchFile(filepath string) ([]yara.MatchRule, error) {
	if c.ruleset == nil {
		return nil, fmt.Errorf("Ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	var m yara.MatchRules
	err := c.ruleset.ScanFile(filepath, 0, 0, &m)
	return m, err
}

func (c *CompilerState) MatchDir(dirpath string) ([]yara.MatchRule, error) {
	if c.ruleset == nil {
		return nil, fmt.Errorf("Ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	var m yara.MatchRules
	files, err := utils.ListDir(dirpath)
	for _, filepath := range files {
		if err = c.ruleset.ScanFile(filepath, 0, 0, &m); err != nil {
			return nil, err
		}
	}
	return m, err
}

func (c *CompilerState) Save() error {
	if c.ruleset == nil {
		return fmt.Errorf("Ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	c.ruleset.Save(c.ruleStorePath)
	return nil
}
