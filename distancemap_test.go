package clustering

import "testing"

func TestDistanceMapClustering1(t *testing.T) {
	d := NewDistanceMapClusterSet(nil)
	if d == nil {
		t.Errorf("could not create empty DistanceMapClusterSet")
	}

	d = NewDistanceMapClusterSet(DistanceMap{"a": {"b": 0.0}})
	if d == nil {
		t.Errorf("could not create 2-node DistanceMapClusterSet")
	}
	if d.Count() != 2 {
		t.Errorf("2-node DistanceMapClusterSet doesn't start with 2 clusters")
	}
	n1, n2 := 0, 0
	d.EachCluster(-1, func(cluster int) {
		n1++

		d.EachItem(cluster, func(x ClusterItem) {
			n2++
		})
	})
	if n1 != 2 || n2 != 2 {
		t.Errorf("2-node DistanceMapClusterSet didn't enumerate 2 clusters w/start=-1")
	}
	n1, n2 = 0, 0
	d.EachCluster(0, func(cluster int) {
		n1++

		d.EachItem(cluster, func(x ClusterItem) {
			n2++
		})
	})
	if n1 != 1 || n2 != 1 {
		t.Errorf("2-node DistanceMapClusterSet didn't enumerate 1 clusters w/start=0")
	}
	n1 = 0
	d.EachCluster(1, func(cluster int) {
		n1++
	})
	if n1 != 0 {
		t.Errorf("2-node DistanceMapClusterSet didn't enumerate 0 clusters w/start=1")
	}

	if d.Distance(0, 1, "a", "b") != 0.0 {
		t.Errorf("2-node DistanceMapClusterSet gave wrong distance")
	}

	d.Merge(0, 1)

	if d.Count() != 1 {
		t.Errorf("after Merge(0,1), 2-node DistanceMapClusterSet isn't 1 cluster")
	}
	n1, n2 = 0, 0
	d.EachCluster(-1, func(cluster int) {
		n1++

		d.EachItem(cluster, func(x ClusterItem) {
			n2++
		})
	})
	if n1 != 1 || n2 != 2 {
		t.Errorf("after Merge(0,1), 2-node DistanceMapClusterSet isn't 1 cluster with 2 items")
	}
}

func TestDistanceMapClustering2(t *testing.T) {
	d := NewDistanceMapClusterSet(DistanceMap{"a": {"b": 0.0, "c": 0.0}})
	if d == nil {
		t.Errorf("could not create 3-node DistanceMapClusterSet")
	}
	if d.Count() != 3 {
		t.Errorf("3-node DistanceMapClusterSet doesn't start with 3 clusters")
	}
	n1, n2 := 0, 0
	d.EachCluster(-1, func(cluster int) {
		n1++

		d.EachItem(cluster, func(x ClusterItem) {
			n2++
		})
	})
	if n1 != 3 || n2 != 3 {
		t.Errorf("3-node DistanceMapClusterSet didn't enumerate 3 clusters w/start=-1")
	}

	Cluster(d, Threshold(1.0), CompleteLinkage())

	if d.Count() != 1 {
		t.Errorf("after clustering, 3-node DistanceMapClusterSet isn't 1 cluster")
	}
	n1, n2 = 0, 0
	d.EachCluster(-1, func(cluster int) {
		n1++

		d.EachItem(cluster, func(x ClusterItem) {
			n2++
		})
	})
	if n1 != 1 || n2 != 3 {
		t.Errorf("after clustering, 3-node DistanceMapClusterSet isn't 1 cluster with 3 items")
	}
}

func TestDistanceMapClustering3(t *testing.T) {
	d := NewDistanceMapClusterSet(DistanceMap{
		"a": {"b": 0.0, "c": 0.0, "d": 1.0, "e": 0.4},
		"b": {"c": 0.1, "d": 0.9, "e": 0.4},
		"c": {"d": 0.9, "e": 0.2},
		"d": {"e": 0.1},
	})
	if d == nil {
		t.Errorf("could not create 5-node DistanceMapClusterSet")
	}
	if d.Count() != 5 {
		t.Errorf("5-node DistanceMapClusterSet doesn't start with 5 clusters")
	}
	n1, n2 := 0, 0
	d.EachCluster(-1, func(cluster int) {
		n1++

		d.EachItem(cluster, func(x ClusterItem) {
			n2++
		})
	})
	if n1 != 5 || n2 != 5 {
		t.Errorf("5-node DistanceMapClusterSet didn't enumerate 5 clusters w/start=-1")
	}

	Cluster(d, Threshold(0.4), CompleteLinkage())

	if d.Count() != 2 {
		t.Errorf("after clustering, 5-node DistanceMapClusterSet isn't 2 clusters")
	}
	n1, n2 = 0, 0
	c0, c1 := 0, 0
	d.EachCluster(-1, func(cluster int) {
		n1++

		d.EachItem(cluster, func(x ClusterItem) {
			n2++

			if cluster == 0 {
				c0++
			}
			if cluster == 1 {
				c1++
			}
		})
	})
	if n1 != 2 || n2 != 5 {
		t.Errorf("after clustering, 5-node DistanceMapClusterSet isn't 2 clusters with 5 items")
	}
	if (c0 < 2 || c0 > 3) || (c1 < 2 || c1 > 3) || (c0+c1) != 5 {
		t.Errorf("after clustering, 5-node DistanceMapClusterSet should be 2,3")
	}
}
