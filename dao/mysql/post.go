package mysql

import "bluebell/models"

// CreatePost 创建帖子插入MySQL数据库
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
	post_id, author_id,community_id,title, content)
	values (?, ?, ?, ?, ?)
	`
	_, err = db.Exec(sqlStr, p.ID, p.AuthorID, p.CommunityID, p.Title, p.Content)
	return
}

// GetPostById 根据帖子ID查询MySQL数据库
func GetPostById(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select
	post_id, title, content, author_id, community_id, create_time
	from post
	where post_id = ?
	`
	err = db.Get(post, sqlStr, pid)
	return
}

// GetUserByID 根据id获取用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id=?`
	err = db.Get(user, sqlStr, uid)
	return
}

// GetPostList 根据page,size获取分页列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select 
	post_id, title, content, author_id, community_id, create_time
	from post
	limit ?,?
	`
	posts = make([]*models.Post, 0, 2) // 不要写成make([]*models.Post, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}
