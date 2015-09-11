// Package hclust provides a set of Go interfaces and methods to quickly
// implement hierarchichal clustering using simple data types.
//
// To cluster a simple set of data using a map of maps to distances, and
// complete-linkage clustering with a simple threshold cutoff, the following
// code suffices:
//
//    // NB map does not have to be symmetric, both key orders are checked
//    clusters := hclust.NewDistanceMapClusterSet(hclust.DistanceMap{
//      "a": {"b": 0.0, "c": 0.0, "d": 1.0, "e": 0.4},
//      "b": {"c": 0.1, "d": 0.9, "e": 0.4},
//      "c": {"d": 0.9, "e": 0.2},
//      "d": {"e": 0.1},
//    })
//    hclust.Cluster(clusters, hclust.Threshold(0.4), hclust.CompleteLinkage())
//
//    clusters.EachCluster(-1, func(cluster int) {
//      clusters.EachItem(cluster, func(x hclust.ClusterItem) {
//        fmt.Println(cluster, x)
//      }
//    }
//
// Outputs two clusters (ordering may be different due to map enumeration):
//
//    0 d
//    0 e
//    1 a
//    1 b
//    1 c
//
// Other useful linkage types that should be implemented one day:
//   Centroid  -- select clusters where the "centers" are close together.
//   Ward      -- select clusters that reduce the variance of distances.
//
package hclust
