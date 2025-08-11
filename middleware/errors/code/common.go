package code

const (
	// JSONDecodeErrCode request body decode error. err: %s
	JSONDecodeErrCode int32 = 1000

	// MapperActionErrCode  mapper error. action: %s, err: %s
	MapperActionErrCode int32 = 1002
	// FileNotFoundErrCode file not found. file name: %s
	FileNotFoundErrCode int32 = 1003
	// RawErrWrapErrCode raw error wrap: %v
	RawErrWrapErrCode int32 = 1004
)
