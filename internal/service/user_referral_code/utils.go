package user_referral_code

var CODE_SET = []rune{
    '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
    'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}
const (
    PRIME_SET = 3       //与CODE_SET长度34互质的数
    CODE_LENGTH = 6     //邀请码固定长度6位
    PRIME_CODE = 5      //与邀请码长度互质
    SALT = 3291301845   //一个固定随机数
)

func GenerateReferralCode(seed uint64) string {
    // 放大 + 加盐
	seed = seed * PRIME_SET + SALT

	var code []rune
	slIdx := make([]byte, CODE_LENGTH)

	// 扩散
	for i := 0; i < CODE_LENGTH; i++ {
		slIdx[i] = byte(seed % uint64(len(CODE_SET)))                           // 获取 34 进制的每一位值
		slIdx[i] = (slIdx[i] + byte(i) * slIdx[0]) % byte(len(CODE_SET))        // 其他位与个位加和再取余（让个位的变化影响到所有位）
		seed = seed / uint64(len(CODE_SET))                                     // 相当于右移一位(34进制)
	}

	// 混淆
	for i := 0; i < CODE_LENGTH; i++ {
		idx := (byte(i) * PRIME_CODE) % byte(CODE_LENGTH)
		code = append(code, CODE_SET[slIdx[idx]])
	}
	return string(code) 
}
