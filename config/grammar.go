package config

type Config struct {
	Populations []*Population `json:"populations"`
	Jobs        []*Job        `json:"jobs"`
}

type Population struct {
	Name          string  `json:"name"`
	Population    int     `json:"population"`
	Crossover     string  `json:"xover"`
	Mutate        float64 `json:"mutationRate"`
	CrossoverRate float64 `json:"xoverRate"`
	Epsilon       int     `json:"epsilon"`
	Checkpoint    bool    `json:"checkpoint"`
	Seed          int     `json:"seed"`
}

type Job struct {
	Name       string `json:"name"`
	GraphFile  string `json:"graph"`
	OutputDir  string `json:"out"`
	Population string `json:"population"`
}
