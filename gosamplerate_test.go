package gosamplerate

import (
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Testing Package gosamplerate")
	m.Run()
	log.Println("Tests Completed")
}

func TestConverter(t *testing.T) {
	tests := []struct {
		name       string
		args       int
		wanterr    bool
		errmessage string
	}{
		// Current Test
		{"Test Converter: Test Get Converter Name for SRC_LINEAR", SRC_LINEAR, false, ""},
		{"Test Converter: Test Get converter Name Error", 5, true, "unknown samplerate converter"},

		// Proposed Tests
		{"Test Converter: Test Get Name for SRC_ZERO_ORDER_HOLD", SRC_ZERO_ORDER_HOLD, false, ""},
		{"Test Converter: Test Get Name for SRC_SINC_FASTEST", SRC_SINC_FASTEST, false,""},
		{"Test Converter: Test Get Name for SRC_SINC_MEDIUM_QUALITY", SRC_SINC_MEDIUM_QUALITY, false, ""},
		{"Test Converter: Test Get Name for SRC_SINC_BEST_QUALITY", SRC_SINC_BEST_QUALITY, false, ""},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := GetName(C.int(test.args))
			if exist(err) != test.wanterr {
				t.Logf("Expected %v for Erros got Error:%v", test.wanterr, err)
				t.Fail()
			}

			if test.wanterr && test.errmessage != err.Error() {
				t.Logf("Unexpected error %v", err)
				t.Fail()
			}
		})
	}
}

func TestConverterDescription(t *testing.T) {
	tests := []struct {
		name       string
		args       int
		desc       string
		wanterr    bool
		errmessage string
	}{
		// Current Tests
		{"Test Description: Test Converter Description", SRC_LINEAR, "Linear interpolator, very fast, poor quality.", false, ""},
		{"Test Description: Test Description Error Message", SRC_ZERO_ORDER_HOLD, "Linear interpolator, very fast, poor quality.", true, "unknown samplerate converter"},

		// Proposed Tests
		{"Test Description: SRC_LINEAR", SRC_LINEAR, "Linear interpolator, very fast, poor quality.", false,""},
		{"Test Description: SRC_ZERO_ORDER_HOLD", SRC_ZERO_ORDER_HOLD, "Zero order hold, very fast, poor quality.", false, ""},
		{"Test Description: SRC_SINC_FASTEST", SRC_SINC_FASTEST, "Sinc interpolator (fastest), medium quality.", false, ""},
		{"Test Description: SRC_SINC_MEDIUM_QUALITY", SRC_SINC_MEDIUM_QUALITY, "Sinc interpolator (medium), good quality.", false, ""},
		{"Test Description: SRC_SINC_BEST_QUALITY", SRC_SINC_BEST_QUALITY, "Sinc interpolator (best quality), slow.", false, ""},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			desc, err := GetDescription(C.int(test.args))
			if exist(err) != test.wanterr {
				t.Logf("Expected %v for Errors got Error:%v", test.wanterr, err)
				t.Fail()
			}

			//test description
			if !reflect.DeepEqual(desc, test.desc) {
				t.Logf("Expected %v as Description got %v", test.desc, desc)
			}

			//test error message
			if test.wanterr && test.errmessage != err.Error() {
				t.Logf("Unexpected error %v", err)
				t.Fail()
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if !strings.Contains(version, "libsamplerate-") {
		t.Logf("Error: Unexpected string")
		t.Fail()
	}
}

func TestInitAndDestroy(t *testing.T) {
	type args struct {
		channelno     int
		convertertype int
		buffer        int
	}
	tests := []struct {
		name    string
		wanterr bool
		errmsg  string
		args
	}{
		{"Test Init And Destroy: Normal", false, "", args{2, SRC_LINEAR, 100}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init
			src, err := New(test.convertertype, test.channelno, test.buffer) //i should probably test the full suite of SRC_SINK
			if exist(err) != test.wanterr {
				t.Logf("Expected %v for errors, got Error:%v, Failed to create converter object", test.wanterr, err)
				t.Fail()
			}

			// retrieve channels
			chans, err := src.GetChannels()
			if exist(err) != test.wanterr {
				t.Logf("Expected %v for errors, got Error:%v, Failed to retrieve channels", test.wanterr, err)
				t.Fail()
			}

			// verify channels
			if test.channelno != chans {
				t.Logf("Expected %v channels got Error:%v, Failed channel no verification", test.channelno, chans)
				t.Fail()
			}

			// Reset converter
			if err := src.Reset(); exist(err) != test.wanterr {
				t.Logf("Expected %v for errors, got Error:%v, Failed to reset converter", test.wanterr, err)
				t.Fail()
			}

			// verify mem de-alloc
			if err := Delete(src); exist(err) != test.wanterr {
				t.Logf("Expected %v for errors, got Error:%v, Failed to cleanup converter object", test.wanterr, err)
				t.Fail()
			}

		})

	}
}

// Here I choose to replace TestInvalidSrcObject with a more comprehensive test TestNew

func TestNew(t *testing.T) {
	type args struct {
		channelno     int
		convertertype int
		buffer        int
	}
	tests := []struct {
		name    string
		wanterr bool
		errmsg  string
		args
	}{
		{"Test New: Invalid Converter Object", true, "Could not initialize samplerate converter object", args{5, 2, 100}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := New(test.convertertype, test.channelno, test.buffer)

			if exist(err) != test.wanterr {
				if test.wanterr && err.Error() != test.errmsg {
					t.Logf("Expected %v for errors, got Error:%v, Failed to create converter object", test.wanterr, err)
					t.Fail()
				}

				t.Logf("Expected %v for errors, got Error:%v", test.wanterr, err)
				t.Fail()
			}

		})
	}
}

