// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

// Test that random number sequences generated by a specific seed
// do not change from version to version.
//
// Do NOT make changes to the golden outputs. If bugs need to be fixed
// in the underlying code, find ways to fix them that do not affect the
// outputs.

package rand_test

import (
	"flag"
	"fmt"
	"github.com/gozelle/rand"
	"reflect"
	"testing"
)

var (
	printgolden = flag.Bool("printgolden", false, "print golden results for regression test")
	skipregress = flag.Bool("skipregress", false, "skip the regression test")
)

func TestRegress(t *testing.T) {
	if *skipregress {
		t.Skip("-skipregress specified")
	}
	
	var int32s = []int32{1, 10, 32, 1 << 20, 1<<20 + 1, 1000000000, 1 << 30, 1<<31 - 2, 1<<31 - 1}
	var uint32s = []uint32{1, 10, 32, 1 << 20, 1<<20 + 1, 1000000000, 1 << 30, 1<<31 - 2, 1<<31 - 1, 1<<32 - 2, 1<<32 - 1}
	var int64s = []int64{1, 10, 32, 1 << 20, 1<<20 + 1, 1000000000, 1 << 30, 1<<31 - 2, 1<<31 - 1, 1000000000000000000, 1 << 60, 1<<63 - 2, 1<<63 - 1}
	var uint64s = []uint64{1, 10, 32, 1 << 20, 1<<20 + 1, 1000000000, 1 << 30, 1<<31 - 2, 1<<31 - 1, 1000000000000000000, 1 << 60, 1<<63 - 2, 1<<63 - 1, 1<<64 - 2, 1<<64 - 1}
	var permSizes = []int{0, 1, 5, 8, 9, 10, 16}
	var readBufferSizes = []int{0, 1, 7, 8, 9, 10}
	var shuffleSliceSizes = []int{0, 1, 7, 8, 9, 10, 239}
	r := rand.New(0)
	
	rv := reflect.ValueOf(r)
	n := rv.NumMethod()
	p := 0
	if *printgolden {
		fmt.Printf("var regressGolden = []interface{}{\n")
	}
	for i := 0; i < n; i++ {
		m := rv.Type().Method(i)
		mv := rv.Method(i)
		mt := mv.Type()
		if m.Name == "Get" || m.Name == "Seed" || m.Name == "UnmarshalBinary" {
			continue
		}
		for repeat := 0; repeat < 17; repeat++ {
			var args []reflect.Value
			var argstr string
			if m.Name == "Shuffle" {
				n := shuffleSliceSizes[repeat%len(shuffleSliceSizes)]
				x := make([]int, n)
				args = append(args, reflect.ValueOf(n))
				args = append(args, reflect.ValueOf(func(i, j int) {
					x[i], x[j] = x[j], x[i]
				}))
			} else if mt.NumIn() == 1 {
				var x interface{}
				switch mt.In(0).Kind() {
				default:
					t.Fatalf("unexpected argument type for r.%s", m.Name)
				
				case reflect.Int:
					if m.Name == "Perm" {
						x = permSizes[repeat%len(permSizes)]
						break
					}
					big := int64s[repeat%len(int64s)]
					if int64(int(big)) != big {
						r.Int63n(big) // what would happen on 64-bit machine, to keep stream in sync
						if *printgolden {
							fmt.Printf("\tskipped, // must run printgolden on 64-bit machine\n")
						}
						p++
						continue
					}
					x = int(big)
				
				case reflect.Int32:
					x = int32s[repeat%len(int32s)]
				
				case reflect.Uint32:
					x = uint32s[repeat%len(uint32s)]
				
				case reflect.Int64:
					x = int64s[repeat%len(int64s)]
				
				case reflect.Uint64:
					x = uint64s[repeat%len(uint64s)]
				
				case reflect.Slice:
					if m.Name == "Read" {
						n := readBufferSizes[repeat%len(readBufferSizes)]
						x = make([]byte, n)
					}
				}
				argstr = fmt.Sprint(x)
				args = append(args, reflect.ValueOf(x))
			}
			
			ret := mv.Call(args)
			if m.Name == "Shuffle" {
				continue // we only run Shuffle for the side effects
			}
			out := ret[0].Interface()
			if m.Name == "Int" || m.Name == "Intn" {
				out = int64(out.(int))
			}
			if m.Name == "Read" {
				out = args[0].Interface().([]byte)
			}
			if *printgolden {
				var val string
				big := int64(1 << 60)
				if int64(int(big)) != big && (m.Name == "Int" || m.Name == "Intn") {
					// 32-bit machine cannot print 64-bit results
					val = "truncated"
				} else if reflect.TypeOf(out).Kind() == reflect.Slice {
					val = fmt.Sprintf("%#v", out)
				} else {
					val = fmt.Sprintf("%T(%v)", out, out)
				}
				fmt.Printf("\t%s, // %s(%s)\n", val, m.Name, argstr)
			} else {
				want := regressGolden[p]
				if m.Name == "Int" {
					want = int64(int(uint(want.(int64)) << 1 >> 1))
				}
				if !reflect.DeepEqual(out, want) {
					t.Errorf("r.%s(%s) = %v, want %v", m.Name, argstr, out, want)
				}
			}
			p++
		}
	}
	if *printgolden {
		fmt.Printf("}\n")
	}
}

