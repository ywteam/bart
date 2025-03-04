package ydk


type YdkError struct {
	Code    int
	Message string	
}
type YdkErrorAggregator struct {
	Errors []YdkError
}
func (e *YdkErrorAggregator) Has() bool {
	return len(e.Errors) > 0
}
func (e *YdkErrorAggregator) Error(code int, message string) {
	e.Errors = append(e.Errors, YdkError{Code: code, Message: message})
}