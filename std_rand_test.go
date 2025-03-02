// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

package rand_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/gozelle/rand"
	"io"
	"math"
	"os"
	"runtime"
	"testing"
	"testing/iotest"
)

const (
	numTestSamples = 10000
)

var (
	printtables = flag.Bool("printtables", false, "print ziggurat tables")
)

var rn, kn, wn, fn = rand.GetNormalDistributionParameters()
var re, ke, we, fe = rand.GetExponentialDistributionParameters()

type statsResults struct {
	mean        float64
	stddev      float64
	closeEnough float64
	maxError    float64
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func nearEqual(a, b, closeEnough, maxError float64) bool {
	absDiff := math.Abs(a - b)
	if absDiff < closeEnough { // Necessary when one value is zero and one value is close to zero.
		return true
	}
	return absDiff/(max(math.Abs(a), math.Abs(b))+math.SmallestNonzeroFloat64) < maxError
}

var testSeeds = []int64{1, 1754801282, 1698661970, 1550503961}

// checkSimilarDistribution returns success if the mean and stddev of the
// two statsResults are similar.
func (sr *statsResults) checkSimilarDistribution(expected *statsResults) error {
	if !nearEqual(sr.mean, expected.mean, expected.closeEnough, expected.maxError) {
		s := fmt.Sprintf("mean %v != %v (allowed error %v, %v)", sr.mean, expected.mean, expected.closeEnough, expected.maxError)
		fmt.Println(s)
		return errors.New(s)
	}
	if !nearEqual(sr.stddev, expected.stddev, expected.closeEnough, expected.maxError) {
		s := fmt.Sprintf("stddev %v != %v (allowed error %v, %v)", sr.stddev, expected.stddev, expected.closeEnough, expected.maxError)
		fmt.Println(s)
		return errors.New(s)
	}
	return nil
}

func getStatsResults(samples []float64) *statsResults {
	res := new(statsResults)
	var sum, squaresum float64
	for _, s := range samples {
		sum += s
		squaresum += s * s
	}
	res.mean = sum / float64(len(samples))
	res.stddev = math.Sqrt(squaresum/float64(len(samples)) - res.mean*res.mean)
	return res
}

func checkSampleDistribution(t *testing.T, samples []float64, expected *statsResults) {
	t.Helper()
	actual := getStatsResults(samples)
	err := actual.checkSimilarDistribution(expected)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func checkSampleSliceDistributions(t *testing.T, samples []float64, nslices int, expected *statsResults) {
	t.Helper()
	chunk := len(samples) / nslices
	for i := 0; i < nslices; i++ {
		low := i * chunk
		var high int
		if i == nslices-1 {
			high = len(samples) - 1
		} else {
			high = (i + 1) * chunk
		}
		checkSampleDistribution(t, samples[low:high], expected)
	}
}

//
// Normal distribution tests
//

func generateNormalSamples(nsamples int, mean, stddev float64, seed int64) []float64 {
	r := rand.New(uint64(seed))
	samples := make([]float64, nsamples)
	for i := range samples {
		samples[i] = r.NormFloat64()*stddev + mean
	}
	return samples
}

func testNormalDistribution(t *testing.T, nsamples int, mean, stddev float64, seed int64) {
	//fmt.Printf("testing nsamples=%v mean=%v stddev=%v seed=%v\n", nsamples, mean, stddev, seed);
	
	samples := generateNormalSamples(nsamples, mean, stddev, seed)
	errorScale := max(1.0, stddev) // Error scales with stddev
	expected := &statsResults{mean, stddev, 0.10 * errorScale, 0.08 * errorScale}
	
	// Make sure that the entire set matches the expected distribution.
	checkSampleDistribution(t, samples, expected)
	
	// Make sure that each half of the set matches the expected distribution.
	checkSampleSliceDistributions(t, samples, 2, expected)
	
	// Make sure that each 7th of the set matches the expected distribution.
	checkSampleSliceDistributions(t, samples, 7, expected)
}

// Actual tests

func TestStandardNormalValues(t *testing.T) {
	for _, seed := range testSeeds {
		testNormalDistribution(t, numTestSamples, 0, 1, seed)
	}
}

func TestNonStandardNormalValues(t *testing.T) {
	sdmax := 1000.0
	mmax := 1000.0
	if testing.Short() {
		sdmax = 5
		mmax = 5
	}
	for sd := 0.5; sd < sdmax; sd *= 2 {
		for m := 0.5; m < mmax; m *= 2 {
			for _, seed := range testSeeds {
				testNormalDistribution(t, numTestSamples, m, sd, seed)
				if testing.Short() {
					break
				}
			}
		}
	}
}

//
// Exponential distribution tests
//

func generateExponentialSamples(nsamples int, rate float64, seed int64) []float64 {
	r := rand.New(uint64(seed))
	samples := make([]float64, nsamples)
	for i := range samples {
		samples[i] = r.ExpFloat64() / rate
	}
	return samples
}

func testExponentialDistribution(t *testing.T, nsamples int, rate float64, seed int64) {
	//fmt.Printf("testing nsamples=%v rate=%v seed=%v\n", nsamples, rate, seed);
	
	mean := 1 / rate
	stddev := mean
	
	samples := generateExponentialSamples(nsamples, rate, seed)
	errorScale := max(1.0, 1/rate) // Error scales with the inverse of the rate
	expected := &statsResults{mean, stddev, 0.10 * errorScale, 0.20 * errorScale}
	
	// Make sure that the entire set matches the expected distribution.
	checkSampleDistribution(t, samples, expected)
	
	// Make sure that each half of the set matches the expected distribution.
	checkSampleSliceDistributions(t, samples, 2, expected)
	
	// Make sure that each 7th of the set matches the expected distribution.
	checkSampleSliceDistributions(t, samples, 7, expected)
}

// Actual tests

func TestStandardExponentialValues(t *testing.T) {
	for _, seed := range testSeeds {
		testExponentialDistribution(t, numTestSamples, 1, seed)
	}
}

func TestNonStandardExponentialValues(t *testing.T) {
	for rate := 0.05; rate < 10; rate *= 2 {
		for _, seed := range testSeeds {
			testExponentialDistribution(t, numTestSamples, rate, seed)
			if testing.Short() {
				break
			}
		}
	}
}

//
// Table generation tests
//

func initNorm() (testKn []uint64, testWn, testFn []float64) {
	const (
		m1 = 1 << 52
		vn = 0.00492867323399
	)
	var (
		dn = rn
		tn = dn
	)
	
	testKn = make([]uint64, 256)
	testWn = make([]float64, 256)
	testFn = make([]float64, 256)
	
	q := vn / math.Exp(-0.5*dn*dn)
	testKn[0] = uint64((dn / q) * m1)
	testKn[1] = 0
	testWn[0] = q / m1
	testWn[255] = dn / m1
	testFn[0] = 1.0
	testFn[255] = math.Exp(-0.5 * dn * dn)
	for i := 254; i >= 1; i-- {
		dn = math.Sqrt(-2.0 * math.Log(vn/dn+math.Exp(-0.5*dn*dn)))
		testKn[i+1] = uint64((dn / tn) * m1)
		tn = dn
		testFn[i] = math.Exp(-0.5 * dn * dn)
		testWn[i] = dn / m1
	}
	return
}

func initExp() (testKe []uint64, testWe, testFe []float64) {
	const (
		m2 = 1 << 53
		ve = 0.0039496598225815571993
	)
	var (
		de = re
		te = de
	)
	
	testKe = make([]uint64, 256)
	testWe = make([]float64, 256)
	testFe = make([]float64, 256)
	
	q := ve / math.Exp(-de)
	testKe[0] = uint64((de / q) * m2)
	testKe[1] = 0
	testWe[0] = q / m2
	testWe[255] = de / m2
	testFe[0] = 1.0
	testFe[255] = math.Exp(-de)
	for i := 254; i >= 1; i-- {
		de = -math.Log(ve/de + math.Exp(-de))
		testKe[i+1] = uint64((de / te) * m2)
		te = de
		testFe[i] = math.Exp(-de)
		testWe[i] = de / m2
	}
	return
}

// compareUint64Slices returns the first index where the two slices
// disagree, or <0 if the lengths are the same and all elements
// are close to each other.
func compareUint64Slices(s1, s2 []uint64) int {
	if len(s1) != len(s2) {
		if len(s1) > len(s2) {
			return len(s2) + 1
		}
		return len(s1) + 1
	}
	for i := range s1 {
		if !nearEqual(float64(s1[i]), float64(s2[i]), 0, 1e-12) {
			return i
		}
	}
	return -1
}

// compareFloat64Slices returns the first index where the two slices
// disagree, or <0 if the lengths are the same and all elements
// are close to each other.
func compareFloat64Slices(s1, s2 []float64) int {
	if len(s1) != len(s2) {
		if len(s1) > len(s2) {
			return len(s2) + 1
		}
		return len(s1) + 1
	}
	for i := range s1 {
		if !nearEqual(s1[i], s2[i], 0, 1e-12) {
			return i
		}
	}
	return -1
}

func TestNormTables(t *testing.T) {
	testKn, testWn, testFn := initNorm()
	if *printtables {
		fmt.Printf("var kn = %#v\n", testKn)
		fmt.Printf("var wn = %#v\n", testWn)
		fmt.Printf("var fn = %#v\n", testFn)
		return
	}
	
	if i := compareUint64Slices(kn[0:], testKn); i >= 0 {
		t.Errorf("kn disagrees at index %v; %v != %v", i, kn[i], testKn[i])
	}
	if i := compareFloat64Slices(wn[0:], testWn); i >= 0 {
		t.Errorf("wn disagrees at index %v; %v != %v", i, wn[i], testWn[i])
	}
	if i := compareFloat64Slices(fn[0:], testFn); i >= 0 {
		t.Errorf("fn disagrees at index %v; %v != %v", i, fn[i], testFn[i])
	}
}

func TestExpTables(t *testing.T) {
	testKe, testWe, testFe := initExp()
	if *printtables {
		fmt.Printf("var ke = %#v\n", testKe)
		fmt.Printf("var we = %#v\n", testWe)
		fmt.Printf("var fe = %#v\n", testFe)
		return
	}
	
	if i := compareUint64Slices(ke[0:], testKe); i >= 0 {
		t.Errorf("ke disagrees at index %v; %v != %v", i, ke[i], testKe[i])
	}
	if i := compareFloat64Slices(we[0:], testWe); i >= 0 {
		t.Errorf("we disagrees at index %v; %v != %v", i, we[i], testWe[i])
	}
	if i := compareFloat64Slices(fe[0:], testFe); i >= 0 {
		t.Errorf("fe disagrees at index %v; %v != %v", i, fe[i], testFe[i])
	}
}

func hasSlowFloatingPoint() bool {
	switch runtime.GOARCH {
	case "arm":
		return os.Getenv("GOARM") == "5"
	case "mips", "mipsle", "mips64", "mips64le":
		// Be conservative and assume that all mips boards
		// have emulated floating point.
		// TODO: detect what it actually has.
		return true
	}
	return false
}

func TestFloat32_6721(t *testing.T) {
	// For issue 6721, the problem came after 7533753 calls, so check 10e6.
	num := int(10e6)
	// But do the full amount only on builders (not locally).
	// But ARM5 floating point emulation is slow (Issue 10749), so
	// do less for that builder:
	if testing.Short() && hasSlowFloatingPoint() {
		num /= 100 // 1.72 seconds instead of 172 seconds
	}
	
	r := rand.New(1)
	for ct := 0; ct < num; ct++ {
		f := r.Float32()
		if f >= 1 {
			t.Fatal("Float32() should be in range [0,1). ct:", ct, "f:", f)
		}
	}
}

func testReadUniformity(t *testing.T, n int, seed int64) {
	r := rand.New(uint64(seed))
	buf := make([]byte, n)
	nRead, err := r.Read(buf)
	if err != nil {
		t.Errorf("Read err %v", err)
	}
	if nRead != n {
		t.Errorf("Read returned unexpected n; %d != %d", nRead, n)
	}
	
	// Expect a uniform distribution of byte values, which lie in [0, 255].
	var (
		mean       = 255.0 / 2
		stddev     = 256.0 / math.Sqrt(12.0)
		errorScale = stddev / math.Sqrt(float64(n))
	)
	
	expected := &statsResults{mean, stddev, 0.10 * errorScale, 0.08 * errorScale}
	
	// Cast bytes as floats to use the common distribution-validity checks.
	samples := make([]float64, n)
	for i, val := range buf {
		samples[i] = float64(val)
	}
	// Make sure that the entire set matches the expected distribution.
	checkSampleDistribution(t, samples, expected)
}

func TestReadUniformity(t *testing.T) {
	testBufferSizes := []int{
		2, 4, 7, 64, 1024, 1 << 16, 1 << 20,
	}
	for _, seed := range testSeeds {
		for _, n := range testBufferSizes {
			testReadUniformity(t, n, seed)
		}
	}
}

func TestReadEmpty(t *testing.T) {
	r := rand.New(1)
	buf := make([]byte, 0)
	n, err := r.Read(buf)
	if err != nil {
		t.Errorf("Read err into empty buffer; %v", err)
	}
	if n != 0 {
		t.Errorf("Read into empty buffer returned unexpected n of %d", n)
	}
}

func TestReadByOneByte(t *testing.T) {
	r := rand.New(1)
	b1 := make([]byte, 100)
	_, err := io.ReadFull(iotest.OneByteReader(r), b1)
	if err != nil {
		t.Errorf("read by one byte: %v", err)
	}
	r = rand.New(1)
	b2 := make([]byte, 100)
	_, err = r.Read(b2)
	if err != nil {
		t.Errorf("read: %v", err)
	}
	if !bytes.Equal(b1, b2) {
		t.Errorf("read by one byte vs single read:\n%x\n%x", b1, b2)
	}
}

func TestReadSeedReset(t *testing.T) {
	r := rand.New(42)
	b1 := make([]byte, 128)
	_, err := r.Read(b1)
	if err != nil {
		t.Errorf("read: %v", err)
	}
	r.Seed(42)
	b2 := make([]byte, 128)
	_, err = r.Read(b2)
	if err != nil {
		t.Errorf("read: %v", err)
	}
	if !bytes.Equal(b1, b2) {
		t.Errorf("mismatch after re-seed:\n%x\n%x", b1, b2)
	}
}

func TestShuffleSmall(t *testing.T) {
	// Check that Shuffle allows n=0 and n=1, but that swap is never called for them.
	r := rand.New(1)
	for n := 0; n <= 1; n++ {
		r.Shuffle(n, func(i, j int) { t.Fatalf("swap called, n=%d i=%d j=%d", n, i, j) })
	}
}

// encodePerm converts from a permuted slice of length n, such as Perm generates, to an int in [0, n!).
// See https://en.wikipedia.org/wiki/Lehmer_code.
// encodePerm modifies the input slice.
func encodePerm(s []int) int {
	// Convert to Lehmer code.
	for i, x := range s {
		r := s[i+1:]
		for j, y := range r {
			if y > x {
				r[j]--
			}
		}
	}
	// Convert to int in [0, n!).
	m := 0
	fact := 1
	for i := len(s) - 1; i >= 0; i-- {
		m += s[i] * fact
		fact *= len(s) - i
	}
	return m
}

// TestUniformFactorial tests several ways of generating a uniform value in [0, n!).
func TestUniformFactorial(t *testing.T) {
	r := rand.New(uint64(testSeeds[0]))
	top := 6
	if testing.Short() {
		top = 3
	}
	for n := 3; n <= top; n++ {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			// Calculate n!.
			nfact := 1
			for i := 2; i <= n; i++ {
				nfact *= i
			}
			
			// Test a few different ways to generate a uniform distribution.
			p := make([]int, n) // re-usable slice for Shuffle generator
			tests := [...]struct {
				name string
				fn   func() int
			}{
				{name: "Int31n", fn: func() int { return int(r.Int31n(int32(nfact))) }},
				{name: "int31n", fn: func() int { return int(rand.Int31nForTest(r, int32(nfact))) }},
				{name: "Perm", fn: func() int { return encodePerm(r.Perm(n)) }},
				{name: "Shuffle", fn: func() int {
					// Generate permutation using Shuffle.
					for i := range p {
						p[i] = i
					}
					r.Shuffle(n, func(i, j int) { p[i], p[j] = p[j], p[i] })
					return encodePerm(p)
				}},
				{name: "ShuffleSliceGeneric", fn: func() int {
					if rand.ShuffleSliceGeneric == nil {
						return int(r.Uint32n(uint32(nfact)))
					}
					// Generate permutation using generic Shuffle.
					for i := range p {
						p[i] = i
					}
					rand.ShuffleSliceGeneric(r, p)
					return encodePerm(p)
				}},
			}
			
			for _, test := range tests {
				t.Run(test.name, func(t *testing.T) {
					// Gather chi-squared values and check that they follow
					// the expected normal distribution given n!-1 degrees of freedom.
					// See https://en.wikipedia.org/wiki/Pearson%27s_chi-squared_test and
					// https://www.johndcook.com/Beautiful_Testing_ch10.pdf.
					nsamples := 10 * nfact
					if nsamples < 200 {
						nsamples = 200
					}
					samples := make([]float64, nsamples)
					for i := range samples {
						// Generate some uniformly distributed values and count their occurrences.
						const iters = 1000
						counts := make([]int, nfact)
						for i := 0; i < iters; i++ {
							counts[test.fn()]++
						}
						// Calculate chi-squared and add to samples.
						want := iters / float64(nfact)
						var χ2 float64
						for _, have := range counts {
							err := float64(have) - want
							χ2 += err * err
						}
						χ2 /= want
						samples[i] = χ2
					}
					
					// Check that our samples approximate the appropriate normal distribution.
					dof := float64(nfact - 1)
					expected := &statsResults{mean: dof, stddev: math.Sqrt(2 * dof)}
					errorScale := max(1.0, expected.stddev)
					expected.closeEnough = 0.10 * errorScale
					expected.maxError = 0.08 // TODO: What is the right value here? See issue 21211.
					checkSampleDistribution(t, samples, expected)
				})
			}
		})
	}
}
