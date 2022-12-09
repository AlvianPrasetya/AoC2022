package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
)

type DisjointSet struct {
	parentMap map[string]string
	rankMap   map[string]int
}

func NewDisjointSet(vList []string) *DisjointSet {
	s := &DisjointSet{
		parentMap: make(map[string]string),
		rankMap:   make(map[string]int),
	}

	for _, v := range vList {
		s.parentMap[v] = v
	}

	return s
}

func (s *DisjointSet) FindRoot(v string) string {
	if v != s.parentMap[v] {
		// Path compression
		s.parentMap[v] = s.FindRoot(s.parentMap[v])
	}

	return s.parentMap[v]
}

func (s *DisjointSet) Merge(v string, w string) {
	vRoot := s.FindRoot(v)
	wRoot := s.FindRoot(w)

	// Union by rank
	if s.rankMap[vRoot] == s.rankMap[wRoot] {
		// rank[v] == rank[w], merge w to v, rank[v]++
		s.parentMap[wRoot] = vRoot
		s.rankMap[vRoot]++
	} else if s.rankMap[vRoot] > s.rankMap[wRoot] {
		// rank[v] > rank[w], merge w to v
		s.parentMap[wRoot] = vRoot
	} else {
		// rank[v] < rank[w], merge v to w
		s.parentMap[vRoot] = wRoot
	}
}

type Group struct {
	GroupID string
	UIDs    []string
}

func main() {
	in, _ := os.Open("in_2.csv")
	defer in.Close()

	// Parse input csv
	records, _ := csv.NewReader(in).ReadAll()
	vList := make([]string, 0, 2*len(records))
	for _, record := range records {
		vList = append(vList, record[0], record[1])
	}

	// Perform UFDS
	s := NewDisjointSet(vList)
	uidMap := make(map[string]bool)
	for _, record := range records {
		uid, adid := record[0], record[1]
		s.Merge(uid, adid)
		uidMap[uid] = true
	}

	// Map UIDs to their groups
	groupMap := make(map[string]*Group)
	for uid := range uidMap {
		groupID := s.FindRoot(uid)

		if _, ok := groupMap[groupID]; !ok {
			groupMap[groupID] = &Group{
				GroupID: groupID,
			}
		}
		groupMap[groupID].UIDs = append(groupMap[groupID].UIDs, uid)
	}

	groupList := make([]*Group, 0, len(groupMap))
	for _, group := range groupMap {
		// Sort UIDs in a group
		sort.Slice(group.UIDs, func(i, j int) bool {
			return group.UIDs[i] < group.UIDs[j]
		})
		groupList = append(groupList, group)
	}

	// Sort groups
	sort.Slice(groupList, func(i, j int) bool {
		return groupList[i].UIDs[0] < groupList[j].UIDs[0]
	})

	fmt.Println(len(groupList))

	// Write output csv
	out, _ := os.Create("out_2.csv")
	defer out.Close()

	w := csv.NewWriter(out)
	for i, group := range groupList {
		if i != 0 {
			// Empty line between groups
			w.Write(nil)
		}
		for _, uid := range group.UIDs {
			w.Write([]string{uid})
		}
	}
	w.Flush()
}
