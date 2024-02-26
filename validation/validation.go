package validation

import (
	"fmt"
	"os"
)

// adapted from this https://dev.to/ansu/best-practices-for-building-a-validation-layer-in-go-59j9

type Rule func() error

type Rules []Rule

type Validator struct {
	rules Rules
}

// NewValidator automatically constructs a validator object
// with the given rules
func NewValidator(rules []Rule) *Validator {
	return &Validator{rules}
}

// Add adds a new rule to the validator
func (v *Validator) Add(rule Rule) {
	v.rules = append(v.rules, rule)
}

// Validate runs all the rules and returns a list of errors
func (v *Validator) Validate() []error {
	var errors []error
	for _, rule := range v.rules {
		if err := rule(); err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

// MustValidate runs all the rules and exits if any of them fail
func (v *Validator) MustValidate() {
	errors := v.Validate()
	fmt.Println("errors encountered during input validation. See -help for more information.")
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		os.Exit(1)
	}

}
