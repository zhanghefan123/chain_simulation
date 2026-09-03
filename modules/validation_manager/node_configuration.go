package validation_manager

import "chain_simulation/entities"

// ModifyBloomFilter updates Bloom filter parameters on one validation node.
func ModifyBloomFilter(nodeIndex int, params *entities.ModifyBloomFilter) error {
	return postToNode(nodeIndex, modifyBloomFilterPath, params, "modify Bloom filter")
}

// InsertSessionTableEntries preloads session table entries on one validation node.
func InsertSessionTableEntries(nodeIndex int, params *entities.InsertSessionTableEntriesInstance) error {
	return postToNode(
		nodeIndex,
		insertSessionTableEntriesPath,
		params,
		"insert session table entries",
	)
}
