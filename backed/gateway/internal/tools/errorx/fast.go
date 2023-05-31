package errorx

func NewPermsNotExist(args ...interface{}) *CodeError {
	if len(args) == 0 {
		return New(Perms, "permission dose not exist")
	} else if len(args) == 1 {
		return NewCodeMsgF(Perms, "%s permission dose not exist", args)
	}
	return NewCodeMsgF(Perms, "%s permission %s dose not exist", args)
}

func NewPermsExist(args ...interface{}) *CodeError {
	if len(args) == 0 {
		return New(Perms, "permission is exist")
	} else if len(args) == 1 {
		return NewCodeMsgF(Perms, "%s permission is exist", args)
	}
	return NewCodeMsgF(Perms, "%s permission %s is exist", args)
}

func NewPermsAddFail(args ...interface{}) *CodeError {
	return NewCodeMsgF(Perms, "%s permission add fail", args)
}
