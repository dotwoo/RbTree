package rbtree

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"sort"
	"testing"

	"github.com/maxim2266/csvplus"
)

type key uint32

func (n key) LessThan(b interface{}) bool {
	value, _ := b.(key)
	return n < value
}

func Test_Preorder(t *testing.T) {
	tree := NewTree()

	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")
	if tree.Size() != 6 {
		t.Error("Error size")
		return
	}
	tree.Preorder()
}

func Test_Find(t *testing.T) {

	tree := NewTree()

	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")

	n := tree.FindIt(key(4))
	if n.Value != "dfa3" {
		t.Error("Error value")
		return
	}
	n.Value = "bdsf"
	if n.Value != "bdsf" {
		t.Error("Error value modify")
		return
	}
	value := tree.Find(key(5)).(string)
	if value != "jcd4" {
		t.Error("Error value after modifyed other node")
		return
	}
}
func Test_Iterator(t *testing.T) {
	tree := NewTree()

	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")

	it := tree.Iterator()

	for it != nil {
		it = it.Next()
	}

}

func Test_Delete(t *testing.T) {
	tree := NewTree()

	tree.Insert(key(1), "123")
	tree.Insert(key(3), "234")
	tree.Insert(key(4), "dfa3")
	tree.Insert(key(6), "sd4")
	tree.Insert(key(5), "jcd4")
	tree.Insert(key(2), "bcd4")
	for i := 1; i <= 6; i++ {
		tree.Delete(key(i))
		if tree.Size() != 6-i {
			t.Error("Delete Error")
		}
		if tree.FindIt(key(i)) != nil {
			t.Error("Delete Error")
		}
	}
	tree.Insert(key(1), "bcd4")
	tree.Clear()
	tree.Preorder()
	if tree.Find(key(1)) != nil {
		t.Error("Can't clear")
		return
	}
}

func TestDeleteLarge(t *testing.T) {
	int2Bool := ReadGeo()
	keys, all := KeysOfMap(int2Bool)
	fmt.Println("all item:", all)
	//build
	tree := NewTree()
	limit := 100

	for i, k := range keys {
		tree.Insert(key(k), int2Bool[k])
		if i > limit {
			break
		}
	}
	//check
	for i, k := range keys {
		v := int2Bool[k]
		r := tree.FindIt(key(k))
		if r == nil || r.Value.(bool) != v {
			t.Error("Cant't find:", k, v, r)
			return
		}
		if i > limit {
			break
		}
	}
	//return
	//delete
	for i, k := range keys {
		v := int2Bool[k]
		if !v {
			tree.Delete(key(k))
		}
		if i > limit {
			break
		}
	}

	//delete check
	for i, k := range keys {
		v := int2Bool[k]
		r := tree.FindIt(key(k))
		if v && (r == nil || r.Value.(bool) != v) {
			t.Error("Cant't find:", k, v, r)
			return
		} else if r != nil {
			t.Error("Cant't find:", k, v, r)
			return
		}
		if i > limit {
			break
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
