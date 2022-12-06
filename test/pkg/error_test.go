package test

import (
	"BRGS/pkg/e"
	"fmt"
	"testing"
)

func TestErrTranslate(t *testing.T) {
	fmt.Printf("e.TranslateError(e.ERROR_FUNCTION): %v\n", e.TranslateError(e.ERROR_FUNCTION))
	fmt.Printf("e.TranslateError(e.ERROR_TRANSLATE): %v\n", e.TranslateError(e.ERROR_TRANSLATE))
	fmt.Printf("e.TranslateToError(e.ERROR_CREATE): %v\n", e.TranslateToError(e.ERROR_CREATE))
}
