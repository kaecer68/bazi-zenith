package basis

// ShenSha represents an auspicious or inauspicious symbolic star.
type ShenSha string

const (
	TianYi    ShenSha = "天乙貴人"
	YiMa      ShenSha = "驛馬"
	TaoHua    ShenSha = "桃花"
	LuShen    ShenSha = "祿神"
	YangRen   ShenSha = "羊刃"
	WenChang  ShenSha = "文昌"
	HuaGai    ShenSha = "華蓋"
	JiangXing ShenSha = "將星"
	HongLuan  ShenSha = "紅鸞"
	TianXi    ShenSha = "天喜"
)

// GetTianYi returns true if the branch is a Tian Yi star for the stem.
func GetTianYi(me Stem, b Branch) bool {
	// 口訣：甲戊並牛羊, 乙己鼠猴鄉, 丙丁豬雞位, 壬癸兔蛇藏, 庚辛逢馬虎
	switch me {
	case Jia, Wu:
		return b == Chou || b == Wei
	case Yi, Ji:
		return b == Zi || b == Shen
	case Bing, Ding:
		return b == Hai || b == You
	case Ren, Gui:
		return b == Mao || b == Si
	case Geng, Xin:
		return b == WuB || b == YinB
	}
	return false
}

// GetYiMa returns true if the target branch is an Yi Ma star relative to source branch.
func GetYiMa(source Branch, target Branch) bool {
	// 申子辰馬在寅, 寅午戌馬在申, 巳酉丑馬在亥, 亥卯未馬在巳
	switch source {
	case Shen, Zi, Chen:
		return target == YinB
	case YinB, WuB, Xu:
		return target == Shen
	case Si, You, Chou:
		return target == Hai
	case Hai, Mao, Wei:
		return target == Si
	}
	return false
}

// GetTaoHua returns true if the target branch is a Tao Hua star relative to source branch.
func GetTaoHua(source Branch, target Branch) bool {
	// 申子辰在酉, 寅午戌在卯, 巳酉丑在午, 亥卯未在子
	switch source {
	case Shen, Zi, Chen:
		return target == You
	case YinB, WuB, Xu:
		return target == Mao
	case Si, You, Chou:
		return target == WuB
	case Hai, Mao, Wei:
		return target == Zi
	}
	return false
}

// GetLuShen returns true if the branch is Lu Shen for the stem.
func GetLuShen(me Stem, b Branch) bool {
	// 甲祿在寅, 乙祿在卯, 丙戊祿在巳, 丁己祿在午, 庚祿在申, 辛祿在酉, 壬祿在亥, 癸祿在子
	switch me {
	case Jia:
		return b == YinB
	case Yi:
		return b == Mao
	case Bing, Wu:
		return b == Si
	case Ding, Ji:
		return b == WuB
	case Geng:
		return b == Shen
	case Xin:
		return b == You
	case Ren:
		return b == Hai
	case Gui:
		return b == Zi
	}
	return false
}

// GetYangRen returns true if the branch is Yang Ren for the stem.
func GetYangRen(me Stem, b Branch) bool {
	// 甲羊刃在卯, 丙戊羊刃在午, 庚羊刃在酉, 壬羊刃在子
	// (Typically for Yang stems, some schools include Yin stems)
	switch me {
	case Jia:
		return b == Mao
	case Bing, Wu:
		return b == WuB
	case Geng:
		return b == You
	case Ren:
		return b == Zi
		// Yin stems optional: 乙在辰, 丁己在未, 辛在戌, 癸在丑 (Simplified usually Omits)
	}
	return false
}

// GetWenChang returns true if the branch is Wen Chang for the stem.
func GetWenChang(me Stem, b Branch) bool {
	// 甲乙巳午報君知, 丙戊申宮丁己雞, 庚豬辛鼠壬逢虎, 癸人見兔入雲梯
	switch me {
	case Jia:
		return b == Si
	case Yi:
		return b == WuB
	case Bing, Wu:
		return b == Shen
	case Ding, Ji:
		return b == You
	case Geng:
		return b == Hai
	case Xin:
		return b == Zi
	case Ren:
		return b == YinB
	case Gui:
		return b == Mao
	}
	return false
}

// GetHuaGai returns true if target is Hua Gai relative to source.
func GetHuaGai(source Branch, target Branch) bool {
	// 申子辰在辰, 寅午戌在戌, 巳酉丑在丑, 亥卯未在未
	switch source {
	case Shen, Zi, Chen:
		return target == Chen
	case YinB, WuB, Xu:
		return target == Xu
	case Si, You, Chou:
		return target == Chou
	case Hai, Mao, Wei:
		return target == Wei
	}
	return false
}

// GetJiangXing returns true if target is Jiang Xing relative to source.
func GetJiangXing(source Branch, target Branch) bool {
	// 申子辰在子, 寅午戌在午, 巳酉丑在酉, 亥卯未在卯
	switch source {
	case Shen, Zi, Chen:
		return target == Zi
	case YinB, WuB, Xu:
		return target == WuB
	case Si, You, Chou:
		return target == You
	case Hai, Mao, Wei:
		return target == Mao
	}
	return false
}

// GetHongLuan returns true if target is Hong Luan relative to year branch.
func GetHongLuan(year Branch, target Branch) bool {
	// 鼠在卯, 牛在寅, 虎在丑, 兔在子, 龍在亥, 蛇在戌, 馬在酉, 羊在申, 猴在未, 雞在午, 狗在巳, 豬在辰
	branchOrder := []Branch{Zi, Chou, YinB, Mao, Chen, Si, WuB, Wei, Shen, You, Xu, Hai}
	for i, b := range branchOrder {
		if b == year {
			// (12 - (i - 1) + 12) % 12
			targetIdx := (16 - i) % 12 // Simplified formula: Zi(0) -> Mao(3)
			return branchOrder[targetIdx] == target
		}
	}
	return false
}

// GetTianXi returns true if target is Tian Xi relative to year branch (opposite to Hong Luan).
func GetTianXi(year Branch, target Branch) bool {
	branchOrder := []Branch{Zi, Chou, YinB, Mao, Chen, Si, WuB, Wei, Shen, You, Xu, Hai}
	for i, b := range branchOrder {
		if b == year {
			targetIdx := (16 - i + 6) % 12
			return branchOrder[targetIdx] == target
		}
	}
	return false
}
