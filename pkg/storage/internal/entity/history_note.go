package entity

var InternalizeActionHistoryNote = "internalizeAction"

func UserIDHistoryAttr(userID int) map[string]any {
	return map[string]any{
		"userId": userID,
	}
}
