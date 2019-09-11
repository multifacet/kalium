// Copyright 2019 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package header_test

import (
	"bytes"
	"testing"
	"time"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/header"
)

// TestNDPNeighborSolicit tests the functions of header.NDPNeighborSolicit.
func TestNDPNeighborSolicit(t *testing.T) {
	b := []byte{
		0, 0, 0, 0,
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	// Test getting the Target Address.
	ns := header.NDPNeighborSolicit(b)
	addr := tcpip.Address("\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10")
	if got := ns.TargetAddress(); got != addr {
		t.Fatalf("got ns.TargetAddress = %s, want %s", got, addr)
	}

	// Test updating the Target Address.
	addr2 := tcpip.Address("\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x11")
	ns.SetTargetAddress(addr2)
	if got := ns.TargetAddress(); got != addr2 {
		t.Fatalf("got ns.TargetAddress = %s, want %s", got, addr2)
	}
	// Make sure the address got updated in the backing buffer.
	if got := tcpip.Address(b[header.NDPNSTargetAddessOffset:][:header.IPv6AddressSize]); got != addr2 {
		t.Fatalf("got targetaddress buffer = %s, want %s", got, addr2)
	}
}

// TestNDPNeighborAdvert tests the functions of header.NDPNeighborAdvert.
func TestNDPNeighborAdvert(t *testing.T) {
	b := []byte{
		160, 0, 0, 0,
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	}

	// Test getting the Target Address.
	na := header.NDPNeighborAdvert(b)
	addr := tcpip.Address("\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10")
	if got := na.TargetAddress(); got != addr {
		t.Fatalf("got TargetAddress = %s, want %s", got, addr)
	}

	// Test getting the Router Flag.
	if got := na.RouterFlag(); !got {
		t.Fatalf("got RouterFlag = %d, want = %d", true)
	}

	// Test getting the Solicited Flag.
	if got := na.SolicitedFlag(); got {
		t.Fatalf("got SolicitedFlag = %d, want = %d", false)
	}

	// Test getting the Override Flag.
	if got := na.OverrideFlag(); !got {
		t.Fatalf("got OverrideFlag = %d, want = %d", true)
	}

	// Test updating the Target Address.
	addr2 := tcpip.Address("\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x11")
	na.SetTargetAddress(addr2)
	if got := na.TargetAddress(); got != addr2 {
		t.Fatalf("got TargetAddress = %s, want %s", got, addr2)
	}
	// Make sure the address got updated in the backing buffer.
	if got := tcpip.Address(b[header.NDPNATargetAddessOffset:][:header.IPv6AddressSize]); got != addr2 {
		t.Fatalf("got targetaddress buffer = %s, want %s", got, addr2)
	}

	// Test updating the Router Flag.
	na.SetRouterFlag(false)
	if got := na.RouterFlag(); got {
		t.Fatalf("got RouterFlag = %d, want = %d", false)
	}

	// Test updating the Solicited Flag.
	na.SetSolicitedFlag(true)
	if got := na.SolicitedFlag(); !got {
		t.Fatalf("got SolicitedFlag = %d, want = %d", true)
	}

	// Test updating the Override Flag.
	na.SetOverrideFlag(false)
	if got := na.OverrideFlag(); got {
		t.Fatalf("got OverrideFlag = %d, want = %d", false)
	}

	// Make sure flags got updated in the backing buffer.
	if got := b[header.NDPNAFlagsOffset]; got != 64 {
		t.Fatalf("got flags byte = %d, want = 64")
	}
}

// TestNDPRouterAdvert tests the functions of header.NDPRouterAdvert.
func TestNDPRouterAdvert(t *testing.T) {
	b := []byte{
		64, 128, 1, 2,
		3, 4, 5, 6,
		7, 8, 9, 10,
	}

	ra := header.NDPRouterAdvert(b)

	// Test getting the Curr Hop Limit.
	if got := ra.CurrHopLimit(); got != 64 {
		t.Fatalf("got ra.CurrHopLimit = %d, want = 64", got)
	}

	// Test getting the ManagedAddrConfFlag.
	if got := ra.ManagedAddrConfFlag(); !got {
		t.Fatalf("got ManagedAddrConfFlag = %d, want = 1", got)
	}

	// Test getting the OtherConfFlag.
	if got := ra.OtherConfFlag(); got {
		t.Fatalf("got OtherConfFlag = %d, want = 0", got)
	}

	// Test getting the RouterLifetime.
	if got, want := ra.RouterLifetime(), time.Second*258; got != want {
		t.Fatalf("got ra.RouterLifetime = %d, want = %d", got, want)
	}

	// Test getting the ReachableTime.
	if got, want := ra.ReachableTime(), time.Millisecond*50595078; got != want {
		t.Fatalf("got ra.ReachableTime = %d, want = %d", got, want)
	}

	// Test getting the RetransTimer.
	if got, want := ra.RetransTimer(), time.Millisecond*117967114; got != want {
		t.Fatalf("got ra.RetransTimer = %d, want = %d", got, want)
	}
}

// TestNDPTargetLinkLayerAddressOptionSerialize tests serializing a
// header.NDPTargetLinkLayerAddressOption.
func TestNDPTargetLinkLayerAddressOptionSerialize(t *testing.T) {
	tests := []struct {
		name        string
		buf         []byte
		expectedBuf []byte
		addr        tcpip.LinkAddress
	}{
		{
			"Ethernet",
			[]byte{0, 0, 0, 0, 0, 0, 0, 0},
			[]byte{2, 1, 1, 2, 3, 4, 5, 6},
			tcpip.LinkAddress("\x01\x02\x03\x04\x05\x06"),
		},
		{
			"Padding",
			[]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			[]byte{2, 2, 1, 2, 3, 4, 5, 6, 7, 8, 0, 0, 0, 0, 0, 0},
			tcpip.LinkAddress("\x01\x02\x03\x04\x05\x06\x07\x08"),
		},
		{
			"Empty",
			[]byte{},
			[]byte{},
			tcpip.LinkAddress(""),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			opts := header.NDPOptions(test.buf)
			serializer := header.NDPOptionsSerializer{
				header.NDPTargetLinkLayerAddressOption(test.addr),
			}
			if got, want := int(serializer.Length()), len(test.expectedBuf); got != want {
				t.Fatalf("got Length = %d, want = %d", got, want)
			}
			opts.Serialize(serializer)
			if !bytes.Equal(test.buf, test.expectedBuf) {
				t.Fatalf("got b = %s, want = %d", test.buf, test.expectedBuf)
			}
		})
	}
}
