package pokeapi

type LocationAreas struct {
	Next     *string              `json:"next"`
	Previous *string              `json:"previous"`
	Results  []LocationAreaResult `json:"results"`
}

type LocationAreaResult struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type PokemonInfoResponse struct {
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func (resp *PokemonInfoResponse) ConvertToDTO() PokemonDTO {
	var p PokemonDTO
	p.Name = resp.Name
	p.Height = resp.Height
	p.Weight = resp.Weight
	p.Experience = resp.BaseExperience

	for _, stat := range resp.Stats {
		p.Stats = append(p.Stats, struct {
			Name   string
			Effort int
		}{
			Name:   stat.Stat.Name,
			Effort: stat.Effort,
		})
	}

	for _, t := range resp.Types {
		p.Types = append(p.Types, struct {
			Name string
			Slot int
		}{
			Name: t.Type.Name,
			Slot: t.Slot,
		})
	}

	return p
}

type PokemonDTO struct {
	Name   string
	Height int
	Weight int
	Stats  []struct {
		Name   string
		Effort int
	}
	Types []struct {
		Name string
		Slot int
	}
	Experience int
}
