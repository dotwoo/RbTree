// Copyright 2015, Hu Keping. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbtree

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"reflect"
	"sort"
	"testing"

	"github.com/maxim2266/csvplus"
)

func TestInsertAndDelete(t *testing.T) {
	rbt := New()

	m := 0
	n := 1000
	for m < n {
		rbt.Insert(Int(m))
		m++
	}
	if rbt.Len() != uint(n) {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), n)
	}

	for m > 0 {
		rbt.Delete(Int(m))
		m--
	}
	if rbt.Len() != 1 {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), 1)
	}
}

type testStruct struct {
	id   int
	text string
}

func (ts *testStruct) Less(than Item) bool {
	return ts.id < than.(*testStruct).id
}

func TestInsertOrGet(t *testing.T) {
	rbt := New()

	items := []*testStruct{
		{1, "this"},
		{2, "is"},
		{3, "a"},
		{4, "test"},
	}

	for i := range items {
		rbt.Insert(items[i])
	}

	newItem := &testStruct{items[0].id, "not"}
	newItem = rbt.InsertOrGet(newItem).(*testStruct)

	if newItem.text != items[0].text {
		t.Errorf("tree.InsertOrGet = {id: %d, text: %s}, expect {id %d, text %s}", newItem.id, newItem.text, items[0].id, items[0].text)
	}

	newItem = &testStruct{5, "new"}
	newItem = rbt.InsertOrGet(newItem).(*testStruct)

	if newItem.text != "new" {
		t.Errorf("tree.InsertOrGet = {id: %d, text: %s}, expect {id %d, text %s}", newItem.id, newItem.text, 5, "new")
	}
}

func TestInsertString(t *testing.T) {
	rbt := New()

	rbt.Insert(String("go"))
	rbt.Insert(String("lang"))

	if rbt.Len() != 2 {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), 2)
	}
}

// Test for duplicate
func TestInsertDup(t *testing.T) {
	rbt := New()

	rbt.Insert(String("go"))
	rbt.Insert(String("go"))
	rbt.Insert(String("go"))

	if rbt.Len() != 1 {
		t.Errorf("tree.Len() = %d, expect %d", rbt.Len(), 1)
	}
}

func TestDescend(t *testing.T) {
	rbt := New()

	m := 0
	n := 10
	for m < n {
		rbt.Insert(Int(m))
		m++
	}

	var ret []Item

	rbt.Descend(Int(1), func(i Item) bool {
		ret = append(ret, i)
		return true
	})
	expected := []Item{Int(1), Int(0)}
	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("expected %v but got %v", expected, ret)
	}

	ret = nil
	rbt.Descend(Int(10), func(i Item) bool {
		ret = append(ret, i)
		return true
	})
	expected = []Item{Int(9), Int(8), Int(7), Int(6), Int(5), Int(4), Int(3), Int(2), Int(1), Int(0)}
	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("expected %v but got %v", expected, ret)
	}
}

func TestGet(t *testing.T) {
	rbt := New()

	rbt.Insert(Int(1))
	rbt.Insert(Int(2))
	rbt.Insert(Int(3))

	no := rbt.Get(Int(100))
	ok := rbt.Get(Int(1))

	if no != nil {
		t.Errorf("100 is expect not exists")
	}

	if ok == nil {
		t.Errorf("1 is expect exists")
	}
}

func TestAscend(t *testing.T) {
	rbt := New()

	rbt.Insert(String("a"))
	rbt.Insert(String("b"))
	rbt.Insert(String("c"))
	rbt.Insert(String("d"))

	rbt.Delete(rbt.Min())

	var ret []Item
	rbt.Ascend(rbt.Min(), func(i Item) bool {
		ret = append(ret, i)
		return true
	})

	expected := []Item{String("b"), String("c"), String("d")}
	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("expected %v but got %v", expected, ret)
	}
}

func TestMax(t *testing.T) {
	rbt := New()

	rbt.Insert(String("z"))
	rbt.Insert(String("h"))
	rbt.Insert(String("a"))

	expected := String("z")
	if rbt.Max() != expected {
		t.Errorf("expected Max of tree as %v but got %v", expected, rbt.Max())
	}
}

func TestAscendRange(t *testing.T) {
	rbt := New()

	strings := []String{"a", "b", "c", "aa", "ab", "ac", "abc", "acb", "bac"}
	for _, v := range strings {
		rbt.Insert(v)
	}

	var ret []Item
	rbt.AscendRange(String("ab"), String("b"), func(i Item) bool {
		ret = append(ret, i)
		return true
	})
	expected := []Item{String("ab"), String("abc"), String("ac"), String("acb")}

	if !reflect.DeepEqual(ret, expected) {
		t.Errorf("expected %v but got %v", expected, ret)
	}
}

