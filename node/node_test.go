package node


import "testing"

func TestSample(t *testing.T) {
	t.Log(" test log 1!")
	t.Fail()

	t.Log(" test log 2!")
	t.Fail()
}

//test table sample
var tests = []struct{   // Test table
	in  string
	out string

}{
	{"in1", "exp1"},
	{"in2", "exp2"},
	{"in3", "exp3"},
}


func verify(t *testing.T, testnum int, testcase, input, output, expected string) {
	if expected != output {
		t.Errorf("%d. %s with input = %s: output %s != %s", testnum, testcase, input, output, expected)
	}
}


func TestFunction(t *testing.T) {
	//for i, tt := range tests {
		//s := FuncToBeTested(tt.in)
		//verify(t, i, “FuncToBeTested: “, tt.in, s, tt.out)
	//}
}