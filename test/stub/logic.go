package stub

var SexExchangeFunc = func(input Teacher) (output Teacher) {
	output = input
	output.StudentData.Sex = !input.StudentData.Sex
	return
}

var DoubleFunc = func(input NumRequest) (output NumResponse) {
	output.Num = input.Num * 2
	return
}