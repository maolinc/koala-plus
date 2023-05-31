package errorx

const (
	ERROR uint32 = 1001
	Param uint32 = 1004
	Db    uint32 = 1005
	Token uint32 = 1006
	User  uint32 = 1007
	Perms uint32 = 1008
	App   uint32 = 1009
)

var (
	SystemError        *CodeError
	RequestParamError  *CodeError
	DbError            *CodeError
	TokenGenError      *CodeError
	UserNotExistError  *CodeError
	PermsNotEmptyError *CodeError
	AppNotExistError   *CodeError
)

func init() {
	SystemError = New(ERROR, "System error, please try again later")
	DbError = New(Db, "The database is busy, please try again later")
	TokenGenError = New(Token, "Token invalid, please log in again")
	UserNotExistError = New(User, "User does not exist")
	PermsNotEmptyError = New(Perms, "Permission does not empty")
	AppNotExistError = New(App, "App does not exist")
	RequestParamError = New(Param, "Request error")
}
