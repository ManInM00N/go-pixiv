package ToolFunc

func StringToInt64(s string) int64 {
	var num int64 = 0
	for i := 0; i < len(s); i++ {
		num = num*10 + int64(s[i]-'0')
	}
	return num
}
func Int64ToString(num int64) string {
	var s string = ""
	for ; num > 0; num /= 10 {
		s = s + string(num%10+'0')
	}
	return s
}
