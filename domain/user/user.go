package user

// User encapsulates required properties related to user
type User struct {
	ID             string
	Username       string
	CommitCount    int
	EventTypeCount map[string]int
}

// base heap structure for User
type baseUserHeap []User

// Swap is required to swap the elements of the heap (Sort interface)
func (r baseUserHeap) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Len is required to know the length (Sort interface)
func (r baseUserHeap) Len() int {
	return len(r)
}

// Push adds an element to heap (Heap interface)
func (r *baseUserHeap) Push(x interface{}) {
	*r = append(*r, x.(User))
}

// Pop take out the top element from the heap (Heap interface)
func (r *baseUserHeap) Pop() interface{} {
	old := *r
	n := len(old)
	x := old[n-1]
	*r = old[0 : n-1]
	return x
}

// concrete heap structure for user
type userHeap struct {
	baseUserHeap
	less func(i, j User) bool
}

// Less defines the way to determine lesser element (Sort interface)
func (r userHeap) Less(i, j int) bool {
	return r.less(r.baseUserHeap[i], r.baseUserHeap[j])
}
