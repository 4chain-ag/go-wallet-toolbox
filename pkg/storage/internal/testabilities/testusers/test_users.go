package testusers

type User struct {
	// Name of the user just for in tests logging purpose
	Name    string
	ID      int
	PrivKey string
}

var Alice = User{
	Name:    "Alice",
	ID:      1,
	PrivKey: "143ab18a84d3b25e1a13cefa90038411e5d2014590a2a4a57263d1593c8dee1c",
}

var Bob = User{
	Name:    "Bob",
	ID:      2,
	PrivKey: "0881208859876fc227d71bfb8b91814462c5164b6fee27e614798f6e85d2547d",
}
