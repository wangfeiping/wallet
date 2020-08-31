package types

type SeedOutput struct {
	Seed string `json:"seed"`
}

type Adapter interface {
	CreateSeed() string
}
