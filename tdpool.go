package main

// TDPool is a traffic director pool info
type TDPool struct {
	Name   string
	URL    string
	Status string
}

// offline first
type byStatus []TDPool

func (b byStatus) Len() int {
	return len(b)
}

func (b byStatus) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b byStatus) Less(i, j int) bool {
	return b[i].Status == "offline" || false
}
