package app_test

import (
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/app"
)

var (
	//go:embed testdata/app_minimal.go
	AppMinimalFile []byte
	//go:embed testdata/app_generic.go
	AppGenericFile []byte
	//go:embed testdata/no_app.go
	NoAppFile []byte
	//go:embed testdata/two_app.go
	TwoAppFile []byte
	//go:embed testdata/app_full.go
	AppFullFile []byte
)

func TestCheckKeeper(t *testing.T) {
	tests := []struct {
		name          string
		appFile       []byte
		keeperName    string
		expectedError string
	}{
		{
			name:       "minimal app",
			appFile:    AppMinimalFile,
			keeperName: "FooKeeper",
		},
		{
			name:       "generic app",
			appFile:    AppGenericFile,
			keeperName: "FooKeeper",
		},
		{
			name:          "no app",
			appFile:       NoAppFile,
			keeperName:    "FooKeeper",
			expectedError: "app.go should contain a single app (got 0)",
		},
		{
			name:          "two apps",
			appFile:       TwoAppFile,
			keeperName:    "FooKeeper",
			expectedError: "app.go should contain a single app (got 2)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "app.go")
			err := os.WriteFile(tmpFile, tt.appFile, 0o644)
			require.NoError(t, err)

			err = app.CheckKeeper(tmpDir, tt.keeperName)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetRegisteredModules(t *testing.T) {
	tmpDir := t.TempDir()

	tmpFile := filepath.Join(tmpDir, "app.go")
	err := os.WriteFile(tmpFile, AppFullFile, 0o644)
	require.NoError(t, err)

	tmpNoAppFile := filepath.Join(tmpDir, "someOtherFile.go")
	err = os.WriteFile(tmpNoAppFile, NoAppFile, 0o644)
	require.NoError(t, err)

	registeredModules, err := app.FindRegisteredModules(tmpDir)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{
		"github.com/cosmos/cosmos-sdk/x/auth",
		"github.com/cosmos/cosmos-sdk/x/genutil",
		"github.com/cosmos/cosmos-sdk/x/bank",
		"github.com/cosmos/cosmos-sdk/x/capability",
		"github.com/cosmos/cosmos-sdk/x/staking",
		"github.com/cosmos/cosmos-sdk/x/mint",
		"github.com/cosmos/cosmos-sdk/x/distribution",
		"github.com/cosmos/cosmos-sdk/x/gov",
		"github.com/cosmos/cosmos-sdk/x/params",
		"github.com/cosmos/cosmos-sdk/x/crisis",
		"github.com/cosmos/cosmos-sdk/x/slashing",
		"github.com/cosmos/cosmos-sdk/x/feegrant/module",
		"github.com/cosmos/ibc-go/v5/modules/core",
		"github.com/cosmos/cosmos-sdk/x/upgrade",
		"github.com/cosmos/cosmos-sdk/x/evidence",
		"github.com/cosmos/ibc-go/v5/modules/apps/transfer",
		"github.com/cosmos/cosmos-sdk/x/auth/vesting",
		"github.com/tendermint/testchain/x/testchain",
		"github.com/tendermint/testchain/x/queryonlymod",
		"github.com/cosmos/cosmos-sdk/x/auth/tx",
		"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
	}, registeredModules)
}

func TestParseAppModules(t *testing.T) {
	basicModules := []string{
		"github.com/cosmos/cosmos-sdk/x/auth",
		"github.com/cosmos/cosmos-sdk/x/bank",
		"github.com/cosmos/cosmos-sdk/x/staking",
		"github.com/username/test/x/foo",
	}

	cases := []struct {
		name            string
		path            string
		expectedModules []string
	}{
		{
			name:            "new basic manager arguments",
			path:            "testdata/modules/arguments",
			expectedModules: basicModules,
		},
		{
			name:            "same file variable",
			path:            "testdata/modules/file_variable",
			expectedModules: basicModules,
		},
		{
			name:            "same package variable",
			path:            "testdata/modules/package_variable",
			expectedModules: basicModules,
		},
		{
			name:            "other package variable",
			path:            "testdata/modules/external_variable",
			expectedModules: basicModules,
		},
		{
			name: "with api routes",
			path: "testdata/modules/api_routes",
			expectedModules: append(
				basicModules,
				"github.com/cosmos/cosmos-sdk/x/auth/tx",
				"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
			),
		},
		{
			name: "gaia",
			path: "testdata/modules/gaia",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/genutil",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/cosmos/cosmos-sdk/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/feegrant/module",
				"github.com/cosmos/cosmos-sdk/x/authz/module",
				"github.com/cosmos/cosmos-sdk/x/group/module",
				"github.com/cosmos/ibc-go/v5/modules/core",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v5/modules/apps/transfer",
				"github.com/cosmos/cosmos-sdk/x/auth/vesting",
				"github.com/gravity-devs/liquidity/v2/x/liquidity",
				"github.com/strangelove-ventures/packet-forward-middleware/v2/router",
				"github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts",
				"github.com/cosmos/gaia/v8/x/icamauth",
				"github.com/cosmos/gaia/v8/x/globalfee",
				"github.com/cosmos/cosmos-sdk/x/auth/tx",
				"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
			},
		},
		{
			name: "spn",
			path: "testdata/modules/spn",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/genutil",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/ignite/modules/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/feegrant/module",
				"github.com/cosmos/cosmos-sdk/x/authz/module",
				"github.com/cosmos/ibc-go/v5/modules/core",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v5/modules/apps/transfer",
				"github.com/cosmos/cosmos-sdk/x/auth/vesting",
				"github.com/tendermint/spn/x/participation",
				"github.com/ignite/modules/x/claim",
				"github.com/tendermint/spn/x/profile",
				"github.com/tendermint/spn/x/launch",
				"github.com/tendermint/spn/x/campaign",
				"github.com/tendermint/spn/x/monitoringc",
				"github.com/tendermint/spn/x/monitoringp",
				"github.com/tendermint/spn/x/reward",
				"github.com/tendermint/fundraising/x/fundraising",
				"github.com/cosmos/cosmos-sdk/x/auth/tx",
				"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
			},
		},
		{
			name: "juno",
			path: "testdata/modules/juno",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/genutil",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/CosmosContracts/juno/v10/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/ibc-go/v3/modules/core",
				"github.com/cosmos/cosmos-sdk/x/feegrant/module",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v3/modules/apps/transfer",
				"github.com/cosmos/cosmos-sdk/x/auth/vesting",
				"github.com/cosmos/cosmos-sdk/x/authz/module",
				"github.com/CosmWasm/wasmd/x/wasm",
				"github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts",
				"github.com/cosmos/cosmos-sdk/x/auth/tx",
				"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
			},
		},
		{
			name: "osmosis",
			path: "testdata/modules/osmosis",
			expectedModules: []string{
				"github.com/cosmos/cosmos-sdk/x/auth",
				"github.com/cosmos/cosmos-sdk/x/genutil",
				"github.com/cosmos/cosmos-sdk/x/bank",
				"github.com/cosmos/cosmos-sdk/x/capability",
				"github.com/cosmos/cosmos-sdk/x/staking",
				"github.com/osmosis-labs/osmosis/v12/x/mint",
				"github.com/cosmos/cosmos-sdk/x/distribution",
				"github.com/cosmos/cosmos-sdk/x/gov",
				"github.com/cosmos/cosmos-sdk/x/params",
				"github.com/cosmos/cosmos-sdk/x/crisis",
				"github.com/cosmos/cosmos-sdk/x/slashing",
				"github.com/cosmos/cosmos-sdk/x/authz/module",
				"github.com/cosmos/ibc-go/v3/modules/core",
				"github.com/cosmos/cosmos-sdk/x/upgrade",
				"github.com/cosmos/cosmos-sdk/x/evidence",
				"github.com/cosmos/ibc-go/v3/modules/apps/transfer",
				"github.com/cosmos/cosmos-sdk/x/auth/vesting",
				"github.com/osmosis-labs/osmosis/v12/x/gamm",
				"github.com/osmosis-labs/osmosis/v12/x/twap/twapmodule",
				"github.com/osmosis-labs/osmosis/v12/x/txfees",
				"github.com/osmosis-labs/osmosis/v12/x/incentives",
				"github.com/osmosis-labs/osmosis/v12/x/lockup",
				"github.com/osmosis-labs/osmosis/v12/x/pool-incentives",
				"github.com/osmosis-labs/osmosis/v12/x/epochs",
				"github.com/osmosis-labs/osmosis/v12/x/superfluid",
				"github.com/osmosis-labs/osmosis/v12/x/tokenfactory",
				"github.com/CosmWasm/wasmd/x/wasm",
				"github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts",
				"github.com/cosmos/cosmos-sdk/x/auth/tx",
				"github.com/cosmos/cosmos-sdk/client/grpc/tmservice",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			m, err := app.FindRegisteredModules(tt.path)

			require.NoError(t, err)
			require.Equal(t, tt.expectedModules, m)
		})
	}
}
