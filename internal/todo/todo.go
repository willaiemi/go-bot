package todo

import "fmt"

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

func MarkTodoDone(userID string, itemID uint32) (Todo, error) {
	if _, exists := todos[userID]; !exists {
		return Todo{}, fmt.Errorf("no to-do items found, create one with `/add`")
	}

	for i, todo := range todos[userID] {
		if todo.ID == itemID {
			todos[userID][i].Done = true
			return todos[userID][i], nil
		}
	}

	return Todo{}, fmt.Errorf("TO-DO item with (ID: %d) does not exist", itemID)
}
