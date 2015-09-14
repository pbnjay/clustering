package clustering

// LinkageType is an interface that defines how two clusters are scored
// based on the pairwise distances of their items.
type LinkageType interface {
	// Reset clears the internal state for this linkage type.
	Reset()

	// Put adds a new distance observation for the item-pair.
	Put(item1, item2 ClusterItem, dist float64)

	// Get returns the current value of the linkage based on all the observed
	// values so far.
	Get() float64

	// LWParams returns the lance-williams parameters for updating clusters
	// after a merge. Return (alpha_i, alpha_j, Beta, gamma), if return value
	// is not 4 floats, then clustering falls back to recomputing at each pass.
	LWParams() []float64
}

// CompleteLinkage implements complete-linkage clustering, which is defined as
// the maximum distance between any pair of items from the two clusters.
func CompleteLinkage() LinkageType {
	return &maxLinkage{}
}

// SingleLinkage implements single-linkage clustering, which is defined as
// the minimum distance between any pair of items from the two clusters.
func SingleLinkage() LinkageType {
	return &minLinkage{}
}

// AverageLinkage implements average-linkage (sometimes referred to as UPGMA)
// clustering, which is defined as the average of all distances between all
// pairs of items in the two clusters.
func AverageLinkage() LinkageType {
	return &avgLinkage{}
}

// WeightedAverageLinkage implements WPGMA linkage agglomeration method
// clustering, which is defined as the average of all distances between pairs
// of items in the two clusters. It weights clusters equally regardless of
// number of items.
func WeightedAverageLinkage() LinkageType {
	return &avgLinkage{isWeighted: true}
}

////////////////

type maxLinkage struct {
	maxDist float64
}

func (c *maxLinkage) Reset() {
	c.maxDist = -1.0
}

func (c *maxLinkage) Get() float64 {
	return c.maxDist
}

func (c *maxLinkage) Put(a, b ClusterItem, dist float64) {
	if dist > c.maxDist || c.maxDist < 0.0 {
		c.maxDist = dist
	}
}

func (c *maxLinkage) LWParams() []float64 {
	return []float64{0.5, 0.5, 0.0, 0.5}
}

////////////////

type minLinkage struct {
	minDist float64
}

func (c *minLinkage) Reset() {
	c.minDist = -1.0
}

func (c *minLinkage) Get() float64 {
	return c.minDist
}

func (c *minLinkage) Put(a, b ClusterItem, dist float64) {
	if dist < c.minDist || c.minDist < 0.0 {
		c.minDist = dist
	}
}

func (c *minLinkage) LWParams() []float64 {
	return []float64{0.5, 0.5, 0.0, -0.5}
}

////////////////

type avgLinkage struct {
	avgDist    float64
	totalPairs float64

	isWeighted  bool
	leftCounts  map[ClusterItem]struct{}
	rightCounts map[ClusterItem]struct{}
}

func (c *avgLinkage) Reset() {
	c.avgDist = 0.0
	c.totalPairs = 0.0
	if !c.isWeighted {
		c.leftCounts = make(map[ClusterItem]struct{})
		c.rightCounts = make(map[ClusterItem]struct{})
	}
}

func (c *avgLinkage) Get() float64 {
	if c.totalPairs <= 0.0 {
		return 0.0
	}
	return c.avgDist / c.totalPairs
}

func (c *avgLinkage) Put(a, b ClusterItem, dist float64) {
	c.avgDist += dist
	c.totalPairs++
	if !c.isWeighted {
		c.leftCounts[a] = struct{}{}
		c.rightCounts[b] = struct{}{}
	}
}

func (c *avgLinkage) LWParams() []float64 {
	if c.isWeighted {
		return []float64{0.5, 0.5, 0.0, 0.0}
	}
	ni := float64(len(c.leftCounts))
	nj := float64(len(c.rightCounts))
	return []float64{ni / (ni + nj), nj / (ni + nj), 0.0, 0.0}
}
