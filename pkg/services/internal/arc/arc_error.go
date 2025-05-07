package arc

import "fmt"

// APIError represents an error returned by the ARC API when status code is 4xx.
type APIError struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	Status    int    `json:"status"`
	Detail    string `json:"detail"`
	Instance  string `json:"instance"`
	TxID      string `json:"txid"`
	ExtraInfo string `json:"extraInfo"`
}

// Error returns the error string it's the implementation of the error interface.
func (a *APIError) Error() string {
	if a.IsEmpty() {
		return "ARC error: empty (or not in json) response"
	}
	return fmt.Sprintf("ARC error: %s <txID: %s> %s", a.Title, a.TxID, a.Detail)
}

// IsEmpty checks if the error is empty indicating that we could not parse the error response.
func (a *APIError) IsEmpty() bool {
	return a == nil || a.Status == 0
}
