package DAO

type NotGood struct {
	S string
}
type AgeLimit struct {
	S string
}
type TooFastRequest struct {
	S string
}

func (i *TooFastRequest) Error() string {
	return i.S
}
func (i *NotGood) Error() string {
	return i.S
}
func (i *AgeLimit) Error() string {
	return i.S
}
