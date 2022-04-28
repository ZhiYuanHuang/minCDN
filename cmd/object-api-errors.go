package cmd

type OperationTimedOut struct {
}

func (e OperationTimedOut) Error() string {
	return "Operation timed out"
}
