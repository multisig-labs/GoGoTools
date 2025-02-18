package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerifyBLSKeys(t *testing.T) {
	pubkey := "0x8d543b279b9bd69c5b6754a09bce1eab2de2d9135eff7e391e42583fca4c19c6007e864971c2baba777dfa312ca7994e"
	pop := "0xa99743e050b543f2482c1010e2908b848c5894080a5e5ac9db96111b721d753efbe106eb12496b56be14f11feaeb4d9605507ca0cb726b3832c45448603025de7ba425d7c054f5e664922d6a5dfab8b9f7df681941d25be420ea4f973b79d041"

	err := ValidateBLSKeys(pubkey, pop)
	require.NoError(t, err)
}