func TestSimple(t *testing.T) {
	type (
		args struct {
			channels      int
			convertertype int
			ratio         float64
		}
		values struct {
			datain   []float32
			expected []float32
		}
	)

	tests := []struct {
		name    string
		wanterr bool
		errmsg  string
		args
		values
	}{
		//current tests
		{"Test Simple: Normal", false, "", args{1, SRC_LINEAR, 1.5}, values{[]float32{0.1, -0.5, 0.3, 0.4, 0.1}, []float32{0.1, 0.1, -0.10000001, -0.5, 0.033333343, 0.33333334, 0.4, 0.2}}},
		{"Test Simple: Simple less than one", false, "", args{1, SRC_LINEAR, 0.5}, values{[]float32{}, []float32{0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3, 0.1, -0.5, 0.4, 0.1, 0.3}}},
		{"Test Simple: Test Error Message", true, "Error code: 6; SRC ratio outside [1/256, 256] range.", args{1, SRC_LINEAR, -5.3}, values{[]float32{0.1, 0.9}, []float32{}}},
	}

	//run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := Simple(test.datain, test.ratio, test.chans, test.convertertype)

			if exist(err) != test.wanterr {
				if test.wanterr && err.Error != test.errmsg {
					t.Logf("Expected Error message %v got  Error: %v", test.errmsg, err)
					t.Fail()
				}
				t.Logf("Expected %v for errors, got Error:%v", test.wanterr, err)
				t.Fail()
			}

			//test output
			if !closeEnough(output, test.expected) {
				t.Logf("Error: expected values and output do not match and its not close")
				t.Fail()
			}

		})

	}
}

