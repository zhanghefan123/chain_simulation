package entities

type Link struct {
	SourceNode Node `json:"source_node"`
	TargetNode Node `json:"target_node"`
}
