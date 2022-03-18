package cli

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	zsw "github.com/zhongshuwen/zswchain-go"
	"github.com/zhongshuwen/zswchain-go/ecc"
)

var reValidAccount = regexp.MustCompile(`[a-z12345]*`)

// ToAccountName converts a eos valid name string (in) into an eos-go
// AccountName struct
func ToAccountName(in string) (out zsw.AccountName, err error) {
	if !reValidAccount.MatchString(in) {
		err = fmt.Errorf("invalid characters in %q, allowed: 'a' through 'z', and '1', '2', '3', '4', '5'", in)
		return
	}

	val, _ := zsw.StringToName(in)
	if zsw.NameToString(val) != in {
		err = fmt.Errorf("invalid name, 13 characters maximum")
		return
	}

	if len(in) == 0 {
		err = fmt.Errorf("empty")
		return
	}

	return zsw.AccountName(in), nil
}

// ToAsset converts a eos valid asset string (in) into an eos-go
// Asset struct
func ToAsset(in string) (out zsw.Asset, err error) {
	return zsw.NewAsset(in)
}

// ToName converts a valid eos name string (in) into an eos-go
// Name struct
func ToName(in string) (out zsw.Name, err error) {
	name, err := ToAccountName(in)
	if err != nil {
		return
	}
	return zsw.Name(name), nil
}

var shortFormTopLevelRE = regexp.MustCompile(`((\d{1,3})\s*=\s*)?(.*)`)
var shortFormKeyOrAccountRE = regexp.MustCompile(`\s*(([A-Za-z0-9]{48,64})|(([a-z1-5\.]{1,13})(@([a-z1-5\.]{1,13}))?))(\s*\+\s*(\d{1,3}))?\s*`)

func ParseShortFormAuth(in string) (*zsw.Authority, error) {

	match := shortFormTopLevelRE.FindStringSubmatch(in)
	if match == nil {
		return nil, fmt.Errorf(`invalid expression %q, example: "3=EOSKey1...,EOSKey2+2,account@perm+2"`, in)
	}

	threshold := uint32(1)
	if t, err := strconv.Atoi(match[2]); err == nil {
		threshold = uint32(t)
	}

	if threshold == 0 {
		return nil, fmt.Errorf("threshold cannot be 0")
	}

	auth := &zsw.Authority{
		Threshold: uint32(threshold),
		Waits:     []zsw.WaitWeight{},
		Accounts:  []zsw.PermissionLevelWeight{},
		Keys:      []zsw.KeyWeight{},
	}

	rest := match[3]

	for _, part := range strings.Split(rest, ",") {
		match = shortFormKeyOrAccountRE.FindStringSubmatch(part)
		if match == nil {
			return nil, fmt.Errorf(`invalid expression %q, example: "EOSKey1...+2" or "account@perm+2"`, part)
		}

		//fmt.Printf("match %q\n", match)

		key := match[2]

		weight := match[8]
		newWeight := uint16(1)
		if weight != "" {
			w, _ := strconv.Atoi(weight)
			newWeight = uint16(w)
		}

		if key != "" {
			pubKey, err := ecc.NewPublicKey(key)
			if err != nil {
				return nil, fmt.Errorf("invalid key %q: %w", key, err)
			}

			auth.Keys = append(auth.Keys, zsw.KeyWeight{
				PublicKey: pubKey,
				Weight:    newWeight,
			})
		} else {
			account := match[4]
			permission := match[6]
			if permission == "" {
				permission = "active"
			}

			if !validateName(account) {
				return nil, fmt.Errorf("invalid account name encoding %q", account)
			}

			if !validateName(permission) {
				return nil, fmt.Errorf("invalid permission name encoding %q", permission)
			}

			auth.Accounts = append(auth.Accounts, zsw.PermissionLevelWeight{
				Permission: zsw.PermissionLevel{
					Actor:      zsw.AccountName(account),
					Permission: zsw.PermissionName(permission),
				},
				Weight: newWeight,
			})
		}
	}

	sort.Slice(auth.Keys, func(i, j int) bool {
		return bytes.Compare(auth.Keys[i].PublicKey.Content, auth.Keys[j].PublicKey.Content) == -1
	})
	sort.Slice(auth.Accounts, func(i, j int) bool {
		return isPermissionLess(auth.Accounts[i].Permission, auth.Accounts[j].Permission)
	})

	return auth, nil
}

func validateName(in string) bool {
	val, err := zsw.StringToName(in)
	if err != nil {
		return false
	}
	check := zsw.NameToString(val)
	return check == in
}

func isPermissionLess(left, right zsw.PermissionLevel) bool {
	actor1 := zsw.MustStringToName(string(left.Actor))
	actor2 := zsw.MustStringToName(string(right.Actor))
	perm1 := zsw.MustStringToName(string(left.Permission))
	perm2 := zsw.MustStringToName(string(right.Permission))

	if actor1 < actor2 {
		return true
	}
	if actor1 > actor2 {
		return false
	}
	if perm1 < perm2 {
		return true
	}
	return false
}
