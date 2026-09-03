package topology

type RatioDistribution struct {
	Start float64 `json:"start"` // 起始值
	End   float64 `json:"end"`   // 结束值
}

type MaliciousParams struct {
	CorruptRatio        RatioDistribution `json:"corrupt_ratio"`
	CorruptSpecialRatio RatioDistribution `json:"corrupt_special_ratio"`
}

type SpecialParams struct {
	InnerRouterType int             `json:"inner_router_type"` // 内部的类型 (即这是一个非可信的路由器还是一个可信的路径验证路由器)
	MaliciousParams MaliciousParams `json:"malicious_params"`  // 恶意参数
}

type Node struct {
	Index         int           `json:"index"`
	Type          string        `json:"type"`
	X             float64       `json:"x"`
	Y             float64       `json:"y"`
	SpecialParmas SpecialParams `json:"special_params"`
}
