package ethereum

import (
	"testing"

	"ethereum_fetcher/internal/services/transactions/ethereum"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeHashes(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedHashes []string
		expectError    bool
	}{
		{
			name:  "Valid single hash RLP with 0x prefix",
			input: "0xb842307834383630336637616466663766626663326131306232326136373130333331656536386632653464316364373361353834643537633838323164663739333536",
			expectedHashes: []string{
				"0x48603f7adff7fbfc2a10b22a6710331ee68f2e4d1cd73a584d57c8821df79356",
			},
			expectError: false,
		},
		{
			name:  "Valid list of hashes RLP with 0x prefix",
			input: "f90110b842307866633262336236646233386135316462336239636239356465323962373139646538646562393936333036323665346234623939646630353666666237663265b842307834383630336637616466663766626663326131306232326136373130333331656536386632653464316364373361353834643537633838323164663739333536b842307863626339323065376262383963626362353430613436396131363232366266313035373832353238336162386561633366343564303038313165656638613634b842307836643630346666633634346132383266636138636238653737386531653366383234356438626431643439333236653330313661336338373862613063626264",
			expectedHashes: []string{
				"0xfc2b3b6db38a51db3b9cb95de29b719de8deb99630626e4b4b99df056ffb7f2e",
				"0x48603f7adff7fbfc2a10b22a6710331ee68f2e4d1cd73a584d57c8821df79356",
				"0xcbc920e7bb89cbcb540a469a16226bf1057825283ab8eac3f45d00811eef8a64",
				"0x6d604ffc644a282fca8cb8e778e1e3f8245d8bd1d49326e3016a3c878ba0cbbd",
			},
			expectError: false,
		},
		{
			name:           "Invalid RLP encoding",
			input:          "invalid_hex",
			expectedHashes: nil,
			expectError:    true,
		},
		{
			name:           "Empty input",
			input:          "",
			expectedHashes: nil,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hashes, err := ethereum.DecodeHashes(tc.input)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, hashes)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedHashes, hashes)
		})
	}
}