type key struct {
	Key   uint32
	Value bool
}

func (x key) Less(than Item) bool {
	value, _ := than.(key)
	return x.Key < value.Key

}

func TestInsertLarge(t *testing.T) {
	int2Bool := ReadGeo()
	limit := 269649
	sub := SubOfMap(limit, int2Bool)
	fmt.Println("all item:", len(int2Bool))
	//build
	tree := New()

	for k, v := range sub {
		tree.Insert(key{k, v})
	}
	//check
	for k, v := range sub {
		r := tree.Get(key{k, v})
		if r == nil || r.(key).Value != v {
			t.Error("Cant't find:", k, v, r)
			return
		}
	}
}
func TestIndexLarge(t *testing.T) {
	int2Bool := ReadGeo()
	index, _ := KeysOfMap(int2Bool)
	//build
	tree := New()

	for k, _ := range int2Bool {
		tree.Insert(Uint32(k))
	}
	first := tree.First()
	for _, k := range index {
		if k != uint32((first.Item).(Uint32)) {
			t.Error("Cant't find:", k, first.Item)
			return
		}
		//fmt.Println(k, first.Item)
		first = tree.Next(first)
	}
}

func TestDeleteLarge(t *testing.T) {
	int2Bool := ReadGeo()
	limit := 2002000
	sub := SubOfMap(limit, int2Bool)
	fmt.Println("all item:", len(int2Bool))
	//build
	tree := New()

	for k, v := range sub {
		tree.Insert(key{k, v})
	}
	//check
	for k, v := range sub {
		r := tree.Get(key{k, v})
		if r == nil || r.(key).Value != v {
			t.Error("Cant't find1:", k, v, r)
			return
		}
	}
	//return
	//delete
	for k, v := range sub {
		if !v {
			tree.Delete(key{k, v})
		}
	}

	//delete check
	for k, v := range sub {
		r := tree.Get(key{k, v})
		if v && (r == nil || r.(key).Value != v) {
			t.Error("Cant't find2:", k, v, r)
			return
		} else if !v && r != nil {
			t.Error("Cant't find3:", k, v, r)
			return
		}
	}
}

func ReadGeo() map[uint32]bool {
	gb := "GeoLite2-Country-Blocks-IPv4.csv"
	BlockCSV := csvplus.FromFile(gb).SelectColumns("network", "geoname_id")
	ip2BoolMap := make(map[uint32]bool, 0)
	pre := ""
	preInt := uint32(0)

	BlockCSV.Iterate(func(r csvplus.Row) error {
		gid := r.SafeGetValue("geoname_id", "-1")
		nw := r.SafeGetValue("network", "")

		end := uint32(0)
		start, n, err := net.ParseCIDR(nw)
		if err != nil {
			log.Println("geo network parse error:", nw, err.Error())
			return nil
		} else {
			end = Ip2Uint(Endip(n)) + 1
		}
		startInt := Ip2Uint(start)
		isTrue := (pre == gid)
		if preInt < startInt {
			ip2BoolMap[preInt] = isTrue
		}
		ip2BoolMap[startInt] = isTrue
		preInt = end
		return nil
	})

	return ip2BoolMap
}

//获取网段最后ip
func Endip(n *net.IPNet) net.IP {
	nip := n.IP.To4()
	if nip == nil {
		mast := n.Mask
		for i := 0; i < 16; i++ {
			nip[i] = nip[i] | (^mast[i])
		}
		return nip
	}
	mast := n.Mask
	for i := 0; i < 4; i++ {
		nip[i] = nip[i] | (^mast[i])
	}
	return nip
}

func Ip2Uint(i net.IP) uint32 {
	nip := i.To4()
	if nip == nil {
		return 0
	}
	return binary.BigEndian.Uint32(nip)
}

type Uint32Slice []uint32

func (c Uint32Slice) Len() int {
	return len(c)
}
func (c Uint32Slice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Uint32Slice) Less(i, j int) bool {
	return c[i] < c[j]
}

func KeysOfMap(m map[uint32]bool) (keys Uint32Slice, cnt int) {
	keys = make(Uint32Slice, len(m))
	cnt = 0
	for key := range m {
		keys[cnt] = key
		cnt++
	}

	sort.Sort(keys)
	return keys, cnt
}

func SubOfMap(limit int, m map[uint32]bool) map[uint32]bool {
	i := 0
	out := make(map[uint32]bool, 0)
	for k, v := range m {
		out[k] = v
		i++
		if i >= limit {
			break
		}
	}
	return out
}
