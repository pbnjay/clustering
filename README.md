# clustering
Some basic clustering algorithm implementations for Go.

[![GoDoc](https://godoc.org/github.com/pbnjay/clustering?status.svg)](https://godoc.org/github.com/pbnjay/clustering)

# Quick Start
To cluster a simple set of data using a map of maps to distances, and
complete-linkage hierarchical clustering with a simple threshold cutoff, the
following code suffices:

    // NB map can be asymmetric like this, both key orderings are checked if necessary
    clusters := clustering.NewDistanceMapClusterSet(clustering.DistanceMap{
      "a": {"b": 0.0, "c": 0.0, "d": 1.0, "e": 0.4},
      "b": {"c": 0.1, "d": 0.9, "e": 0.4},
      "c": {"d": 0.9, "e": 0.2},
      "d": {"e": 0.1},
    })
    clustering.Cluster(clusters, clustering.Threshold(0.4), clustering.CompleteLinkage())

    // Enumerate clusters and print members
    clusters.EachCluster(-1, func(cluster int) {
      clusters.EachItem(cluster, func(x clustering.ClusterItem) {
        fmt.Println(cluster, x)
      }
    }

 Outputs two clusters (ordering may be different due to map enumeration):

    0 d
    0 e
    1 a
    1 b
    1 c

# Supported Data sources

I highly recommend implementing the [`ClusterSet` interface](http://godoc.org/github.com/pbnjay/clustering#ClusterSet) to work with your existing data, it will be much more efficient and give you better tools to tweak things. For smaller data sets, using the included [`DistanceMap`](http://godoc.org/github.com/pbnjay/clustering#DistanceMap) is probably good enough for most purposes.

# Supported Hierarchical Clustering Linkage methods

* **Complete Linkage (Maximum Linkage)** - Uses the maximum distance between any 2 items in the 2 clusters as the cluster-pair's linkage score. i.e. the 2 clusters with the smallest distance between the *furthest* two items are selected.

* **Single Linkage (Minimum Linkage)** - Uses the minimum distance between any 2 items in the 2 clusters as the cluster-pair's linkage score. i.e. the 2 clusters with the smallest distance between the *closest* two items are selected.

* **Average Linkage (UPGMA)** - Uses the average distance between all pairs of items in the 2 clusters as the cluster-pair's linkage score. i.e. the 2 clusters with the smallest average distance across all pairs of items are selected.

## License and Contributions

This code is available under the MIT license. Contributions are welcome if following the [standard Go style conventions](https://github.com/golang/go/wiki/CodeReviewComments).