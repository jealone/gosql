package gosql

type LoadBalance interface {
	Select(int) int
}

type RoundRobin struct {
}

func (lb *RoundRobin) Select(i int) int {
	panic("implement me")
}

func ProvideRoundRobin() *RoundRobin {
	return &RoundRobin{}
}
