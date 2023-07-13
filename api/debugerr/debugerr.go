package debugerr

import "fmt"

func WrapMsg(msg string, where string) string {
	return fmt.Sprintf("couldn't %s at %s", msg, where)
}
