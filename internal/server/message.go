package server

const (
	addItemCommand     = "AddItem"
	removeItemCommand  = "RemoveItem"
	getItemCommand     = "GetItem"
	getAllItemsCommand = "GetAllItems"
)

type Message struct {
	Command string `json:"command"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}
