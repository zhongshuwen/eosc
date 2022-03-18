package cli

import (
	"testing"

	zsw "github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/ecc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	pubKey1, _ := ecc.NewPublicKey("EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV")
	pubKey2, _ := ecc.NewPublicKey("EOS4xenzB8vAWjwHxnk8eGLkPumXDAEA1Sgq11U2muX3kJ8n7v2KA")

	tests := []struct {
		name   string
		input  string
		expect *zsw.Authority
	}{
		{
			name:  "full",
			input: "3=EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV+2,abourget@perm+3,bob",
			expect: &zsw.Authority{
				Threshold: uint32(3),
				Waits:     []zsw.WaitWeight{},
				Accounts: []zsw.PermissionLevelWeight{
					zsw.PermissionLevelWeight{
						Permission: zsw.PermissionLevel{
							Actor:      zsw.AccountName("abourget"),
							Permission: zsw.PermissionName("perm"),
						},
						Weight: 3,
					},
					zsw.PermissionLevelWeight{
						Permission: zsw.PermissionLevel{
							Actor:      zsw.AccountName("bob"),
							Permission: zsw.PermissionName("active"),
						},
						Weight: 1,
					},
				},
				Keys: []zsw.KeyWeight{
					zsw.KeyWeight{
						PublicKey: pubKey1,
						Weight:    2,
					},
				},
			},
		},
		{
			name:  "full of spaces",
			input: "3  =  EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV  +  2  ,  abourget@perm + 3 , bob",
			expect: &zsw.Authority{
				Threshold: uint32(3),
				Waits:     []zsw.WaitWeight{},
				Accounts: []zsw.PermissionLevelWeight{
					zsw.PermissionLevelWeight{
						Permission: zsw.PermissionLevel{
							Actor:      zsw.AccountName("abourget"),
							Permission: zsw.PermissionName("perm"),
						},
						Weight: 3,
					},
					zsw.PermissionLevelWeight{
						Permission: zsw.PermissionLevel{
							Actor:      zsw.AccountName("bob"),
							Permission: zsw.PermissionName("active"),
						},
						Weight: 1,
					},
				},
				Keys: []zsw.KeyWeight{
					zsw.KeyWeight{
						PublicKey: pubKey1,
						Weight:    2,
					},
				},
			},
		},
		{
			name:  "single account",
			input: "abourget",
			expect: &zsw.Authority{
				Threshold: uint32(1),
				Waits:     []zsw.WaitWeight{},
				Accounts: []zsw.PermissionLevelWeight{
					zsw.PermissionLevelWeight{
						Permission: zsw.PermissionLevel{
							Actor:      zsw.AccountName("abourget"),
							Permission: zsw.PermissionName("active"),
						},
						Weight: 1,
					},
				},
				Keys: []zsw.KeyWeight{},
			},
		},
		{
			name:  "single key",
			input: "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV",
			expect: &zsw.Authority{
				Threshold: uint32(1),
				Waits:     []zsw.WaitWeight{},
				Accounts:  []zsw.PermissionLevelWeight{},
				Keys: []zsw.KeyWeight{
					zsw.KeyWeight{
						PublicKey: pubKey1,
						Weight:    1,
					},
				},
			},
		},
		{
			name:  "sorted keys",
			input: "EOS6MRyAjQq8ud7hVNYcfnVPJqcVpscN5So8BhtHuGYqET5GDW5CV, EOS4xenzB8vAWjwHxnk8eGLkPumXDAEA1Sgq11U2muX3kJ8n7v2KA",
			expect: &zsw.Authority{
				Threshold: uint32(1),
				Waits:     []zsw.WaitWeight{},
				Accounts:  []zsw.PermissionLevelWeight{},
				Keys: []zsw.KeyWeight{
					zsw.KeyWeight{
						PublicKey: pubKey2,
						Weight:    1,
					},
					zsw.KeyWeight{
						PublicKey: pubKey1,
						Weight:    1,
					},
				},
			},
		},
		{
			name:  "sorted accounts",
			input: "alex, bob",
			expect: &zsw.Authority{
				Threshold: uint32(1),
				Waits:     []zsw.WaitWeight{},
				Accounts: []zsw.PermissionLevelWeight{
					zsw.PermissionLevelWeight{
						Permission: zsw.PermissionLevel{
							Actor:      zsw.AccountName("alex"),
							Permission: zsw.PermissionName("active"),
						},
						Weight: 1,
					},
					zsw.PermissionLevelWeight{
						Permission: zsw.PermissionLevel{
							Actor:      zsw.AccountName("bob"),
							Permission: zsw.PermissionName("active"),
						},
						Weight: 1,
					},
				},
				Keys: []zsw.KeyWeight{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := ParseShortFormAuth(test.input)
			require.NoError(t, err)
			assert.Equal(t, test.expect, res)
		})
	}
}
