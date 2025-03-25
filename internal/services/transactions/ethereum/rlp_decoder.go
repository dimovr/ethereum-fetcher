package ethereum

import (
	"encoding/hex"
	"strings"

	errorstypes "ethereum_fetcher/internal/services/transactions/types"

	"github.com/ethereum/go-ethereum/rlp"
)

func DecodeHashes(rlpHex string) ([]string, error) {
	decoders := []Decoder{
		decodeWith(StringDecoder{}),
		decodeWith(StringArrayDecoder{}),
	}

	for _, decoder := range decoders {
		hashes, err := decoder(rlpHex)
		if err == nil {
			return hashes, nil
		}
	}

	return nil, errorstypes.InvalidRlpEncoding
}

func validateAndFormatHash(hashHex string) ([]string, error) {
	hashHex = strings.TrimPrefix(hashHex, "0x")
	if len(hashHex) != 40 && len(hashHex) != 64 {
		return nil, errorstypes.InvalidHashLength
	}
	return []string{"0x" + hashHex}, nil
}

type HashValue interface {
	~string | ~[]string
}

type RlpDecoder[T HashValue] interface {
	Decode(rlpEncoded string) (T, error)
	Validate(value T) ([]string, error)
}

type GenericDecoder[T HashValue] struct{}

type StringDecoder = GenericDecoder[string]
type StringArrayDecoder = GenericDecoder[[]string]

func (d GenericDecoder[T]) Decode(rlpEncoded string) (T, error) {
	rlpBytes, err := hex.DecodeString(rlpEncoded)
	if err != nil {
		var zero T
		return zero, errorstypes.InvalidHexEncoding
	}

	var value T
	err = rlp.DecodeBytes(rlpBytes, &value)
	if err != nil {
		var zero T
		return zero, errorstypes.InvalidRlpEncoding
	}

	return value, nil
}

func (d GenericDecoder[T]) Validate(value T) ([]string, error) {
	switch v := any(value).(type) {
	case string:
		return validateAndFormatHash(v)
	case []string:
		result := make([]string, 0, len(v))
		for _, str := range v {
			validated, err := validateAndFormatHash(str)
			if err != nil {
				return nil, err
			}
			result = append(result, validated...)
		}
		return result, nil
	default:
		return nil, errorstypes.InvalidRlpEncoding
	}
}

type Decoder func(string) ([]string, error)

func decodeWith[T HashValue](decoder RlpDecoder[T]) Decoder {
	return func(rlpEncoded string) ([]string, error) {
		rlpEncoded = strings.TrimPrefix(rlpEncoded, "0x")
		value, err := decoder.Decode(rlpEncoded)
		if err != nil {
			return nil, err
		}
		return decoder.Validate(value)
	}
}
