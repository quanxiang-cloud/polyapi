package arrange

import (
	"bytes"
	"fmt"
	"testing"
)

const testGenCases = false // generate test case data

func testNewNode(name string) *Node {
	return &Node{Name: name}
}
func assert(ok bool) {
	if !ok {
		panic("assert fail")
	}
}
func testVerify(q *nodeDeque, i, head, tail, size, _cap int, all string, _panic bool) bool {
	ok := q.head == head && q.tail == tail && q.Size() == size && q.Cap() == _cap && testGetAll(q) == all
	if !ok {
		fmt.Printf(`case:%d---- %d, %d, %d, %d, %s%s%s,%s`,
			i+1, q.head, q.tail, q.Size(), len(q.d), "`", testGetAll(q), "`", "\n")
		if false {
			for i, v := range q.d {
				fmt.Println(i, v)
			}
		}
	}
	if _panic {
		assert(ok)
		if q.Empty() {
			_, ok1 := q.Front()
			_, ok2 := q.Back()
			assert(!ok1 && !ok2)
		} else {
			n1, ok1 := q.Front()
			n2, ok2 := q.Back()
			assert(ok1 && ok2 && n1.Name == all[:1] && n2.Name == all[len(all)-1:])
		}
	}
	return ok
}

func testGetAll(q *nodeDeque) string {
	var b bytes.Buffer
	for i := q.head; i != q.tail; i = q.next(i) {
		b.WriteString(q.d[i].Name)
	}
	ret := b.String()
	return ret
}

