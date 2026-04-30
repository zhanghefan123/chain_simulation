package entities

type ModifyBloomFilter struct {
	BfEffectiveBits int `json:"bf_effective_bits"`
}

func NewModifyBloomFilter(bfEffectiveBits int) *ModifyBloomFilter {
	return &ModifyBloomFilter{
		BfEffectiveBits: bfEffectiveBits,
	}
}
