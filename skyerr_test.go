package skyerr

import (
	"errors"
	"testing"
)

func TestSkyError(t *testing.T) {
	err := SkyError("test error")
	t.Log(err.String())
	err = SkyError(nil)
	if err == nil {
		t.Log("err is nil", err.String())
	}
	err = SkyError(errors.New("json marshal wrong"))
	t.Log(err.String())
}

func TestSkyErrorM(t *testing.T) {
	err := SkyErrorM(-1001, "test error %d", 9999)
	t.Log(err.String())
}

func MyName() SkyErr {
	return SkyErrorF(-1005, errors.New("json marshal wrong"), "test throw %d", 9999)
}

func UseMyName() SkyErr {
	err := MyName()
	if err != nil {
		return SkyError(err)
	}
	return nil
}

func TestSkyThrowF(t *testing.T) {
	err := UseMyName()
	if err != nil {
		err = SkyError(err)
	}
	t.Log(err.String())
}
