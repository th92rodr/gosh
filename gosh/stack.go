package main

type stack []string

func (t *terminal) newStack() {
	t.history = make([]string, 0)
}

func (t *terminal) push(input string) {
	t.history = append(t.history, input)
}

func (t *terminal) top() string {
	return t.history[len(t.history)-1]
}
