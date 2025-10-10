package entities

type Node struct {
	Index int     `json:"index"`
	Type  string  `json:"type"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}
