package yarastore

import (
	"fmt"
	"os"

	"github.com/hillu/go-yara/v4"
)

// Compiler State that stores compiler and rules.
type CompilerState struct {
	// One time compiler to compiler added rules.
	compiler *yara.Compiler
	// Compiled ruleset of add yara rules. By default
	// will be `nil` until rules are added and compiled.
	ruleset *yara.Rules
	// Path to store the compiled rulset.
	ruleStorePath string
}

// Create a new `yara.Compiler` instance and set it.
// Note: This function sets `ruleset` to `nil` by default
func NewCompilerState(ruleStorePath string) (*CompilerState, error) {
	compiler, err := yara.NewCompiler()
	return &CompilerState{compiler, nil, ruleStorePath}, err
}

// Compile added rules and set the `ruleset` attribute
func (c *CompilerState) Compile() error {
	ruleset, err := c.compiler.GetRules()
	c.ruleset = ruleset
	return err
}

// Compile and add strings.
// Note: Will error if rule contains a syntax error and
// the compiler object will become unusable.
func (c *CompilerState) ReadString(ruleString string) error {
	return c.compiler.AddString(ruleString, "input/text")
}

// Compile and add a file contaning yara rules.
// Note: Will error if rule contains a syntax error and
// the compiler object will become unusable.
func (c *CompilerState) ReadFile(filepath string) error {
	file, err := os.Open(filepath)
	defer file.Close()
	if err != nil {
		return err
	}
	return c.compiler.AddFile(file, file.Name())
}

// Read pre-compiled yara ruleset and set `ruleset` attribute.
func (c *CompilerState) ReadCompiled(filename string) error {
	ruleset, err := yara.LoadRules(filename)
	c.ruleset = ruleset
	return err
}

// Apply compiled yara rules and find matches in a string.
// Note: `ruleset` must be compiled and set before matching.
func (c *CompilerState) MatchString(testString string) ([]yara.MatchRule, error) {
	if c.ruleset == nil {
		return nil, fmt.Errorf("Ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	testBytes := []byte(testString)
	var m yara.MatchRules
	err := c.ruleset.ScanMem(testBytes, 0, 0, &m)
	return m, err
}

// Apply compiled yara rules and find matches in a file.
// Note: `ruleset` must be compiled and set before matching.
func (c *CompilerState) MatchFile(filepath string) ([]yara.MatchRule, error) {
	if c.ruleset == nil {
		return nil, fmt.Errorf("Ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	var m yara.MatchRules
	err := c.ruleset.ScanFile(filepath, 0, 0, &m)
	return m, err
}

// Save the compiled rules in a file.
// Note: `ruleset` must be compiled and set before matching.
func (c *CompilerState) Save() error {
	if c.ruleset == nil {
		return fmt.Errorf("Ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	c.ruleset.Save(c.ruleStorePath)
	return nil
}
