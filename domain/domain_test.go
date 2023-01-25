package domain

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	// 测试时使用, 减少加密密码的时间
	enctyptPasswordCost = bcrypt.MinCost
	m.Run()
}
