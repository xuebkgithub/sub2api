package db

import (
	"errors"
	"strings"

	"github.com/lib/pq"
)

// IsUniqueConstraintViolation 判断错误是否为唯一约束冲突。
//
// 支持多种检测方式：
//  1. PostgreSQL 特定错误码 23505（唯一约束冲突）
//  2. 错误消息中包含的通用关键词
//
// 这种多层次的检测确保了对不同数据库驱动和 ORM 的兼容性。
func IsUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
	}

	// 优先检测 PostgreSQL 特定错误码（最精确）。
	// 错误码 23505 对应 unique_violation。
	// 参考：https://www.postgresql.org/docs/current/errcodes-appendix.html
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	// 回退到错误消息检测（兼容其他场景）。
	// 这些关键词覆盖了 PostgreSQL、MySQL 等主流数据库的错误消息。
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "unique constraint") ||
		strings.Contains(msg, "duplicate entry")
}
