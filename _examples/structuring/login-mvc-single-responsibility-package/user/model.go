package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

//模型是我们的用户示例模型

// Model is our User example model.
type Model struct {
	ID        int64  `json:"id"`
	Firstname string `json:"firstname"`
	Username  string `json:"username"`

	//密码是客户端提供的密码，不会存储在服务器中的任何地方
	//它仅用于注册和更新密码之类的操作，
	//因为我们接受`DataSource#InsertOrUpdate`函数中的Model实例

	// password is the client-given password
	// which will not be stored anywhere in the server.
	// It's here only for actions like registration and update password,
	// because we caccept a Model instance
	// inside the `DataSource#InsertOrUpdate` function.
	password       string
	HashedPassword []byte    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

// GeneratePassword将根据用户输入为我们生成一个哈希密码

// GeneratePassword will generate a hashed password for us based on the
// user's input.
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}

// ValidatePassword将检查密码是否匹配

// ValidatePassword will check if passwords are matched.
func ValidatePassword(userPassword string, hashed []byte) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(userPassword)); err != nil {
		return false, err
	}
	return true, nil
}
