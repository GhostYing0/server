package e

const (
	SUCCESS                        = 200
	ERROR_AUTH_CHECK_TOKEN_EMPTY   = 10000
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 10001
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 10002
	ERROR_AUTH_TOKEN_PARSE         = 10003
	ERROR_AUTH_TOKEN_DIFF          = 10004
)

const (
	Pass       = 1
	Reject     = 2
	Processing = 3
	Revoked    = 4
)

const (
	EnrollOpen  = 1
	EnrollClose = 2
)

const (
	StudentRole = 1
	TeacherRole = 2
)
