package basis

import "testing"

func TestShenSha(t *testing.T) {
	// Tian Yi: Jia Day, find Chou or Wei
	if !GetTianYi(Jia, Chou) {
		t.Error("Jia should have TianYi at Chou")
	}
	if !GetTianYi(Jia, Wei) {
		t.Error("Jia should have TianYi at Wei")
	}

	// Yi Ma: Zi year/day, find Yin
	if !GetYiMa(Zi, YinB) {
		t.Error("Zi should have YiMa at YinB")
	}

	// Tao Hua: Zi year/day, find You
	if !GetTaoHua(Zi, You) {
		t.Error("Zi should have TaoHua at You")
	}
}
