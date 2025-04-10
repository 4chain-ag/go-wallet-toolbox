package primitives

// ReqHistoryNote is the history representation of the request
type ReqHistoryNote struct {
	When        *string
	What        string
	ExtraFields map[string]interface{}
}
