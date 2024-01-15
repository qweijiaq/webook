package repository

type UserRepository struct {
}

func (r *UserRepository) FindById(int642 int64) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 找到了会写 cache
}
