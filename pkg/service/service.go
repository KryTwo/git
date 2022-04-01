package service

type Service struct {
	TodoPeoples
}

type TodoPeoples interface {
	Create()
	GetAll()
	GetById()
	Delete()
	Update()
}

//type Repository struct {    то же самое что и в repository
//	TodoList
//}

//type Service struct {
//	TodoList
//}
//
//func NewService(repos *repository.Repository) *Service {
//	return &Service{
//		TodoList: NewTodoListService(repos.TodoList),
//	}
//}
