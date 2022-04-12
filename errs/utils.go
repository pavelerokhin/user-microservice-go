package errs

import "fmt"

func FErrJSON(msg string) string {
	return fmt.Sprintf("{\"error\":\"%s\"}", msg)
}
