package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/classic/crypto/ed25519"
	"github.com/tendermint/go-amino-x"

	sdk "github.com/tendermint/classic/sdk/types"
	authtypes "github.com/tendermint/classic/sdk/x/auth/types"
)

var (
	priv = ed25519.GenPrivKey()
	addr = sdk.AccAddress(priv.PubKey().Address())
)

func TestParseQueryResponse(t *testing.T) {
	sdkResBytes := amino.MustMarshalBinaryLengthPrefixed(sdk.Result{GasUsed: 10})
	gas, err := parseQueryResponse(sdkResBytes)
	assert.Equal(t, gas, uint64(10))
	assert.Nil(t, err)
	gas, err = parseQueryResponse([]byte("fuzzy"))
	assert.Equal(t, gas, uint64(0))
	assert.Error(t, err)
}

func TestCalculateGas(t *testing.T) {
	makeQueryFunc := func(gasUsed uint64, wantErr bool) func(string, []byte) ([]byte, int64, error) {
		return func(string, []byte) ([]byte, int64, error) {
			if wantErr {
				return nil, 0, errors.New("")
			}
			return amino.MustMarshalBinaryLengthPrefixed(sdk.Result{GasUsed: gasUsed}), 0, nil
		}
	}
	type args struct {
		queryFuncGasUsed uint64
		queryFuncWantErr bool
		adjustment       float64
	}
	tests := []struct {
		name         string
		args         args
		wantEstimate uint64
		wantAdjusted uint64
		wantErr      bool
	}{
		{"error", args{0, true, 1.2}, 0, 0, true},
		{"adjusted gas", args{10, false, 1.2}, 10, 12, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queryFunc := makeQueryFunc(tt.args.queryFuncGasUsed, tt.args.queryFuncWantErr)
			gotEstimate, gotAdjusted, err := CalculateGas(queryFunc, []byte(""), tt.args.adjustment)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, gotEstimate, tt.wantEstimate)
			assert.Equal(t, gotAdjusted, tt.wantAdjusted)
		})
	}
}

func TestDefaultTxEncoder(t *testing.T) {

	defaultEncoder := authtypes.DefaultTxEncoder()
	encoder := GetTxEncoder()

	compareEncoders(t, defaultEncoder, encoder)
}

func TestConfiguredTxEncoder(t *testing.T) {

	customEncoder := func(tx sdk.Tx) ([]byte, error) {
		return json.Marshal(tx)
	}

	config := sdk.GetConfig()
	config.SetTxEncoder(customEncoder)

	encoder := GetTxEncoder()

	compareEncoders(t, customEncoder, encoder)
}

func TestReadStdTxFromFile(t *testing.T) {

	// Build a test transaction
	fee := authtypes.NewStdFee(50000, sdk.Coins{sdk.NewInt64Coin("atom", 150)})
	stdTx := authtypes.NewStdTx([]sdk.Msg{}, fee, []authtypes.StdSignature{}, "foomemo")

	// Write it to the file
	encodedTx, _ := amino.MarshalJSON(stdTx)
	jsonTxFile := writeToNewTempFile(t, string(encodedTx))
	defer os.Remove(jsonTxFile.Name())

	// Read it back
	decodedTx, err := ReadStdTxFromFile(jsonTxFile.Name())
	require.NoError(t, err)
	require.Equal(t, decodedTx.Memo, "foomemo")
}

func compareEncoders(t *testing.T, expected sdk.TxEncoder, actual sdk.TxEncoder) {
	msgs := []sdk.Msg{sdk.NewTestMsg(addr)}
	tx := authtypes.NewStdTx(msgs, authtypes.StdFee{}, []authtypes.StdSignature{}, "")

	defaultEncoderBytes, err := expected(tx)
	require.NoError(t, err)
	encoderBytes, err := actual(tx)
	require.NoError(t, err)
	require.Equal(t, defaultEncoderBytes, encoderBytes)
}

func writeToNewTempFile(t *testing.T, data string) *os.File {
	fp, err := ioutil.TempFile(os.TempDir(), "client_tx_test")
	require.NoError(t, err)

	_, err = fp.WriteString(data)
	require.NoError(t, err)

	return fp
}