var regressGolden = []interface{}{
	float64(0.22067985252185793), // ExpFloat64()
	float64(1.9687711464165194),  // ExpFloat64()
	float64(0.09365679875798526), // ExpFloat64()
	float64(0.14517501157814602), // ExpFloat64()
	float64(0.49508896017758675), // ExpFloat64()
	float64(0.19460162662744554), // ExpFloat64()
	float64(1.772112345348705),   // ExpFloat64()
	float64(0.6731399041877683),  // ExpFloat64()
	float64(0.9608592383348641),  // ExpFloat64()
	float64(1.6377580380236019),  // ExpFloat64()
	float64(0.746790875739628),   // ExpFloat64()
	float64(0.7046262185514),     // ExpFloat64()
	float64(1.2004224748791037),  // ExpFloat64()
	float64(0.2862998393251507),  // ExpFloat64()
	float64(0.06920911706531854), // ExpFloat64()
	float64(0.8560046295086123),  // ExpFloat64()
	float64(1.022440348964754),   // ExpFloat64()
	float32(0.6771215),           // Float32()
	float32(0.27626145),          // Float32()
	float32(0.8183098),           // Float32()
	float32(0.3243996),           // Float32()
	float32(0.67201096),          // Float32()
	float32(0.4681297),           // Float32()
	float32(0.023567796),         // Float32()
	float32(0.087473094),         // Float32()
	float32(0.0034111738),        // Float32()
	float32(0.65722114),          // Float32()
	float32(0.046393096),         // Float32()
	float32(0.21173078),          // Float32()
	float32(0.47271806),          // Float32()
	float32(0.29274207),          // Float32()
	float32(0.27181208),          // Float32()
	float32(0.6496809),           // Float32()
	float32(0.74196166),          // Float32()
	float64(0.856433858351397),   // Float64()
	float64(0.7891435426818407),  // Float64()
	float64(0.2733668469637417),  // Float64()
	float64(0.09475695109948656), // Float64()
	float64(0.9273195412198052),  // Float64()
	float64(0.4249010634878422),  // Float64()
	float64(0.434481617284035),   // Float64()
	float64(0.24533397715360217), // Float64()
	float64(0.22545626444238742), // Float64()
	float64(0.7962420121491581),  // Float64()
	float64(0.9245530787008205),  // Float64()
	float64(0.8394583155312959),  // Float64()
	float64(0.4300312870817893),  // Float64()
	float64(0.2487366685162612),  // Float64()
	float64(0.4381898278658328),  // Float64()
	float64(0.592397672040487),   // Float64()
	float64(0.14746941299436844), // Float64()
	int64(5754373348782608125),   // Int()
	int64(7748491296369333668),   // Int()
	int64(572057954588715219),    // Int()
	int64(6655530453728205615),   // Int()
	int64(7746168941076259749),   // Int()
	int64(2065021622730388476),   // Int()
	int64(7739025699315706832),   // Int()
	int64(1416132004977955628),   // Int()
	int64(2672183821718751310),   // Int()
	int64(1467583583146080573),   // Int()
	int64(6526556134661863112),   // Int()
	int64(1498962930278429112),   // Int()
	int64(3564578358808135765),   // Int()
	int64(7493566175953169584),   // Int()
	int64(2164480193314143082),   // Int()
	int64(8892254210449407921),   // Int()
	int64(752890949371391472),    // Int()
	int32(1205287211),            // Int31()
	int32(404925465),             // Int31()
	int32(1867989579),            // Int31()
	int32(151674396),             // Int31()
	int32(1265122101),            // Int31()
	int32(408483400),             // Int31()
	int32(1543085239),            // Int31()
	int32(1850147509),            // Int31()
	int32(2102981969),            // Int31()
	int32(1217480144),            // Int31()
	int32(2146262991),            // Int31()
	int32(689039740),             // Int31()
	int32(44876493),              // Int31()
	int32(1190852950),            // Int31()
	int32(1593076892),            // Int31()
	int32(1948965381),            // Int31()
	int32(1582074401),            // Int31()
	int32(0),                     // Int31n(1)
	int32(6),                     // Int31n(10)
	int32(29),                    // Int31n(32)
	int32(171754),                // Int31n(1048576)
	int32(662959),                // Int31n(1048577)
	int32(902730596),             // Int31n(1000000000)
	int32(174711228),             // Int31n(1073741824)
	int32(1236167451),            // Int31n(2147483646)
	int32(1417043963),            // Int31n(2147483647)
	int32(0),                     // Int31n(1)
	int32(8),                     // Int31n(10)
	int32(6),                     // Int31n(32)
	int32(207436),                // Int31n(1048576)
	int32(651393),                // Int31n(1048577)
	int32(848592667),             // Int31n(1000000000)
	int32(508814525),             // Int31n(1073741824)
	int32(1139808083),            // Int31n(2147483646)
	int64(4913831498199109714),   // Int63()
	int64(9107756857070956389),   // Int63()
	int64(1227799260184772992),   // Int63()
	int64(2150828967340353585),   // Int63()
	int64(960667031188823006),    // Int63()
	int64(5125145001232459059),   // Int63()
	int64(4341096159660331390),   // Int63()
	int64(7892524944240304887),   // Int63()
	int64(9003988926428784094),   // Int63()
	int64(1290403754045170150),   // Int63()
	int64(7648611523255928381),   // Int63()
	int64(6895932085076097687),   // Int63()
	int64(8430236826169566034),   // Int63()
	int64(6560226495627602614),   // Int63()
	int64(1031322271605560397),   // Int63()
	int64(3236959108230395884),   // Int63()
	int64(4967355935137401225),   // Int63()
	int64(0),                     // Int63n(1)
	int64(3),                     // Int63n(10)
	int64(7),                     // Int63n(32)
	int64(1009739),               // Int63n(1048576)
	int64(848369),                // Int63n(1048577)
	int64(606497288),             // Int63n(1000000000)
	int64(187638578),             // Int63n(1073741824)
	int64(1183902487),            // Int63n(2147483646)
	int64(1200900157),            // Int63n(2147483647)
	int64(61991983276636305),     // Int63n(1000000000000000000)
	int64(692963167483433090),    // Int63n(1152921504606846976)
	int64(3912258686940198097),   // Int63n(9223372036854775806)
	int64(1177200405359738371),   // Int63n(9223372036854775807)
	int64(0),                     // Int63n(1)
	int64(4),                     // Int63n(10)
	int64(29),                    // Int63n(32)
	int64(337390),                // Int63n(1048576)
	int64(0),                     // Intn(1)
	int64(1),                     // Intn(10)
	int64(5),                     // Intn(32)
	int64(720876),                // Intn(1048576)
	int64(126152),                // Intn(1048577)
	int64(782208792),             // Intn(1000000000)
	int64(1053629115),            // Intn(1073741824)
	int64(1724409739),            // Intn(2147483646)
	int64(102204766),             // Intn(2147483647)
	int64(350818036186644838),    // Intn(1000000000000000000)
	int64(895031574546959106),    // Intn(1152921504606846976)
	int64(2272837822344028440),   // Intn(9223372036854775806)
	int64(9015800786283557131),   // Intn(9223372036854775807)
	int64(0),                     // Intn(1)
	int64(4),                     // Intn(10)
	int64(23),                    // Intn(32)
	int64(213701),                // Intn(1048576)
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	[]byte{0x6c, 0x7e, 0x6c, 0xb7, 0x4f, 0x80, 0x7a, 0xcc, 0x32, 0x5c, 0xcb, 0xa1, 0x53, 0x59, 0xd9, 0xca, 0xe0, 0x2f, 0xce, 0xf0, 0xc9, 0x14, 0xb0, 0xcb, 0x9d, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x21, 0x8a, 0x4c, 0x5e, 0x5, 0xda, 0x2a, 0xf4, 0x0}, // MarshalBinary()
	float64(-0.8654257554398836),                                // NormFloat64()
	float64(-0.21406829968820063),                               // NormFloat64()
	float64(-1.259634794338612),                                 // NormFloat64()
	float64(0.9792767971163675),                                 // NormFloat64()
	float64(1.079517806578937),                                  // NormFloat64()
	float64(-1.7279815182679379),                                // NormFloat64()
	float64(-0.1091512345583307),                                // NormFloat64()
	float64(1.8756598905697905),                                 // NormFloat64()
	float64(0.1152268468912775),                                 // NormFloat64()
	float64(0.4380076443898085),                                 // NormFloat64()
	float64(-0.6122218559579252),                                // NormFloat64()
	float64(2.203114764815355),                                  // NormFloat64()
	float64(-1.007500691429182),                                 // NormFloat64()
	float64(-0.009209736102766444),                              // NormFloat64()
	float64(1.8994576881568932),                                 // NormFloat64()
	float64(2.077433728093697),                                  // NormFloat64()
	float64(0.058706583568411005),                               // NormFloat64()
	[]int{},                                                     // Perm(0)
	[]int{0},                                                    // Perm(1)
	[]int{0, 2, 4, 3, 1},                                        // Perm(5)
	[]int{7, 5, 6, 0, 4, 3, 2, 1},                               // Perm(8)
	[]int{8, 6, 2, 4, 7, 3, 1, 5, 0},                            // Perm(9)
	[]int{9, 4, 7, 2, 8, 6, 3, 1, 5, 0},                         // Perm(10)
	[]int{6, 8, 4, 2, 9, 10, 5, 3, 15, 1, 12, 7, 13, 0, 14, 11}, // Perm(16)
	[]int{},                             // Perm(0)
	[]int{0},                            // Perm(1)
	[]int{2, 1, 3, 0, 4},                // Perm(5)
	[]int{6, 1, 3, 7, 0, 2, 5, 4},       // Perm(8)
	[]int{1, 8, 7, 2, 6, 0, 3, 5, 4},    // Perm(9)
	[]int{0, 5, 4, 8, 3, 6, 9, 7, 1, 2}, // Perm(10)
	[]int{13, 2, 10, 6, 3, 7, 5, 8, 9, 4, 11, 14, 12, 1, 15, 0}, // Perm(16)
	[]int{},              // Perm(0)
	[]int{0},             // Perm(1)
	[]int{2, 4, 1, 0, 3}, // Perm(5)
	[]byte{},             // Read([])
	[]byte{0x94},         // Read([0])
	[]byte{0xd6, 0xea, 0x86, 0xf4, 0x43, 0x15, 0x49},                   // Read([0 0 0 0 0 0 0])
	[]byte{0xde, 0x73, 0x2f, 0x87, 0x13, 0x33, 0x41, 0x5f},             // Read([0 0 0 0 0 0 0 0])
	[]byte{0x94, 0xe4, 0x85, 0x89, 0x88, 0x35, 0xb7, 0x46, 0xe8},       // Read([0 0 0 0 0 0 0 0 0])
	[]byte{0xb5, 0x60, 0xaf, 0x5f, 0xe6, 0x80, 0xe6, 0x3e, 0xdc, 0x38}, // Read([0 0 0 0 0 0 0 0 0 0])
	[]byte{},     // Read([])
	[]byte{0x89}, // Read([0])
	[]byte{0xba, 0xeb, 0xcf, 0xc5, 0xc8, 0x14, 0x3c},                  // Read([0 0 0 0 0 0 0])
	[]byte{0x8c, 0xd9, 0x9f, 0xb3, 0x5c, 0x85, 0x1a, 0x2},             // Read([0 0 0 0 0 0 0 0])
	[]byte{0x1a, 0x84, 0x2e, 0x8, 0xea, 0x1b, 0x6, 0x82, 0xbe},        // Read([0 0 0 0 0 0 0 0 0])
	[]byte{0xd9, 0xf4, 0xd9, 0x58, 0x5, 0xca, 0x22, 0x1b, 0x78, 0x8b}, // Read([0 0 0 0 0 0 0 0 0 0])
	[]byte{},     // Read([])
	[]byte{0xf1}, // Read([0])
	[]byte{0x97, 0x5, 0xdb, 0x7f, 0xf2, 0xd7, 0xf3},              // Read([0 0 0 0 0 0 0])
	[]byte{0x45, 0x2f, 0xf4, 0x1d, 0xb0, 0x29, 0x59, 0x1a},       // Read([0 0 0 0 0 0 0 0])
	[]byte{0x1b, 0x49, 0xcc, 0x93, 0x4a, 0x93, 0x38, 0x4a, 0x88}, // Read([0 0 0 0 0 0 0 0 0])
	uint32(443144931),            // Uint32()
	uint32(2838888050),           // Uint32()
	uint32(540933917),            // Uint32()
	uint32(3532980411),           // Uint32()
	uint32(3879394529),           // Uint32()
	uint32(2263983371),           // Uint32()
	uint32(485587527),            // Uint32()
	uint32(157177437),            // Uint32()
	uint32(1210876971),           // Uint32()
	uint32(1236730850),           // Uint32()
	uint32(1093477689),           // Uint32()
	uint32(3169312281),           // Uint32()
	uint32(3320883706),           // Uint32()
	uint32(2221532646),           // Uint32()
	uint32(3765772079),           // Uint32()
	uint32(1102721479),           // Uint32()
	uint32(443264971),            // Uint32()
	uint32(0),                    // Uint32n(1)
	uint32(6),                    // Uint32n(10)
	uint32(5),                    // Uint32n(32)
	uint32(454419),               // Uint32n(1048576)
	uint32(348174),               // Uint32n(1048577)
	uint32(388944719),            // Uint32n(1000000000)
	uint32(522616556),            // Uint32n(1073741824)
	uint32(1333373448),           // Uint32n(2147483646)
	uint32(1895299264),           // Uint32n(2147483647)
	uint32(2669655105),           // Uint32n(4294967294)
	uint32(2815593974),           // Uint32n(4294967295)
	uint32(0),                    // Uint32n(1)
	uint32(4),                    // Uint32n(10)
	uint32(24),                   // Uint32n(32)
	uint32(542010),               // Uint32n(1048576)
	uint32(907389),               // Uint32n(1048577)
	uint32(549564619),            // Uint32n(1000000000)
	uint64(623435815602436215),   // Uint64()
	uint64(7091866858325530325),  // Uint64()
	uint64(15646221088092807745), // Uint64()
	uint64(7017598857963742454),  // Uint64()
	uint64(18438963929968280692), // Uint64()
	uint64(6664292895603936092),  // Uint64()
	uint64(3934775071970460260),  // Uint64()
	uint64(3277824236972575889),  // Uint64()
	uint64(6836477321205388868),  // Uint64()
	uint64(16094187032350467526), // Uint64()
	uint64(16591668613370222261), // Uint64()
	uint64(11145758340702467251), // Uint64()
	uint64(11306661243905047112), // Uint64()
	uint64(3920891166178067046),  // Uint64()
	uint64(18441123780112909729), // Uint64()
	uint64(11443767348496673295), // Uint64()
	uint64(16268865858039102658), // Uint64()
	uint64(0),                    // Uint64n(1)
	uint64(1),                    // Uint64n(10)
	uint64(23),                   // Uint64n(32)
	uint64(936393),               // Uint64n(1048576)
	uint64(965321),               // Uint64n(1048577)
	uint64(921068474),            // Uint64n(1000000000)
	uint64(551904612),            // Uint64n(1073741824)
	uint64(115775440),            // Uint64n(2147483646)
	uint64(818025944),            // Uint64n(2147483647)
	uint64(15198654419150629),    // Uint64n(1000000000000000000)
	uint64(908755076137728455),   // Uint64n(1152921504606846976)
	uint64(8143090435608784732),  // Uint64n(9223372036854775806)
	uint64(263966714504933425),   // Uint64n(9223372036854775807)
	uint64(10916874489150940206), // Uint64n(18446744073709551614)
	uint64(14331617103672661280), // Uint64n(18446744073709551615)
	uint64(0),                    // Uint64n(1)
	uint64(6),                    // Uint64n(10)
}
