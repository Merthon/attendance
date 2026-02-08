package service

import (
	"attendance-system/internal/database"
	"attendance-system/internal/model"
	"errors"
	"time"
	"fmt" 
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 定义加密密钥
var jwtKey = []byte("my_secret_key_123456")

type Claims struct {
	UserID uint   `json:"userId"`
	Role   int    `json:"role"`
	jwt.RegisteredClaims
}


// Login 登录业务
func Login(username, password string) (string, *model.User, error) {
	var user model.User
	// 1. 找用户
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", nil, errors.New("用户不存在")
	}
	fmt.Printf(">> 正在尝试登录: %s\n", username)
    fmt.Printf(">> 数据库里的密码(Hash): %s\n", user.Password)
    fmt.Printf(">> 前端传来的密码: %s\n", password)

	// 2. 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("密码错误")
	}

	// 3. 生成 JWT Token
	expirationTime := time.Now().Add(24 * time.Hour) // 24小时有效
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, &user, err
}

// Register 注册业务
func Register(user *model.User) error {
	fmt.Printf(">> [注册阶段] 用户名: %s, 原始密码: %s\n", user.Username, user.Password)
	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return database.DB.Create(user).Error
}