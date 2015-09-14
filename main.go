package hclust

import "math"

// ClusterItem represents a generic cluster item key. For implementation
// purposes, it should be comparable / suitable as a map key.
type ClusterItem interface{}

// ClusterSet is implemented by the user to allow generic clustering data
// sources. Clusters are identified by simple integers, and items within
// clusters are identified by the generic ClusterItem interface. Paired item
// distances are computed by the user code as well.
type ClusterSet interface {
	// Count returns the number of clusters in the set.
	Count() int

	// EachCluster enumerates every cluster id "after" start. Use start=-1 to
	// start enumeration from the beginning.
	EachCluster(start int, cb func(cluster int))

	// EachItem enumerates every item from the cluster.
	EachItem(cluster int, cb func(item ClusterItem))

	// Distance computes the distance between two items in separate clusters.
	Distance(c1, c2 int, item1, item2 ClusterItem) float64

	// Merge the two clusters together. After this step Count() should be
	// reduced by 1. Retuns the cluster that is merged into (kept) and the
	// cluster that is swapped into the place of the merged cluster (typically
	// the last cluster).
	Merge(cluster1, cluster2 int) (kept, swappedIn int)
}

// OptimizedClusterSet allows implementors to optimize distance calculations by
// caching the left-hand cluster/item data. This interface is optional.
type OptimizedClusterSet interface {
	// EachItemDistance is a way for optimized distance calculations, by
	// caching information about (c1, item1) before applying to another cluster.
	// It is equivalent to the following code:
	//    cs.EachItem(c1, func(item1 ClusterItem){
	//      cs.EachCluster(c1, func(c2 int){
	//
	//        // equiv to: cs.EachItemDistance(c1,c2,item1,callback)
	//
	//        cs.EachItem(c2, func(item2 ClusterItem){
	//          dist := cs.Distance(c1,c2,item1,item2)
	//			callback(item2, dist)
	//        })
	//      })
	//
	//    })
	EachItemDistance(c1, c2 int, item1 ClusterItem, cb func(item2 ClusterItem, dist float64))
}

type defaultOptimizedClusterSet struct {
	cs ClusterSet
}

func (x *defaultOptimizedClusterSet) EachItemDistance(c1, c2 int, item1 ClusterItem, cb func(ClusterItem, float64)) {
	x.cs.EachItem(c2, func(item2 ClusterItem) {
		dist := x.cs.Distance(c1, c2, item1, item2)
		cb(item2, dist)
	})
}

// HClustering is a hierarchical clustering wrapper for arbitrary data sets.
type HClustering struct {
	// LinkageType is the method used to select clusters to merge.
	LinkageType LinkageType

	// Checker is used to check stop criteria for the clustering.
	Checker Checker

	// ClusterSet is used to enumerate and manipulate the set of clusters.
	ClusterSet ClusterSet

	lwCache   []float64
	distCache map[int]map[int]float64
}

//////////////////

// Cluster clusters the input set (in-place) using the specified linkage type
// until the provided threshold is hit.
func Cluster(c ClusterSet, chk Checker, lt LinkageType) {
	h := HClustering{
		ClusterSet:  c,
		Checker:     chk,
		LinkageType: lt,
	}

	for h.ClusterSet.Count() > 1 {
		if !h.MergeNext() {
			break
		}
	}
}

// calculate the distance between cluster i and cluster j.
// also caches and reuses prior calculations
func (h *HClustering) dist(i, j int) float64 {
	if h.distCache != nil {
		if i > j {
			i, j = j, i
		}
		if _, f := h.distCache[i]; f {
			if s, f2 := h.distCache[i][j]; f2 {
				return s
			}
		} else {
			// prep for saving to cache
			h.distCache[i] = make(map[int]float64)
		}
	}
	h.LinkageType.Reset()

	ocs, ok := h.ClusterSet.(OptimizedClusterSet)
	if !ok {
		ocs = &defaultOptimizedClusterSet{cs: h.ClusterSet}
	}

	h.ClusterSet.EachItem(i, func(a ClusterItem) {
		ocs.EachItemDistance(i, j, a, func(b ClusterItem, dist float64) {
			h.LinkageType.Put(a, b, dist)
		})
	})

	s := h.LinkageType.Get()
	if h.distCache != nil {
		h.distCache[i][j] = s
	}
	return s
}

// merges clusters i and j, and calculates the new distances resulting from it.
// 1) call ClusterSet.Merge(i,j)
// 2) remove cluster r=ClusterSet.Count() from distance cache
// 3) for each cluster x:
// 3a) update distances for (i,j) and remove r
func (h *HClustering) mergeAndUpdateAll(i, j int) {
	nc := h.ClusterSet.Count()

	diks := []float64{}
	djks := []float64{}
	for k := 0; k < nc; k++ {
		diks = append(diks, h.dist(i, k))
		djks = append(djks, h.dist(j, k))
	}

	origDist := diks[j]
	ni, nj := h.ClusterSet.Merge(i, j)

	if nj != j {
		// where did nj go to?
		r := j
		if ni == j {
			r = i
		}

		//move cached distances from nj into r
		for k := 0; k < nc; k++ {
			if k == nj {
				continue
			}
			x1, y1 := k, r
			if x1 > y1 {
				x1, y1 = r, k
			}

			x2, y2 := k, nj
			if x2 > y2 {
				x2, y2 = nj, k
			}
			h.distCache[x1][y1] = h.distCache[x2][y2]
		}

		// now remove unused cache values
		for k := 0; k < nc; k++ {
			if k == nj {
				delete(h.distCache, nj)
				continue
			}
			if _, f := h.distCache[k]; f {
				delete(h.distCache[k], nj)
			}
		}
	}

	// apply lance-williams update method to all affected pairs
	nc--
	for k := 0; k < nc; k++ {
		dik := diks[k]
		djk := djks[k]
		dd := dik - djk
		if dd < 0.0 {
			dd = -dd
		}

		d := h.lwCache[0]*dik + h.lwCache[1]*djk + h.lwCache[2]*origDist + h.lwCache[3]*dd
		if ni < k {
			h.distCache[ni][k] = d
		} else {
			h.distCache[k][ni] = d
		}
	}

}

// MergeNext finds the next pair of clusters to merge by applying the linkage
// method to all pairs and selecting the best result. It then verifies criteria
// are met before merging them. It returns true if the pair of clusters was
// merged successfully, otherwise false.
func (h *HClustering) MergeNext() bool {
	bestScore := math.MaxFloat64
	var bestPair []int

	if len(h.lwCache) != 4 {
		h.lwCache = h.LinkageType.LWParams()
		h.distCache = make(map[int]map[int]float64)
	}

	h.ClusterSet.EachCluster(-1, func(c1 int) {
		h.ClusterSet.EachCluster(c1, func(c2 int) {
			score := h.dist(c1, c2)
			if score < bestScore {
				bestScore = score
				bestPair = []int{c1, c2}
			}
		})
	})

	if len(bestPair) == 0 || bestScore == math.MaxFloat64 {
		return false
	}

	if !h.Checker.Check(h.ClusterSet, bestPair[0], bestPair[1], bestScore) {
		return false
	}

	if h.distCache == nil {
		h.ClusterSet.Merge(bestPair[0], bestPair[1])
	} else {
		h.mergeAndUpdateAll(bestPair[0], bestPair[1])
	}
	return true
}
