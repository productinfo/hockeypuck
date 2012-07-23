/*
   Hockeypuck - OpenPGP key server
   Copyright (C) 2012  Casey Marshall

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, version 3.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package hockeypuck

import (
	"github.com/bmizerany/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestGetKeyword(t *testing.T) {
	// basic search
	testUrl, err := url.Parse("/pks/lookup?op=get&search=alice")
	assert.Equal(t, err, nil)
	req := &http.Request{
		Method: "GET",
		URL:    testUrl}
	lookup, err := parseLookup(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, Get, lookup.Op)
	assert.Equal(t, "alice", lookup.Search)
	assert.Equal(t, NoOption, lookup.Option)
	assert.Equal(t, false, lookup.Fingerprint)
	assert.Equal(t, false, lookup.Exact)
}

func TestGetFp(t *testing.T) {
	// fp search
	testUrl, err := url.Parse("/pks/lookup?op=get&search=0xdecafbad&options=mr,nm&fingerprint=on&exact=on")
	assert.Equal(t, err, nil)
	req := &http.Request{
		Method: "GET",
		URL:    testUrl}
	lookup, err := parseLookup(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, Get, lookup.Op)
	assert.Equal(t, "0xdecafbad", lookup.Search)
	assert.Equal(t, MachineReadable&lookup.Option, MachineReadable)
	assert.Equal(t, NotModifiable&lookup.Option, NotModifiable)
	assert.Equal(t, true, lookup.Fingerprint)
	assert.Equal(t, true, lookup.Exact)
}

func TestIndex(t *testing.T) {
	// index
	testUrl, err := url.Parse("/pks/lookup?op=index&search=sharin") // as in, foo
	assert.Equal(t, err, nil)
	req := &http.Request{
		Method: "GET",
		URL:    testUrl}
	lookup, err := parseLookup(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, Index, lookup.Op)
}

func TestVindex(t *testing.T) {
	// index
	testUrl, err := url.Parse("/pks/lookup?op=vindex&search=bob") // as in, foo
	assert.Equal(t, err, nil)
	req := &http.Request{
		Method: "GET",
		URL:    testUrl}
	lookup, err := parseLookup(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, Vindex, lookup.Op)
}

func TestMissingSearch(t *testing.T) {
	// index
	testUrl, err := url.Parse("/pks/lookup?op=get")
	assert.Equal(t, err, nil)
	req := &http.Request{
		Method: "GET",
		URL:    testUrl}
	_, err = parseLookup(req)
	// error without search term
	assert.NotEqual(t, err, nil)
}

func TestNoSuchOp(t *testing.T) {
	// index
	testUrl, err := url.Parse("/pks/lookup?op=barrelroll")
	assert.Equal(t, err, nil)
	req := &http.Request{
		Method: "GET",
		URL:    testUrl}
	_, err = parseLookup(req)
	// Unknown operation
	assert.NotEqual(t, err, nil)
}

func TestAdd(t *testing.T) {
	testUrl, err := url.Parse("/pks/add")
	assert.Equal(t, err, nil)
	postData := make(map[string][]string)
	postData["keytext"] = []string{"sus llaves aqui"}
	req := &http.Request{
		Method: "POST",
		URL:    testUrl,
		Form:   url.Values(postData)}
	add, err := parseAdd(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, "sus llaves aqui", add.Keytext)
	assert.Equal(t, NoOption, add.Option)
}

func TestAddOptions(t *testing.T) {
	testUrl, err := url.Parse("/pks/add?options=mr")
	assert.Equal(t, err, nil)
	postData := make(map[string][]string)
	postData["keytext"] = []string{"sus llaves estan aqui"}
	postData["options"] = []string{"mr"}
	req := &http.Request{
		Method: "POST",
		URL:    testUrl,
		Form:   url.Values(postData)}
	add, err := parseAdd(req)
	assert.Equal(t, err, nil)
	assert.Equal(t, "sus llaves estan aqui", add.Keytext)
	assert.Equal(t, MachineReadable&add.Option, MachineReadable)
	assert.Equal(t, NotModifiable&add.Option, NoOption)
}

func TestAddMissingKey(t *testing.T) {
	testUrl, err := url.Parse("/pks/add")
	assert.Equal(t, err, nil)
	postData := make(map[string][]string)
	req := &http.Request{
		Method: "POST",
		URL:    testUrl,
		Form:   url.Values(postData)}
	_, err = parseAdd(req)
	// error without keytext
	assert.NotEqual(t, err, nil)
}
