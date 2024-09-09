package yarastore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hillu/go-yara/v4"

	"yarastore/pkg/config"
)

// CompilerState Compiler State that stores compiler and rules.
type CompilerState struct {
	// One time compiler to compiler added rules.
	compiler *yara.Compiler
	// Compiled ruleset of add yara rules. By default,
	// will be `nil` until rules are added and compiled.
	ruleset *yara.Rules
	// Path to store the compiled ruleset.
	ruleStorePath string
}

// RuleMatch A file-based match state which stores all yara matches in a particular file.
type RuleMatch struct {
	// File which rules were tested against.
	Filename string
	// All rule matches found in the file.
	Matches yara.MatchRules
}

// NewRuleMatch Create a new instance of `RuleMatch`
func NewRuleMatch(filename string, matches yara.MatchRules) *RuleMatch {
	return &RuleMatch{filename, matches}
}

// RuleMatchAsJsonS Serialized `RuleMatch` as json and return string.
func RuleMatchAsJsonS(rules []RuleMatch) (string, error) {
	jsonBytes, err := json.Marshal(rules)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), err
}

// RuleMatchAsJson Serialize `RuleMatch` and dump into a json file.
func RuleMatchAsJson(rules []RuleMatch, filename string) error {
	jsonBytes, err := json.Marshal(rules)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, jsonBytes, 0666); err != nil {
		return err
	}
	return nil
}

// NewCompilerState Create a new `yara.Compiler` instance and set it.
// Note: This function sets `ruleset` to `nil` by default.
func NewCompilerState(ruleStorePath string) (*CompilerState, error) {
	compiler, err := yara.NewCompiler()
	return &CompilerState{compiler, nil, ruleStorePath}, err
}

// Compile added rules and set the `ruleset` attribute.
func (c *CompilerState) Compile() error {
	ruleset, err := c.compiler.GetRules()
	c.ruleset = ruleset
	return err
}

// ReadString Compile and add strings.
// Note: Will error if rule contains a syntax error and
// the compiler object will become unusable.
func (c *CompilerState) ReadString(ruleString string) error {
	return c.compiler.AddString(ruleString, "input/text")
}

// ReadFile Compile and add a file containing yara rules.
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

// ReadConfig Compile files from `config.Config` rules.
func (c *CompilerState) ReadConfig(conf *config.Config) error {
	for _, dir := range conf.Rules.Dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			base := filepath.Base(path)
            if info.IsDir() && dir != path && !conf.Rules.Recursive {
                return filepath.SkipDir
            }
			if info.IsDir() && conf.Rules.IsDirExcluded(base) {
				return filepath.SkipDir
			}
			if !info.IsDir() && conf.Rules.IsFilenameValid(base) {
				if err := c.ReadFile(path); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	for _, file := range conf.Rules.Files {
		if conf.Rules.IsFilenameValid(file) {
			if err := c.ReadFile(file); err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadCompiled Read pre-compiled yara ruleset and set `ruleset` attribute.
func (c *CompilerState) ReadCompiled(filename string) error {
	ruleset, err := yara.LoadRules(filename)
	c.ruleset = ruleset
	return err
}

// MatchString Apply compiled yara rules and find matches in a string.
// Note: `ruleset` must be compiled and set before matching.
func (c *CompilerState) MatchString(testString string) (*RuleMatch, error) {
	if c.ruleset == nil {
		return nil, fmt.Errorf("ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	testBytes := []byte(testString)
	var m yara.MatchRules
	if err := c.ruleset.ScanMem(testBytes, 0, 0, &m); err != nil {
		return nil, err
	}
	ruleMatch := NewRuleMatch("", m)
	return ruleMatch, nil
}

// MatchFile Apply compiled yara rules and find matches in a file.
// Note: `ruleset` must be compiled and set before matching.
func (c *CompilerState) MatchFile(filepath string) (*RuleMatch, error) {
	if c.ruleset == nil {
		return nil, fmt.Errorf("ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	var m yara.MatchRules
	if err := c.ruleset.ScanFile(filepath, 0, 0, &m); err != nil {
		return nil, err
	}
	ruleMatch := NewRuleMatch(filepath, m)
	return ruleMatch, nil
}

// MatchConfig Apply compiled yara rules and find matches in files from `config.Config`.
func (c *CompilerState) MatchConfig(conf *config.Config) ([]RuleMatch, error) {
	if c.ruleset == nil {
		return nil, fmt.Errorf("ruleset not compiled! Please use `.Compile()` before performing this operation")
	}

	var matches []RuleMatch
	for _, dir := range conf.Target.Dirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

            base := filepath.Base(path)
            if info.IsDir() && dir != path && !conf.Target.Recursive {
                return filepath.SkipDir
            }
            if info.IsDir() && conf.Target.IsDirExcluded(base) {
                return filepath.SkipDir
            }
			if !info.IsDir() && conf.Target.IsFilenameValid(base) {
				ruleMatch, err := c.MatchFile(path)
				if err != nil {
					return err
				}
				matches = append(matches, *ruleMatch)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	for _, file := range conf.Target.Files {
		if conf.Target.IsFilenameValid(file) {
			ruleMatch, err := c.MatchFile(file)
			if err != nil {
				return nil, err
			}
			matches = append(matches, *ruleMatch)
		}
	}
	return matches, nil
}

// Save the compiled rules in a file.
// Note: `ruleset` must be compiled and set before matching.
func (c *CompilerState) Save() error {
	if c.ruleset == nil {
		return fmt.Errorf("ruleset not compiled! Please use `.Compile()` before performing this operation")
	}
	return c.ruleset.Save(c.ruleStorePath)
}
