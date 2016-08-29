package handler

import (
	"encoding/json"

	"github.com/go-macaron/binding"
)

func requestErrBytes(files []string, err error) (int, []byte) {
	reqErrs := binding.Errors{
		binding.Error{
			FieldNames:     files,
			Classification: "DeserializationError",
			Message:        err.Error(),
		},
	}
	errOutput, _ := json.Marshal(reqErrs)
	return 400, errOutput
}
