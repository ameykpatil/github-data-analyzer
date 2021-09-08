package repo

// Repo encapsulates required properties related to repo
type Repo struct {
	ID             string
	Name           string
	CommitCount    int
	EventTypeCount map[string]int
}

// base heap structure for Repo
type baseRepoHeap []Repo

// Swap is required to swap the elements of the heap (Sort interface)
func (r baseRepoHeap) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Len is required to know the length (Sort interface)
func (r baseRepoHeap) Len() int {
	return len(r)
}

// Push adds an element to heap (Heap interface)
func (r *baseRepoHeap) Push(x interface{}) {
	*r = append(*r, x.(Repo))
}

// Pop take out the top element from the heap (Heap interface)
func (r *baseRepoHeap) Pop() interface{} {
	old := *r
	n := len(old)
	x := old[n-1]
	*r = old[0 : n-1]
	return x
}

// concrete heap structure for repo
type repoHeap struct {
	baseRepoHeap
	less func(i, j Repo) bool
}

// Less defines the way to determine lesser element (Sort interface)
func (r repoHeap) Less(i, j int) bool {
	return r.less(r.baseRepoHeap[i], r.baseRepoHeap[j])
}
