package primitives

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type mockStruct struct {
	Tx ExplicitByteArray `json:"tx"`
}

func TestExplicitByteArrayUnmarshall(t *testing.T) {
	// given:
	jValue := `{
		"tx": [1, 2, 3, 255]
	}`

	// when:
	var jVal mockStruct
	err := json.Unmarshal([]byte(jValue), &jVal)

	// then:
	require.NoError(t, err)

	assert.Equal(t, 4, len(jVal.Tx))
	assert.Equal(t, byte(1), jVal.Tx[0])
	assert.Equal(t, byte(2), jVal.Tx[1])
	assert.Equal(t, byte(3), jVal.Tx[2])
	assert.Equal(t, byte(255), jVal.Tx[3])
}

func TestExplicitByteArrayUnmarshallOutOfRange(t *testing.T) {
	// given:
	jValue := `{
		"tx": [256]
	}`

	// when:
	var jVal mockStruct
	err := json.Unmarshal([]byte(jValue), &jVal)

	// then:
	require.Error(t, err)
}

func TestExplicitByteArrayMarshall(t *testing.T) {
	// given:
	jVal := mockStruct{
		Tx: ExplicitByteArray{1, 2, 3, 255},
	}

	// when:
	marshalled, err := json.Marshal(jVal)

	// then:
	require.NoError(t, err)

	assert.Equal(t, `{"tx":[1,2,3,255]}`, string(marshalled))
}
