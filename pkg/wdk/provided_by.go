package wdk

type ProvidedBy string

const (
	ProvidedByYou           ProvidedBy = "you"
	ProvidedByStorage       ProvidedBy = "storage"
	ProvidedByYouAndStorage ProvidedBy = "you-and-storage"
)
