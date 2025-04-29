package wdk

const (
	// BasketNameForChange is the name of the output basket that is used to store "change" outputs
	BasketNameForChange = "default"

	// StorageCommissionPurpose is the purpose-string used for tagging storage commission outputs
	StorageCommissionPurpose = "storage-commission"

	// ChangePurpose is the purpose-string used for tagging change outputs
	ChangePurpose = "change"

	// NumberOfDesiredUTXOsForChange is the number of desired UTXOs for the change output basket,
	// it influences the number of change outputs created during createAction
	NumberOfDesiredUTXOsForChange = 32

	// MinimumDesiredUTXOValue is the minimum value of UTXOs in the change output basket,
	// it influences the number of change outputs created during createAction
	MinimumDesiredUTXOValue = 1000

	// DefaultNumberOfDesiredUTXOs is the default number of desired UTXOs for non-change output baskets.
	// During createAction or internalizeAction, when a user provides an output with non-existing basket name, it will be created with this number of desired UTXOs.
	DefaultNumberOfDesiredUTXOs = 0 // TODO: Even though in TS version this is set to 0, we should double check that if this is intentional or not

	// DefaultMinimumDesiredUTXOValue is the default minimum value of UTXOs for non-change output baskets.
	// During createAction or internalizeAction, when a user provides an output with non-existing basket name, it will be created with this minimum value.
	DefaultMinimumDesiredUTXOValue = 0 // TODO: Even though in TS version this is set to 0, we should double check that if this is intentional or not
)