func TestProcess(t *testing.T) {
	type (
		args struct {
			channels      int
			convertertype int
			buffer        int
			ratio         float64
			endOfInput    bool
		}
		values struct {
			datain   []float32
			expected []float32
		}
	)
	tests := []struct {
		name    string
		wanterr bool
		errmsg  string
		args
		values
	}{
		// current test
		{"Test Process: Normal Process", false, "", args{2, SRC_LINEAR, 100, 2.0, false}, values{[]float32{0.1, -0.5, 0.2, -0.3}, []float32{0.1, -0.5, 0.1, -0.5, 0.1, -0.5, 0.15, -0.4}}},
		{"Test Process: Process With End of Input Flag Set", false, "", args{}, values{}},
		{"Test Process: Process Data Slice Bigger than Input buffer", true, "data slice is larger than buffer", args{}, values{[]float32{0.1, -0.5, 0.2, -0.3}, []float32{0.11488709,
			-0.46334597, 0.18373828, -0.48996875, 0.1821644,
			-0.32879135, 0.10804618, -0.11150829}}},
		{"Test Process: Process Error With Invalid Ratio", true, "Error code: 6; SRC ratio outside [1/256, 256] range.", args{}, values{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			src, err := New(test.convertertype, test.channels, test.buffer)

			if exist(err) != test.wanterr {
				t.Logf("Expected %v for errors, got Error:%v, Failed to create converter object", test.wanterr, err)
				t.Fail()
			}

			output, err := src.Process(test.input, test.channels, test.endOfInput)

			if exist(err) != test.wanterr {
				if test.wanterr && err.Error != test.errmsg {
					t.Logf("Expected Error message %v got  Error: %v", test.errmsg, err)
					t.Fail()
				}
				t.Logf("Expected %v for errors, got Error:%v", test.wanterr, err)
				t.Fail()
			}

			// verify output values
			if !reflect.DeepEqual(output, test.expected) {
				t.Log("input:", test.input)
				t.Log("output:", output)
				t.Log("unexpected output")
				t.Fail()
			}

			// Deallocate memory
			err = Delete(src) // test.want err should not apply here
			if err != nil {
				t.Logf("Error: Failed to Deallocate converter object")
				t.Fail()
			}

		})
	}

}

func TestGetChannels(t *testing.T) {
	type args struct {
		channels      int
		convertertype int
		buffer        int
	}

	tests := []struct {
		name    string
		wanterr bool
		errmsg  string
		args
	}{
		{"Test Channels: Test", false, "", args{2, SRC_SINC_FASTEST, 100}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			src, err := New(test.convertertype, test.channels, test.buffer)
			if exist(err) != test.wanterr {
				if test.wanterr && test.errmsg != err.Error() {
					t.Logf("Expected Error message: %v got %v", test.errmsg, err)
					t.Fail()
				}

				t.Logf("Expected %v for Errors got Error: %v", test.wanterr, err)
				t.Fail()
			}

			chanlen, err := src.GetChannels()

			if exist(err) != test.wanterr {
				t.Logf("Expected %v for errors, got Error:%v", test.wanterr, err)
				t.Fail()
			}

			if chanlen != test.channels {
				t.Logf("Error: Expected %v for channel length got %v", test.channels, chanlen)
				t.Fail()
			}

		})
	}

}

func TestSetRatio(t *testing.T) {
	type args struct {
		convertertype int
		channels      int
		buffer        int
	}

	tests := []struct {
		name string
		args
		setratio float32
		wanterr  bool
		errmsg   string
	}{
		{},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			src, err := New(test.convertertype, test.channels, test.buffer)
			if exist(err) != test.wanterr {
				if test.wanterr && test.errmsg != err.Error() {
					t.Logf("Expected Error message: %v got %v", test.errmsg, err)
					t.Fail()
				}

				t.Logf("Expected %v for Errors got Error: %v", test.wanterr, err)
				t.Fail()
			}

			if err = src.SetRatio(test.setratio); err != nil {
				t.Logf("unexpected result; should be valid conversion rate")
				t.Fail()
			}
		})
	}

}

func TestIsValidRatio(t *testing.T) {
	tests := []struct {
		name    string
		value   int
		wanterr bool
	}{
		{"Test IsValidRatio: Normal", 5, false},
		{"Test IsValidRatio: Negative value", -1, true},
		{"Test IsValidRatio: Random value", 257, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !IsValidRatio(test.value) {
				t.Logf("Expected %v for validity got %v", test.wanterr, !test.wanterr)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	channels := 2
	src, err := New(SRC_SINC_FASTEST, channels, 100)
	if err != nil {
		t.Fatal(err)
	}

	errNo := src.ErrorNo()
	if errNo != 0 {
		t.Fatal("unexpected error number")
	}

	errString := Error(0)
	if errString != "No error." {
		t.Fatal("unexpected Error string")
	}

	err = Delete(src)
	if err != nil {
		t.Fatal(err)
	}
}

func closeEnough(a, b []float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v-b[i] > 0.00001 {
			return false
		}
	}
	return true
}

// useful error to bool converter
func exist(err error) bool {
	if err != nil {
		return true
	}
	return false
}
