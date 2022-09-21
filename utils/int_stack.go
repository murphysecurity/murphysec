package utils

type IntStack []int

func NewIntStack() *IntStack {
	r := make(IntStack, 0)
	return &r
}

func (i *IntStack) Len() int {
	return len(*i)
}

func (i *IntStack) Push(el int) {
	*i = append(*i, el)
}

func (i *IntStack) Peek() int {
	return (*i)[i.Len()-1]
}

func (i *IntStack) Pop() int {
	r := i.Peek()
	*i = (*i)[:i.Len()-1]
	return r
}

func (i *IntStack) Empty() bool {
	return i.Len() == 0
}
