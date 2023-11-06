package document

import "time"

type Todo struct {
	Task      string
	Completed bool
	Version   int
	Deleted   bool
	UpdatedAt time.Time
}

// IDを含んだデータ型
type TodoDataset struct {
	ID        int       `json:"id"`
	Task      string    `json:"task"`
	Completed bool      `json:"completed"`
	Deleted   bool      `json:"deleted"`
	Version   int       `json:"version"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Todo型をTodoDatasetに変換する
func ConvertTodoDataset(id int, todo Todo) TodoDataset {
	return TodoDataset{
		ID:        id,
		Task:      todo.Task,
		Completed: todo.Completed,
		Deleted:   todo.Deleted,
		Version:   todo.Version,
		UpdatedAt: todo.UpdatedAt,
	}
}

func ConvertTodo(todoDataset TodoDataset) Todo {
	return Todo{
		Task:      todoDataset.Task,
		Completed: todoDataset.Completed,
		Deleted:   todoDataset.Deleted,
		Version:   todoDataset.Version,
		UpdatedAt: todoDataset.UpdatedAt,
	}
}
