package domain

import "time"

// User 领域对象, 是 DDD 中的 entity
// 也称为 BO (Business Object)
type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}

//type Address struct {
//}
