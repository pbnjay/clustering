package clustering

// DistanceMap is a map of maps from cluster items (pairs) to the distance
// measures between them. A distance map does not have to be symmetric, but it
// is highly recommended to have all pairs defined.
type DistanceMap map[ClusterItem]map[ClusterItem]float64

type distMapClusterSet struct {
	data map[ClusterItem]map[ClusterItem]float64

	clusters [][]ClusterItem
}

// NewDistanceMapClusterSet initializes a new ClusterSet from a distance map by
// creating a singleton cluster for every unique item in the maps.
func NewDistanceMapClusterSet(data DistanceMap) ClusterSet {
	d := &distMapClusterSet{
		data: data,
	}

	allItems := make(map[ClusterItem]struct{})
	for k1, subs := range data {
		if _, done := allItems[k1]; !done {
			allItems[k1] = struct{}{}
			d.clusters = append(d.clusters, []ClusterItem{k1})
		}
		for k2 := range subs {
			if _, done := allItems[k2]; !done {
				allItems[k2] = struct{}{}
				d.clusters = append(d.clusters, []ClusterItem{k2})
			}
		}
	}

	return d
}

func (d *distMapClusterSet) EachCluster(start int, cb func(cluster int)) {
	if start+1 >= len(d.clusters) {
		return
	}

	for i := start + 1; i < len(d.clusters); i++ {
		cb(i)
	}
}

func (d *distMapClusterSet) EachItem(cluster int, cb func(ClusterItem)) {
	for _, x := range d.clusters[cluster] {
		cb(x)
	}
}

func (d *distMapClusterSet) Distance(c1, c2 int, item1, item2 ClusterItem) float64 {
	if x, ok := d.data[item1]; ok {
		if y, ok := x[item2]; ok {
			return y
		}
	}
	if x, ok := d.data[item2]; ok {
		if y, ok := x[item1]; ok {
			return y
		}
	}
	return 1.0
}

func (d *distMapClusterSet) Count() int {
	return len(d.clusters)
}

func (d *distMapClusterSet) Merge(i, j int) (keep, swappedIn int) {
	if j < i {
		j, i = i, j
	}

	// move the to-be-merged cluster to the end of the array
	x := len(d.clusters) - 1
	if j < x {
		d.clusters[x], d.clusters[j] = d.clusters[j], d.clusters[x]
		j = x
	}
	d.clusters[i] = append(d.clusters[i], d.clusters[j]...)
	d.clusters = d.clusters[:j]
	return i, x
}
