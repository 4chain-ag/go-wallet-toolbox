package history

const (
	InternalizeActionHistoryNote = "internalizeAction"
	ProcessActionHistoryNote     = "processAction"
)

func UserIDHistoryAttr(userID int) map[string]any {
	return map[string]any{
		"userId": userID,
	}
}
