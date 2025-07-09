package todo

type Todo struct {
	ID     uint32
	UserID string
	Title  string
	Done   bool
}

var (
	todos = make(map[string][]Todo) // userID -> list of TO-DO items
)

func AddTodo(userID, title string) (Todo, error) {
	if _, exists := todos[userID]; !exists {
		todos[userID] = []Todo{}
	}

	itemID := uint32(len(todos[userID]) + 1) // Simple ID generation

	todo := Todo{
		ID:     itemID,
		UserID: userID,
		Title:  title,
		Done:   false,
	}

	todos[userID] = append(todos[userID], todo)

	return todo, nil
}

func GetTodos(userID string) []Todo {
	if _, exists := todos[userID]; !exists {
		return []Todo{} // Return empty slice if no todos exist for the user
	}

	return todos[userID]
}