func TestDeque(t *testing.T) {
	type testCase struct {
		ID         int
		Op         string
		Dir        string
		Head, Tail int
		Size, Cap  int
		All        string
	}
	tab := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890~!@#$%^&*()_+|{}[]?<>,.:;'"/\`
	tbIdx := 0
	nextCh := func() string {
		idx := tbIdx % len(tab)
		ret := tab[idx : idx+1]
		tbIdx++
		return ret
	}
	verify := func(q *nodeDeque, i int, c *testCase) {
		if testGenCases {
			fmt.Printf(`&testCase{%2d,"%s", "%s", %d, %d, %d, %d, %s%s%s},%s`,
				i+1, c.Op, c.Dir, q.head, q.tail, q.Size(), q.Cap(), "`", testGetAll(q), "`", "\n")
			return
		}
		testVerify(q, i, c.Head, c.Tail, c.Size, c.Cap, c.All, true)
	}
	cases := []*testCase{
		&testCase{1, "+", ">", 0, 1, 1, 5, `a`},
		&testCase{2, "+", "<", 4, 1, 2, 5, `ba`},
		&testCase{3, "+", ">", 4, 2, 3, 5, `bac`},
		&testCase{4, "+", "<", 3, 2, 4, 5, `dbac`},
		&testCase{5, "----", "<", 4, 2, 3, 5, `bac`},
		&testCase{6, "----", ">", 4, 1, 2, 5, `ba`},
		&testCase{7, "+", ">", 4, 2, 3, 5, `bae`},
		&testCase{8, "+", "<", 3, 2, 4, 5, `fbae`},
		&testCase{9, "+", ">", 0, 5, 5, 16, `fbaeg`},
		&testCase{10, "----", "<", 1, 5, 4, 16, `baeg`},
		&testCase{11, "----", ">", 1, 4, 3, 16, `bae`},
		&testCase{12, "+", "<", 0, 4, 4, 16, `hbae`},
		&testCase{13, "+", ">", 0, 5, 5, 16, `hbaei`},
		&testCase{14, "+", "<", 15, 5, 6, 16, `jhbaei`},
		&testCase{15, "+", ">", 15, 6, 7, 16, `jhbaeik`},
		&testCase{16, "----", "<", 0, 6, 6, 16, `hbaeik`},
		&testCase{17, "----", ">", 0, 5, 5, 16, `hbaei`},
		&testCase{18, "+", "<", 15, 5, 6, 16, `lhbaei`},
		&testCase{19, "+", ">", 15, 6, 7, 16, `lhbaeim`},
		&testCase{20, "+", "<", 14, 6, 8, 16, `nlhbaeim`},
		&testCase{21, "+", ">", 14, 7, 9, 16, `nlhbaeimo`},
		&testCase{22, "+", "<", 13, 7, 10, 16, `pnlhbaeimo`},
		&testCase{23, "----", "<", 14, 7, 9, 16, `nlhbaeimo`},
		&testCase{24, "----", ">", 14, 6, 8, 16, `nlhbaeim`},
		&testCase{25, "+", ">", 14, 7, 9, 16, `nlhbaeimq`},
		&testCase{26, "+", "<", 13, 7, 10, 16, `rnlhbaeimq`},
		&testCase{27, "+", ">", 13, 8, 11, 16, `rnlhbaeimqs`},
		&testCase{28, "+", "<", 12, 8, 12, 16, `trnlhbaeimqs`},
		&testCase{29, "+", ">", 12, 9, 13, 16, `trnlhbaeimqsu`},
		&testCase{30, "+", "<", 11, 9, 14, 16, `vtrnlhbaeimqsu`},
		&testCase{31, "----", "<", 12, 9, 13, 16, `trnlhbaeimqsu`},
		&testCase{32, "----", ">", 12, 8, 12, 16, `trnlhbaeimqs`},
		&testCase{33, "+", ">", 12, 9, 13, 16, `trnlhbaeimqsw`},
		&testCase{34, "+", "<", 11, 9, 14, 16, `xtrnlhbaeimqsw`},
		&testCase{35, "+", ">", 11, 10, 15, 16, `xtrnlhbaeimqswy`},
		&testCase{36, "+", "<", 0, 16, 16, 32, `zxtrnlhbaeimqswy`},
		&testCase{37, "+", ">", 0, 17, 17, 32, `zxtrnlhbaeimqswyA`},
		&testCase{38, "+", "<", 31, 17, 18, 32, `BzxtrnlhbaeimqswyA`},
		&testCase{39, "+", ">", 31, 18, 19, 32, `BzxtrnlhbaeimqswyAC`},
		&testCase{40, "----", "<", 0, 18, 18, 32, `zxtrnlhbaeimqswyAC`},
		&testCase{41, "----", ">", 0, 17, 17, 32, `zxtrnlhbaeimqswyA`},
		&testCase{42, "+", "<", 31, 17, 18, 32, `DzxtrnlhbaeimqswyA`},
		&testCase{43, "+", ">", 31, 18, 19, 32, `DzxtrnlhbaeimqswyAE`},
		&testCase{44, "+", "<", 30, 18, 20, 32, `FDzxtrnlhbaeimqswyAE`},
		&testCase{45, "+", ">", 30, 19, 21, 32, `FDzxtrnlhbaeimqswyAEG`},
		&testCase{46, "+", "<", 29, 19, 22, 32, `HFDzxtrnlhbaeimqswyAEG`},
		&testCase{47, "+", ">", 29, 20, 23, 32, `HFDzxtrnlhbaeimqswyAEGI`},
		&testCase{48, "+", "<", 28, 20, 24, 32, `JHFDzxtrnlhbaeimqswyAEGI`},
		&testCase{49, "+", ">", 28, 21, 25, 32, `JHFDzxtrnlhbaeimqswyAEGIK`},
		&testCase{50, "+", "<", 27, 21, 26, 32, `LJHFDzxtrnlhbaeimqswyAEGIK`},
		&testCase{51, "+", ">", 27, 22, 27, 32, `LJHFDzxtrnlhbaeimqswyAEGIKM`},
		&testCase{52, "+", "<", 26, 22, 28, 32, `NLJHFDzxtrnlhbaeimqswyAEGIKM`},
		&testCase{53, "+", ">", 26, 23, 29, 32, `NLJHFDzxtrnlhbaeimqswyAEGIKMO`},
		&testCase{54, "+", "<", 25, 23, 30, 32, `PNLJHFDzxtrnlhbaeimqswyAEGIKMO`},
		&testCase{55, "+", ">", 25, 24, 31, 32, `PNLJHFDzxtrnlhbaeimqswyAEGIKMOQ`},
		&testCase{56, "+", "<", 0, 32, 32, 64, `RPNLJHFDzxtrnlhbaeimqswyAEGIKMOQ`},
		&testCase{57, "+", ">", 0, 33, 33, 64, `RPNLJHFDzxtrnlhbaeimqswyAEGIKMOQS`},
		&testCase{58, "----", "<", 1, 33, 32, 64, `PNLJHFDzxtrnlhbaeimqswyAEGIKMOQS`},
		&testCase{59, "----", ">", 1, 32, 31, 64, `PNLJHFDzxtrnlhbaeimqswyAEGIKMOQ`},
		&testCase{60, "+", "<", 0, 32, 32, 64, `TPNLJHFDzxtrnlhbaeimqswyAEGIKMOQ`},
		&testCase{61, "+", ">", 0, 33, 33, 64, `TPNLJHFDzxtrnlhbaeimqswyAEGIKMOQU`},
		&testCase{62, "+", "<", 63, 33, 34, 64, `VTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQU`},
		&testCase{63, "+", ">", 63, 34, 35, 64, `VTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUW`},
		&testCase{64, "+", "<", 62, 34, 36, 64, `XVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUW`},
		&testCase{65, "+", ">", 62, 35, 37, 64, `XVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY`},
		&testCase{66, "+", "<", 61, 35, 38, 64, `ZXVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY`},
		&testCase{67, "+", ">", 61, 36, 39, 64, `ZXVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY0`},
		&testCase{68, "+", "<", 60, 36, 40, 64, `1ZXVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY0`},
		&testCase{69, "+", ">", 60, 37, 41, 64, `1ZXVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY02`},
		&testCase{70, "+", "<", 59, 37, 42, 64, `31ZXVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY02`},
		&testCase{71, "----", "<", 60, 37, 41, 64, `1ZXVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY02`},
		&testCase{72, "----", ">", 60, 36, 40, 64, `1ZXVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY0`},
		&testCase{73, "----", "<", 61, 36, 39, 64, `ZXVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY0`},
		&testCase{74, "----", ">", 61, 35, 38, 64, `ZXVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY`},
		&testCase{75, "----", "<", 62, 35, 37, 64, `XVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUWY`},
		&testCase{76, "----", ">", 62, 34, 36, 64, `XVTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUW`},
		&testCase{77, "----", "<", 63, 34, 35, 64, `VTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQUW`},
		&testCase{78, "----", ">", 63, 33, 34, 64, `VTPNLJHFDzxtrnlhbaeimqswyAEGIKMOQU`},
		&testCase{79, "----", "<", 0, 33, 33, 64, `TPNLJHFDzxtrnlhbaeimqswyAEGIKMOQU`},
		&testCase{80, "----", ">", 0, 32, 32, 64, `TPNLJHFDzxtrnlhbaeimqswyAEGIKMOQ`},
		&testCase{81, "----", "<", 1, 32, 31, 64, `PNLJHFDzxtrnlhbaeimqswyAEGIKMOQ`},
		&testCase{82, "----", ">", 1, 31, 30, 64, `PNLJHFDzxtrnlhbaeimqswyAEGIKMO`},
		&testCase{83, "+", ">", 1, 32, 31, 64, `PNLJHFDzxtrnlhbaeimqswyAEGIKMO4`},
		&testCase{84, "+", "<", 0, 32, 32, 64, `5PNLJHFDzxtrnlhbaeimqswyAEGIKMO4`},
		&testCase{85, "+", ">", 0, 33, 33, 64, `5PNLJHFDzxtrnlhbaeimqswyAEGIKMO46`},
		&testCase{86, "+", "<", 63, 33, 34, 64, `75PNLJHFDzxtrnlhbaeimqswyAEGIKMO46`},
		&testCase{87, "----", "<", 0, 33, 33, 64, `5PNLJHFDzxtrnlhbaeimqswyAEGIKMO46`},
		&testCase{88, "----", ">", 0, 32, 32, 64, `5PNLJHFDzxtrnlhbaeimqswyAEGIKMO4`},
		&testCase{89, "+", ">", 0, 33, 33, 64, `5PNLJHFDzxtrnlhbaeimqswyAEGIKMO48`},
		&testCase{90, "+", "<", 63, 33, 34, 64, `95PNLJHFDzxtrnlhbaeimqswyAEGIKMO48`},
		&testCase{91, "----", "<", 0, 33, 33, 64, `5PNLJHFDzxtrnlhbaeimqswyAEGIKMO48`},
		&testCase{92, "----", ">", 0, 32, 32, 64, `5PNLJHFDzxtrnlhbaeimqswyAEGIKMO4`},
		&testCase{93, "----", "<", 1, 32, 31, 64, `PNLJHFDzxtrnlhbaeimqswyAEGIKMO4`},
		&testCase{94, "----", ">", 1, 31, 30, 64, `PNLJHFDzxtrnlhbaeimqswyAEGIKMO`},
		&testCase{95, "----", "<", 2, 31, 29, 64, `NLJHFDzxtrnlhbaeimqswyAEGIKMO`},
		&testCase{96, "----", ">", 2, 30, 28, 64, `NLJHFDzxtrnlhbaeimqswyAEGIKM`},
		&testCase{97, "----", "<", 3, 30, 27, 64, `LJHFDzxtrnlhbaeimqswyAEGIKM`},
		&testCase{98, "----", ">", 3, 29, 26, 64, `LJHFDzxtrnlhbaeimqswyAEGIK`},
		&testCase{99, "----", "<", 4, 29, 25, 64, `JHFDzxtrnlhbaeimqswyAEGIK`},
		&testCase{100, "----", ">", 4, 28, 24, 64, `JHFDzxtrnlhbaeimqswyAEGI`},
		&testCase{101, "----", "<", 5, 28, 23, 64, `HFDzxtrnlhbaeimqswyAEGI`},
		&testCase{102, "----", ">", 5, 27, 22, 64, `HFDzxtrnlhbaeimqswyAEG`},
		&testCase{103, "+", ">", 5, 28, 23, 64, `HFDzxtrnlhbaeimqswyAEG0`},
		&testCase{104, "+", "<", 4, 28, 24, 64, `~HFDzxtrnlhbaeimqswyAEG0`},
		&testCase{105, "+", ">", 4, 29, 25, 64, `~HFDzxtrnlhbaeimqswyAEG0!`},
		&testCase{106, "+", "<", 3, 29, 26, 64, `@~HFDzxtrnlhbaeimqswyAEG0!`},
		&testCase{107, "+", ">", 3, 30, 27, 64, `@~HFDzxtrnlhbaeimqswyAEG0!#`},
		&testCase{108, "+", "<", 2, 30, 28, 64, `$@~HFDzxtrnlhbaeimqswyAEG0!#`},
		&testCase{109, "+", ">", 2, 31, 29, 64, `$@~HFDzxtrnlhbaeimqswyAEG0!#%`},
		&testCase{110, "+", "<", 1, 31, 30, 64, `^$@~HFDzxtrnlhbaeimqswyAEG0!#%`},
		&testCase{111, "+", ">", 1, 32, 31, 64, `^$@~HFDzxtrnlhbaeimqswyAEG0!#%&`},
		&testCase{112, "+", "<", 0, 32, 32, 64, `*^$@~HFDzxtrnlhbaeimqswyAEG0!#%&`},
		&testCase{113, "----", "<", 1, 32, 31, 64, `^$@~HFDzxtrnlhbaeimqswyAEG0!#%&`},
		&testCase{114, "----", ">", 1, 31, 30, 64, `^$@~HFDzxtrnlhbaeimqswyAEG0!#%`},
		&testCase{115, "+", ">", 1, 32, 31, 64, `^$@~HFDzxtrnlhbaeimqswyAEG0!#%(`},
		&testCase{116, "+", "<", 0, 32, 32, 64, `)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(`},
		&testCase{117, "+", ">", 0, 33, 33, 64, `)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(_`},
		&testCase{118, "+", "<", 63, 33, 34, 64, `+)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(_`},
		&testCase{119, "----", "<", 0, 33, 33, 64, `)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(_`},
		&testCase{120, "----", ">", 0, 32, 32, 64, `)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(`},
		&testCase{121, "+", ">", 0, 33, 33, 64, `)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|`},
		&testCase{122, "+", "<", 63, 33, 34, 64, `{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|`},
		&testCase{123, "+", ">", 63, 34, 35, 64, `{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|}`},
		&testCase{124, "+", "<", 62, 34, 36, 64, `[{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|}`},
		&testCase{125, "----", "<", 63, 34, 35, 64, `{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|}`},
		&testCase{126, "----", ">", 63, 33, 34, 64, `{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|`},
		&testCase{127, "+", ">", 63, 34, 35, 64, `{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|]`},
		&testCase{128, "+", "<", 62, 34, 36, 64, `?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|]`},
		&testCase{129, "+", ">", 62, 35, 37, 64, `?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|]<`},
		&testCase{130, "+", "<", 61, 35, 38, 64, `>?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|]<`},
		&testCase{131, "----", "<", 62, 35, 37, 64, `?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|]<`},
		&testCase{132, "----", ">", 62, 34, 36, 64, `?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|]`},
		&testCase{133, "+", ">", 62, 35, 37, 64, `?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],`},
		&testCase{134, "+", "<", 61, 35, 38, 64, `.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],`},
		&testCase{135, "+", ">", 61, 36, 39, 64, `.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],:`},
		&testCase{136, "+", "<", 60, 36, 40, 64, `;.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],:`},
		&testCase{137, "+", ">", 60, 37, 41, 64, `;.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],:'`},
		&testCase{138, "+", "<", 59, 37, 42, 64, `";.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],:'`},
		&testCase{139, "----", "<", 60, 37, 41, 64, `;.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],:'`},
		&testCase{140, "----", ">", 60, 36, 40, 64, `;.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],:`},
		&testCase{141, "+", ">", 60, 37, 41, 64, `;.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],:/`},
		&testCase{142, "+", "<", 59, 37, 42, 64, `\;.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],:/`},
	}
	q := newNodeDeque(5)
	assert(q.Empty() && q.Size() == 0 && q.Cap() == 5)
	for i, v := range cases {
		switch v.Op {
		case "+":
			n := testNewNode(nextCh())
			switch v.Dir {
			case ">":
				q.PushBack(n)
			case "<":
				q.PushFront(n)
			default:
				panic("")
			}
		case "----":
			switch v.Dir {
			case ">":
				q.PopBack()
			case "<":
				q.PopFront()
			default:
				panic("")
			}

		default:
			panic("")
		}
		verify(q, i, v)
	}
	testVerify(q, 999, 59, 37, 42, 64, `\;.?{)^$@~HFDzxtrnlhbaeimqswyAEG0!#%(|],:/`, true)
	assert(!q.Empty() && q.Size() != 0)
	q.Clear()
	assert(q.Empty() && q.Size() == 0)
	testVerify(q, 999, 0, 0, 0, 64, "", true)
}

func TestLargeDeque(t *testing.T) {
	tab := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890~!@#$%^&*()_+|{}[]?<>,.:;/`
	tbIdx := 0
	nextCh := func() string {
		idx := tbIdx % len(tab)
		ret := tab[idx : idx+1]
		tbIdx++
		return ret
	}

	q := newNodeDeque(7)
	for i := 0; i < dqTooLarge/2+7; i++ {
		q.PushBack(testNewNode(nextCh()))
		q.PushFront(testNewNode(nextCh()))
	}
	testVerify(q, -1, 8185, 4103, 4110, 8192, "pnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdb/:,<]}|_(&%#!086420YWUSQOMKIGECAywusqomkigeca;.>?[{+)*^$@~97531ZXVTRPNLJHFDBzxvtrpnljhfdbacegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmoqsuwyACEGIKMOQSUWY024680!#%&(_|}]<,:/bdfhjlnprtvxzBDFHJLNPRTVXZ13579~@$^*)+{[?>.;acegikmo",
		true)
}
