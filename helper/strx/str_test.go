package strx

import "testing"

func TestB2s(t *testing.T) {
	B2s([]byte("hello"))
}

func TestS2b(t *testing.T) {
	S2b("hello")
}

func TestStr2Number(t *testing.T) {
	in := "1"
	str2Number := Str2Number[int](in)
	t.Log(str2Number)
	str2Number1 := Str2Number[float64](in)
	t.Log(str2Number1)
	str2Number2 := Str2Number[uint64](in)
	t.Log(str2Number2)
	in = "0"
	str2Number = Str2Number[int](in)
	t.Log(str2Number)
	str2Number1 = Str2Number[float64](in)
	t.Log(str2Number1)
	str2Number2 = Str2Number[uint64](in)
	t.Log(str2Number2)
	in = ""
	str2Number = Str2Number[int](in)
	t.Log(str2Number)
	str2Number1 = Str2Number[float64](in)
	t.Log(str2Number1)
	str2Number2 = Str2Number[uint64](in)
	t.Log(str2Number2)
}

func TestStr2Numbers(t *testing.T) {
	in := []string{"", "10", "11", "", "25"}
	str2Number := Str2Numbers[int](in...)
	t.Log(str2Number)
	str2Number1 := Str2Numbers[float64](in...)
	t.Log(str2Number1)
	str2Number2 := Str2Numbers[uint64](in...)
	t.Log(str2Number2)
}

func TestStr2Map(t *testing.T) {
	var in = `{"cn": "123"}`
	str2Map := S2Map(in)
	t.Log(str2Map)
	str2Map = B2Map([]byte(in))
	t.Log(str2Map)
}
