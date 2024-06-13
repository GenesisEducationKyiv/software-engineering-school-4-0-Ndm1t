package database

type CustomError interface {
	int
	Error() string
}

type DBErorr struct {
	StatusCode int
}

func (e *DBErorr) Error() string {
	return "Database Error"
}

type AlreadySubscribed struct {
	StatusCode int
}

func (e *AlreadySubscribed) Error() string {
	return "Email is already subscribed"
}
