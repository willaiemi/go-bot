package todo

import "fmt"

type Todo struct {
	ID     uint32
	UserID string
	Title  string
	Done   bool
}

func (t Todo) String() string {
	if t.Done {
		return fmt.Sprintf(":white_check_mark: ~~**%s**~~ (ID: %d)", t.Title, t.ID)
	} else {
		return fmt.Sprintf(":black_small_square: **%s** (ID: %d)", t.Title, t.ID)
	}
}

var (
	todos = make(map[string][]Todo) // userID -> list of TO-DO items
)

type ListFilter int

const (
	ListFilterPending ListFilter = iota
	ListFilterCompleted
	ListFilterAll
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

func GetFilteredTodos(userID string, filter ListFilter, highlightTodoId uint32) []Todo {
	userTodos, exists := todos[userID]
	if !exists {
		return []Todo{}
	}

	var filteredTodos []Todo

	switch filter {
	case ListFilterAll:
		return userTodos
	case ListFilterPending:
		for _, todo := range userTodos {
			if todo.ID == highlightTodoId || !todo.Done {
				filteredTodos = append(filteredTodos, todo)
			}
		}
	case ListFilterCompleted:
		for _, todo := range userTodos {
			if todo.Done {
				filteredTodos = append(filteredTodos, todo)
			}
		}
	}

	return filteredTodos
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
