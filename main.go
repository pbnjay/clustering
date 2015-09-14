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

// MergeNext finds the next pair of clusters to merge by applying the linkage
// method to all pairs and selecting the best result. It then verifies criteria
// are met before merging them. It returns true if the pair of clusters was
// merged successfully, otherwise false.
func (h *HClustering) MergeNext() bool {
	bestScore := math.MaxFloat64
	var bestPair []int

	ocs, ok := h.ClusterSet.(OptimizedClusterSet)
	if !ok {
		ocs = &defaultOptimizedClusterSet{cs: h.ClusterSet}
	}

	// TODO: memoize pair scores so we only update what changed
	// instead of re-calculating everything every time
	h.ClusterSet.EachCluster(-1, func(c1 int) {
		h.ClusterSet.EachCluster(c1, func(c2 int) {
			h.LinkageType.Reset()

			h.ClusterSet.EachItem(c1, func(a ClusterItem) {
				ocs.EachItemDistance(c1, c2, a, func(b ClusterItem, dist float64) {
					h.LinkageType.Put(a, b, dist)
				})
			})

			score := h.LinkageType.Get()
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

	h.ClusterSet.Merge(bestPair[0], bestPair[1])
	return true
}
