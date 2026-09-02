package executor

// randomMaliciousBadNodeCount controls how many nodes are attacked in each update when corrupt-ratio-mode=random.
// Nodes are spread across hop groups (e.g. 3 hops: 4->[2,1,1], 5->[2,2,1], 6->[2,2,2]) without emptying any hop.
var randomMaliciousBadNodeCount = 1

func SetRandomMaliciousBadNodeCount(count int) error {
	RefreshMaliciousCandidatesFromConfig()
	if err := ValidateRandomMaliciousBadNodeCount(count, DefaultRandomMaliciousCandidateGroups); err != nil {
		return err
	}
	randomMaliciousBadNodeCount = count
	return nil
}

func GetMaxRandomMaliciousBadNodeCount() int {
	RefreshMaliciousCandidatesFromConfig()
	return MaxRandomMaliciousBadNodeCount(DefaultRandomMaliciousCandidateGroups)
}

func GetRandomMaliciousBadNodeCount() int {
	return randomMaliciousBadNodeCount
}
