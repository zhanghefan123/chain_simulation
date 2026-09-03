package executor

import (
	"chain_simulation/configs"
	"fmt"
	"math/rand"
	"sort"
)

// DefaultRandomMaliciousCandidateNodes is kept for backward compatibility.
// It is refreshed from config via RefreshMaliciousCandidatesFromConfig().
var DefaultRandomMaliciousCandidateNodes = []int{2, 3, 5, 6, 8, 9}

// DefaultRandomMaliciousCandidateGroups groups candidate nodes by hop.
// Constraint: each update selects at most one node from each group (hop).
var DefaultRandomMaliciousCandidateGroups = [][]int{{2, 3}, {5, 6}, {8, 9}}

// DefaultCyclicMaliciousNodes is kept for backward compatibility.
// It is refreshed from config via RefreshMaliciousCandidatesFromConfig().
var DefaultCyclicMaliciousNodes = []int{2, 5, 8, 3, 6, 9}

func computeNormalRouterCandidateGroups(numberOfHops, numberOfIntermediateNodes int) [][]int {
	if numberOfHops <= 0 {
		return nil
	}
	if numberOfIntermediateNodes <= 0 {
		numberOfIntermediateNodes = 2
	}
	groups := make([][]int, 0, numberOfHops)
	for hop := 0; hop < numberOfHops; hop++ {
		// Node indexes are 1-based.
		// For hop=0: source=1, intermediates are 2..(1+numberOfIntermediateNodes)
		// For hop=h: source = 1 + h*(numberOfIntermediateNodes+1)
		sourceIndex := 1 + hop*(numberOfIntermediateNodes+1)
		group := make([]int, 0, numberOfIntermediateNodes)
		for i := 0; i < numberOfIntermediateNodes; i++ {
			group = append(group, sourceIndex+1+i)
		}
		groups = append(groups, group)
	}
	return groups
}

func flattenGroups(groups [][]int) []int {
	var out []int
	for _, g := range groups {
		out = append(out, g...)
	}
	return out
}

// computeCyclicMaliciousNodes builds order:
// hop0[i], hop1[i], ... hop(H-1)[i] for i in [0..K-1].
func computeCyclicMaliciousNodes(groups [][]int) []int {
	if len(groups) == 0 {
		return nil
	}
	maxLen := 0
	for _, g := range groups {
		if len(g) > maxLen {
			maxLen = len(g)
		}
	}
	out := make([]int, 0, len(groups)*maxLen)
	for i := 0; i < maxLen; i++ {
		for hop := 0; hop < len(groups); hop++ {
			if i < len(groups[hop]) {
				out = append(out, groups[hop][i])
			}
		}
	}
	return out
}

// allocateBadNodesPerGroup spreads badNodeCount across hop groups as evenly as possible.
// e.g. 3 groups: 3->[1,1,1], 4->[2,1,1], 5->[2,2,1], 6->[2,2,2].
func allocateBadNodesPerGroup(badNodeCount, numGroups int) []int {
	allocations := make([]int, numGroups)
	base := badNodeCount / numGroups
	remainder := badNodeCount % numGroups
	for i := 0; i < numGroups; i++ {
		allocations[i] = base
		if i < remainder {
			allocations[i]++
		}
	}
	return allocations
}

// MaxRandomMaliciousBadNodeCount returns the max bad nodes per update: sum of (groupSize-1) per hop.
// Constraint: cannot attack every normal router in the same hop at once.
func MaxRandomMaliciousBadNodeCount(groups [][]int) int {
	total := 0
	for _, group := range groups {
		if len(group) > 0 {
			total += len(group) - 1
		}
	}
	return total
}

func validateAllocationFeasible(allocations []int, groups [][]int) error {
	if len(allocations) != len(groups) {
		return fmt.Errorf("allocation length %d != group count %d", len(allocations), len(groups))
	}
	slotCounts := append([]int(nil), allocations...)
	groupCaps := make([]int, len(groups))
	for i, group := range groups {
		groupCaps[i] = len(group) - 1
	}
	sort.Sort(sort.Reverse(sort.IntSlice(slotCounts)))
	sort.Sort(sort.Reverse(sort.IntSlice(groupCaps)))
	for i, need := range slotCounts {
		if need > groupCaps[i] {
			return fmt.Errorf(
				"cannot place %d attacked nodes into a hop with at most %d attackable normal routers",
				need,
				groupCaps[i],
			)
		}
	}
	return nil
}

// ValidateRandomMaliciousBadNodeCount checks count against topology-derived candidate groups.
func ValidateRandomMaliciousBadNodeCount(count int, groups [][]int) error {
	if count < 1 {
		return fmt.Errorf("bad-node-count must be >= 1, got %d", count)
	}
	if len(groups) == 0 {
		return fmt.Errorf("no malicious candidate groups configured")
	}
	maxCount := MaxRandomMaliciousBadNodeCount(groups)
	if count > maxCount {
		return fmt.Errorf(
			"bad-node-count=%d exceeds max=%d (cannot attack all normal routers in any hop)",
			count,
			maxCount,
		)
	}
	allocations := allocateBadNodesPerGroup(count, len(groups))
	if err := validateAllocationFeasible(allocations, groups); err != nil {
		return fmt.Errorf("bad-node-count=%d is not feasible: %w", count, err)
	}
	return nil
}

// pickRandomBadNodes picks badNodeCount distinct normal routers using per-hop allocation.
func pickRandomBadNodes(rng *rand.Rand, groups [][]int, badNodeCount int) []int {
	numGroups := len(groups)
	allocations := allocateBadNodesPerGroup(badNodeCount, numGroups)
	perm := rng.Perm(numGroups)

	badNodes := make([]int, 0, badNodeCount)
	for slot, groupIdx := range perm {
		attackCount := allocations[slot]
		if attackCount == 0 {
			continue
		}
		group := groups[groupIdx]
		pickedIndexes := rng.Perm(len(group))[:attackCount]
		for _, idx := range pickedIndexes {
			badNodes = append(badNodes, group[idx])
		}
	}
	sort.Ints(badNodes)
	return badNodes
}

// RefreshMaliciousCandidatesFromConfig refreshes candidate nodes/groups according to current loaded config.
// If config is not initialized or values are invalid, it keeps previous defaults.
func RefreshMaliciousCandidatesFromConfig() {
	hops := configs.TopConfigInstance.SecPathMabConfig.NumberOfHops
	k := configs.TopConfigInstance.SecPathMabConfig.NumberOfIntermediateNodes
	if hops <= 0 {
		return
	}
	groups := computeNormalRouterCandidateGroups(hops, k)
	if len(groups) == 0 {
		return
	}
	DefaultRandomMaliciousCandidateGroups = groups
	DefaultRandomMaliciousCandidateNodes = flattenGroups(groups)
	DefaultCyclicMaliciousNodes = computeCyclicMaliciousNodes(groups)
}

