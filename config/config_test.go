package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tatrasoft/fyp-backend/config"

)

func TestDBConfig_GetConfig_NoErr(t *testing.T) {
	testConfigStruct := config.DBConfig{}
	conf, err := testConfigStruct.GetConfig("testdata/config_test.json")
	require.NoError(t, err)
	require.NotNil(t, conf)

	assert.Equal(t, "test_db", conf.DatabaseName)
}
