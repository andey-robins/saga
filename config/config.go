package config

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
)

func ParseConfig(configFile string) *Config {
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var config Config
	json.Unmarshal(byteValue, &config)

	if err := config.Validate(); err != nil {
		log.Fatalln("invalid config file:", err)
		panic(err)
	}

	return &config
}

// Validate returns an error if the config is invalid and nil if the config is valid.
func (c *Config) Validate() error {

	validCrossover := []string{"default"}

	in := func(s string, ss []string) bool {
		for _, v := range ss {
			if v == s {
				return true
			}
		}
		return false
	}

	populationNames := make(map[string]bool)
	for _, p := range c.Populations {
		if _, ok := populationNames[p.Name]; ok {
			return errors.New("duplicate population name: " + p.Name)
		}

		populationNames[p.Name] = true

		if p.Population < 4 || p.Population > 10_000 || p.Population%2 != 0 {
			return errors.New("invalid population size: " + p.Name)
		}

		if !in(p.Crossover, validCrossover) {
			return errors.New("invalid crossover method: " + p.Crossover)
		}

		if p.MutationRate < 0 || p.MutationRate > 1 {
			return errors.New("invalid mutation rate: " + p.Name)
		}

		if p.CrossoverRate < 0 || p.CrossoverRate > 1 {
			return errors.New("invalid crossover rate: " + p.Name)
		}

		if p.Epsilon < 0 || p.Epsilon > 1_000_000 {
			return errors.New("invalid epsilon: " + p.Name)
		}
	}

	jobNames := make(map[string]bool)
	for _, j := range c.Jobs {
		if _, ok := jobNames[j.Name]; ok {
			return errors.New("duplicate job name: " + j.Name)
		}

		if _, ok := populationNames[j.Population]; !ok {
			return errors.New("invalid population name: " + j.Population)
		}
	}

	return nil
}
