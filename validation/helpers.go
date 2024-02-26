package validation

import "fmt"

// ValidateNonEmpty will assert that the value is not the empty string
func ValidateNonEmpty(key string, value interface{}) Rule {
	return func() error {
		if value == "" {
			return fmt.Errorf("%s can't be blank", key)
		}
		return nil
	}
}

// TODO: Refactor this and ValidateRangeFloat with generic types
func ValidateRangeInt(min, max int, value interface{}) Rule {
	return func() error {
		num, ok := value.(int)
		if !ok {
			return fmt.Errorf("%s is not a number", value)
		}
		if num < min || num > max {
			return fmt.Errorf("%s must be between %d and %d", value, min, max)
		}
		return nil
	}
}

func ValidateRangeFloat(min, max float64, value interface{}) Rule {
	return func() error {
		num, ok := value.(float64)
		if !ok {
			return fmt.Errorf("%s is not a number", value)
		}
		if num < min || num > max {
			return fmt.Errorf("%s must be between %f and %f", value, min, max)
		}
		return nil
	}
}
