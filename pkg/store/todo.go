package store

// キーバリューストア型 IDをキーとする
type TodoKeyValueStore map[int]*Todo

type Todo struct {
	Task      string
	Completed bool
	Version   int
	Deleted   bool
}

// IDを含んだデータ型
type TodoDateset struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
	Deleted   bool   `json:"deleted"`
	Version   int    `json:"version"`
}

// Todo型をTodoDatasetに変換する
func ConvertTodoDataset(id int, todo Todo) TodoDateset {
	return TodoDateset{
		ID:        id,
		Task:      todo.Task,
		Completed: todo.Completed,
		Deleted:   todo.Deleted,
		Version:   todo.Version,
	}
}
