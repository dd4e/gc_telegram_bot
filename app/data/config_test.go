package data

import "testing"

var testConfig = Configs{
	12345: &Chat{
		Timeout:      60,
		Enabled:      true,
		ChatTitle:    "testChat",
		TimeoutLimit: 65400,
		ChatID:       12345,
	},
}

func TestConfigs_Get(t *testing.T) {
	_, err := testConfig.Get(56789)
	if err == nil {
		t.Error("error with get unknown chat")
	}

	_, err = testConfig.Get(12345)
	if err != nil {
		t.Error("error with getting chat")
	}
}

func TestConfigs_Exist(t *testing.T) {
	if testConfig.Exist(56789) {
		t.Error("error with exist unknown chat")
	}

	if !testConfig.Exist(12345) {
		t.Error("error with exist chat by id")
	}
}

func TestConfigs_ExistAndEnable(t *testing.T) {
	if !testConfig.ExistAndEnable(12345) {
		t.Error("error exist and enable")
	}
}
