package exception

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

func ServiceErrorCatch(tx *gorm.DB, done *context.CancelFunc, funcName string) {
	if err := recover(); err != nil {
		var outErr error
		switch x := (err).(type) {
		case error:
			outErr = x
		case string:
			outErr = fmt.Errorf(x)
		case fmt.Stringer:
			outErr = fmt.Errorf(x.String())
		default:
			outErr = fmt.Errorf("%v", x)
		}
		if tx != nil {
			tx.Rollback()
		}

		panic(StandardRuntimeBadError().
			SetOutPutMessage(outErr.Error()).
			SetFunctionName(funcName).
			SetOriginalError(outErr).SetErrorCode(1))
	} else {
		if tx != nil {
			tx.Commit()
		}
		if done != nil {
			(*done)()
		}
	}
}
