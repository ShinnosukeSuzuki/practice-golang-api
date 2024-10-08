package apperrors

// 独自エラー構造体を定義
type MyAppError struct {
	ErrCode
	Message string
	Err     error `json:"-"` // エラー情報はJSONに含めない
}

// エラーメソッドを定義
func (myErr *MyAppError) Error() string {
	return myErr.Err.Error()
}

func (myErr *MyAppError) Unwrap() error {
	return myErr.Err
}
