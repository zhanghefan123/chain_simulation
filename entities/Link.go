package entities

type Link struct {
	SourceNode Node   `json:"source_node"`
	TargetNode Node   `json:"target_node"`
	LinkType   string `json:"link_type"`
}
