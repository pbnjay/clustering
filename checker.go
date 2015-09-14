package clustering

import "log"

// Checker implements the decision criteria used to stop clustering.
// Note that this interface may also be used to collect the hierarchical
// clustering tree produced by the agglomerations.
type Checker interface {
	// Check decides wether or not to continue merging cluster nodes, based on
	// the current set of clusters and the next best merge score.
	//
	// Returns true to continue clustering, false to stop.
	Check(clusters ClusterSet, i, j int, nextScore float64) bool
}

// MaxClusters returns a Checker that limits total number of output clusters.
func MaxClusters(t int) Checker {
	return limitClustersCount{t}
}

// Threshold returns a Checker that stops before a merge threshold is passed.
func Threshold(t float64) Checker {
	return simpleThreshold{t}
}

// TreeLog prints the merge decisions that occur at each step of the tree.
func TreeLog(c Checker) Checker {
	return clusterTreeLog{c}
}

/////////////

type simpleThreshold struct {
	val float64
}

func (t simpleThreshold) Check(clusters ClusterSet, i, j int, nextScore float64) bool {
	return nextScore <= t.val
}

/////////////

type clusterTreeLog struct {
	chk Checker
}

func (c clusterTreeLog) Check(clusters ClusterSet, i, j int, nextScore float64) bool {
	t := c.chk.Check(clusters, i, j, nextScore)
	if t {
		log.Printf("  merge (%d,%d) ~~ %f %v", i, j, nextScore, clusters)
	} else {
		log.Printf("  STOP  (%d,%d) ~~ %f", i, j, nextScore)
	}
	return t
}

//////////////

type limitClustersCount struct {
	val int
}

func (t limitClustersCount) Check(clusters ClusterSet, i, j int, nextScore float64) bool {
	return clusters.Count() > t.val
}
