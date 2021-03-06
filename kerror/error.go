package kerror

import (
	"github.com/robertkrimen/otto"
	"fmt"
)

func ToError(err error) error {
	if oerr, ok := err.(*otto.Error); ok {
		return fmt.Errorf("%s", oerr.String())
	} else if oerr, ok := err.(otto.Error); ok {
		return fmt.Errorf("%s", oerr.String())
	} else {
		return err
	}
}
