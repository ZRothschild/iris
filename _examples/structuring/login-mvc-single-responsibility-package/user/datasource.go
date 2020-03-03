package user

import (
	"errors"
	"sync"
	"time"
)

// IDGenerator将是我们的用户ID生成器，
// 但是在这里，我们按用户ID保持用户的顺序，
// 因此我们将使用可以轻松写入浏览器的数字来从REST API中获取结果
// var IDGenerator = func() string {
// 	return uuid.NewV4().String()
// }

// IDGenerator would be our user ID generator
// but here we keep the order of users by their IDs
// so we will use numbers that can be easly written
// to the browser to get results back from the REST API.
// var IDGenerator = func() string {
// 	return uuid.NewV4().String()
// }

// DataSource是我们的数据存储示例

// DataSource is our data store example.
type DataSource struct {
	Users map[int64]Model
	mu    sync.RWMutex
}

// NewDataSource返回一个新的用户数据源

// NewDataSource returns a new user data source.
func NewDataSource() *DataSource {
	return &DataSource{
		Users: make(map[int64]Model),
	}
}

// GetBy接收一个查询函数，
// 该函数针对我们虚构数据库中的每个单个用户模型触发
// 当该函数返回true时，它将停止迭代

// GetBy receives a query function
// which is fired for every single user model inside
// our imaginary database.
// When that function returns true then it stops the iteration.

//返回查询的返回最后一个已知布尔值
// 和最后一个已知用户模型，以帮助调用者减少loc。

// It returns the query's return last known boolean value
// and the last known user model
// to help callers to reduce the loc.

//但要小心，调用者应始终检查“找到”的内容，
//因为它可能为假，但用户模型中实际上具有真实数据。

// But be carefully, the caller should always check for the "found"
// because it may be false but the user model has actually real data inside it.

//实际上，这是我想到的一个简单但非常聪明的原型函数，
//此后一直在各处使用，希望您也发现它非常有用

// It's actually a simple but very clever prototype function
// I'm think of and using everywhere since then,
// hope you find it very useful too.
func (d *DataSource) GetBy(query func(Model) bool) (user Model, found bool) {
	d.mu.RLock()
	for _, user = range d.Users {
		found = query(user)
		if found {
			break
		}
	}
	d.mu.RUnlock()
	return
}

// GetByID根据其ID返回用户模型

// GetByID returns a user model based on its ID.
func (d *DataSource) GetByID(id int64) (Model, bool) {
	return d.GetBy(func(u Model) bool {
		return u.ID == id
	})
}

// GetByUsername返回基于用户名的用户模型

// GetByUsername returns a user model based on the Username.
func (d *DataSource) GetByUsername(username string) (Model, bool) {
	return d.GetBy(func(u Model) bool {
		return u.Username == username
	})
}

func (d *DataSource) getLastID() (lastID int64) {
	d.mu.RLock()
	for id := range d.Users {
		if id > lastID {
			lastID = id
		}
	}
	d.mu.RUnlock()

	return lastID
}

// InsertOrUpdate将用户添加或更新到内存存储

// InsertOrUpdate adds or updates a user to the (memory) storage.
func (d *DataSource) InsertOrUpdate(user Model) (Model, error) {
	//无论我们将更新update和insert动作的密码哈希值

	// no matter what we will update the password hash
	// for both update and insert actions.
	hashedPassword, err := GeneratePassword(user.password)
	if err != nil {
		return user, err
	}
	user.HashedPassword = hashedPassword

	// update
	if id := user.ID; id > 0 {
		_, found := d.GetByID(id)
		if !found {
			return user, errors.New("ID should be zero or a valid one that maps to an existing User")
		}
		d.mu.Lock()
		d.Users[id] = user
		d.mu.Unlock()
		return user, nil
	}

	// insert
	id := d.getLastID() + 1
	user.ID = id
	d.mu.Lock()
	user.CreatedAt = time.Now()
	d.Users[id] = user
	d.mu.Unlock()

	return user, nil
}
