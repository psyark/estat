package test

import (
	"encoding/json"
	"fmt"

	"github.com/wI2L/jsondiff"
)

func cyclicTest(data []byte, cycle func() ([]byte, error)) error {
	untyped := map[string]any{}
	if err := json.Unmarshal(data, &untyped); err != nil {
		return err
	}

	data1, _ := json.MarshalIndent(untyped, "", "  ")
	data2, err := cycle()
	if err != nil {
		return err
	}

	patch, err := jsondiff.CompareJSON(data1, data2)
	if err != nil {
		return err
	}

	for _, op := range patch {
		_, _ = fmt.Printf("%v %v: value=%v, oldValue=%v\n", op.Type, op.Path, op.Value, op.OldValue)
	}

	if len(patch) != 0 {
		return fmt.Errorf("unmatch")
	}
	return nil
}
