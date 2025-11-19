package network

import (
	"bytes"
	"net"
	"sort"
)

// Deduplicate removes duplicate networks from a slice of *net.IPNet.
func Deduplicate(networks []*net.IPNet) []*net.IPNet {
	if len(networks) < 2 {
		return networks
	}

	seen := make(map[string]struct{})
	uniqueNetworks := make([]*net.IPNet, 0, len(networks))

	for _, n := range networks {
		// Use the string representation of the CIDR as the key
		cidrStr := n.String()
		if _, found := seen[cidrStr]; !found {
			seen[cidrStr] = struct{}{}
			uniqueNetworks = append(uniqueNetworks, n)
		}
	}

	return uniqueNetworks
}

// Aggregate removes networks that are completely contained within larger networks in the slice.
func Aggregate(networks []*net.IPNet) []*net.IPNet {
	if len(networks) < 2 {
		return networks
	}

	// Sort by prefix length (ascending), then by IP address.
	// This ensures larger networks (smaller prefix) come first.
	sort.Slice(networks, func(i, j int) bool {
		onesI, _ := networks[i].Mask.Size()
		onesJ, _ := networks[j].Mask.Size()
		if onesI != onesJ {
			return onesI < onesJ
		}
		return bytes.Compare(networks[i].IP, networks[j].IP) < 0
	})

	finalList := make([]*net.IPNet, 0, len(networks))
	for _, n := range networks {
		isContained := false
		// Check if the current network `n` is contained by any network already in `finalList`.
		for _, existing := range finalList {
			// If an existing network (which is guaranteed to be a larger or equal-sized prefix)
			// contains the IP of the current network, then `n` is redundant.
			if existing.Contains(n.IP) {
				isContained = true
				break
			}
		}
		if !isContained {
			finalList = append(finalList, n)
		}
	}

	return finalList
}

// Merge combines adjacent networks into larger supernets.
// This process is repeated until no more merges can be performed.
func Merge(networks []*net.IPNet) []*net.IPNet {
	if len(networks) < 2 {
		return networks
	}

	for {
		mergedOnce := false

		// Sort by IP address to find adjacent networks easily.
		sort.Slice(networks, func(i, j int) bool {
			return bytes.Compare(networks[i].IP, networks[j].IP) < 0
		})

		var nextNetworks []*net.IPNet
		skipNext := false

		for i := 0; i < len(networks); i++ {
			if skipNext {
				skipNext = false
				continue
			}

			if i+1 >= len(networks) {
				nextNetworks = append(nextNetworks, networks[i])
				break
			}

			n1 := networks[i]
			n2 := networks[i+1]

			if supernet := tryMerge(n1, n2); supernet != nil {
				nextNetworks = append(nextNetworks, supernet)
				mergedOnce = true
				skipNext = true // n2 was merged, so skip it in the next loop iteration.
			} else {
				nextNetworks = append(nextNetworks, n1)
			}
		}

		networks = nextNetworks

		if !mergedOnce {
			break // No merges in a full pass, so we are done.
		}
	}
	return networks
}

// tryMerge attempts to merge two networks. If they are adjacent and can be
// combined into a single larger network, the new supernet is returned.
// Otherwise, it returns nil.
func tryMerge(n1, n2 *net.IPNet) *net.IPNet {
	ones, bits := n1.Mask.Size()
	if o, b := n2.Mask.Size(); o != ones || b != bits || ones == 0 {
		return nil // Must be same type and have same prefix length.
	}

	// The supernet would have a prefix one smaller.
	supernetPrefix := ones - 1
	supernetMask := net.CIDRMask(supernetPrefix, bits)

	// For two networks to be mergeable, they must both belong to the same potential
	// supernet. The first network's IP must also be the starting IP of the supernet.
	if !n1.IP.Mask(supernetMask).Equal(n1.IP) {
		return nil
	}

	// The second network's IP must be the one that starts the second half
	// of the supernet. We can calculate this expected IP.
	expectedSecondHalfIP := make(net.IP, len(n1.IP))
	copy(expectedSecondHalfIP, n1.IP)
	byteIndex := supernetPrefix / 8
	bitIndex := 7 - (supernetPrefix % 8)
	expectedSecondHalfIP[byteIndex] |= (1 << bitIndex)

	if n2.IP.Equal(expectedSecondHalfIP) {
		return &net.IPNet{IP: n1.IP, Mask: supernetMask}
	}

	return nil
}
