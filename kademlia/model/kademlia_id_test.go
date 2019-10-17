package model

import (
	"crypto/sha1"
	"gotest.tools/assert"
	"testing"
)

func TestKademliaID_NewID(t *testing.T) {
	s := []byte("1111111111133333333334000000000000000000")
	sha := sha1.Sum(s)

	k1 := NewKademliaID(s)
	k2, err := KademliaIDFromBytes(sha[:])
	assert.NilError(t, err)
	assert.Assert(t, k1.equals(k2))
}

func TestKademliaIDFromString_StringShouldReturnTheSame(t *testing.T) {
	s := "1111111111133333333334000000000000000000"
	assert.Equal(t, s, KademliaIDFromString(s).String())
}

func TestKademliaID_Equal(t *testing.T) {
	k := NewRandomKademliaID()
	assert.Assert(t, k.equals(k))
	assert.Assert(t, !k.equals(NewRandomKademliaID()))
}

func TestKademliaID_Less(t *testing.T) {
	max := KademliaIDFromString("ffffffffffffffffffffffffffffffffffffffff")
	mid := KademliaIDFromString("8888888888888888888888888888888888888888")
	min := KademliaIDFromString("0000000000000000000000000000000000000000")

	assert.Assert(t, min.less(mid))
	assert.Assert(t, mid.less(max))
	assert.Assert(t, min.less(max))

	assert.Assert(t, !max.less(min))
	assert.Assert(t, !max.less(mid))
	assert.Assert(t, !mid.less(min))
}

func TestKademliaID_Distance(t *testing.T) {
	var sortedIDs = []string{
		"ffffffff00000000000000000000000000000000",
		"1111111400000000000000000000000000000000",
		"1111111300000000000000000000000000000000",
		"1111111200000000000000000000000000000000",
		"1111111111133333333334000000000000000000",
		"1111111111133333333333000000000000000000",
		"1111111111111111222222000000000000000000",
		"1111111111111111111111000000000000000000",
		"1111111111111111111110000000000000000000",
		"1111111100000000000000000000000000000000",
	}
	target := KademliaIDFromString("ffffffffffffffffffffffffffffffffffffffff")

	for i := range sortedIDs[:len(sortedIDs)-1] {
		k1 := KademliaIDFromString(sortedIDs[i])
		d1 := k1.calcDistance(target)

		k2 := KademliaIDFromString(sortedIDs[i+1])
		d2 := k2.calcDistance(target)

		assert.Assert(t, d1.less(d2))
		assert.Assert(t, !d2.less(d1))
	}
}
