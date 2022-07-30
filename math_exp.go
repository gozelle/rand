// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-go file.

package rand

import (
	"math"
)

/*
 * Exponential distribution
 *
 * See "The Ziggurat Method for Generating Random Variables"
 * (Marsaglia & Tsang, 2000)
 * https://www.jstatsoft.org/v05/i08/paper [pdf]
 *
 * Fixed correlation and increased number of distinct results generated,
 * see https://github.com/flyingmutant/rand/issues/3
 */

const (
	re = 7.69711747013104972
)

// ExpFloat64 returns an exponentially distributed float64 in the range
// (0, +math.MaxFloat64] with an exponential distribution whose rate parameter
// (lambda) is 1 and whose mean is 1/lambda (1).
// To produce a distribution with a different rate parameter,
// callers can adjust the output using:
//
//  sample = ExpFloat64() / desiredRateParameter
//
func (r *Rand) ExpFloat64() float64 {
	for {
		v := r.Uint64()
		j := v >> 11
		i := v & 0xFF
		x := float64(j) * we[i]
		if j < ke[i] {
			return x
		}
		if i == 0 {
			return re - math.Log(r.Float64())
		}
		if fe[i]+r.Float64()*(fe[i-1]-fe[i]) < math.Exp(-x) {
			return x
		}
	}
}

var ke = [256]uint64{
	0x1c5214272497c5, 0x0, 0x137d5bd79c3125, 0x186ef58e3f3bf1,
	0x1a9bb7320eb09b, 0x1bd127f7194472, 0x1c951d0f886513, 0x1d1bfe2d5c3970,
	0x1d7e5bd56b18b2, 0x1dc934dd172c6e, 0x1e0409dfac9dc8, 0x1e337b71d47835,
	0x1e5a8b177cb7a0, 0x1e7b42096f046d, 0x1e970daf08ae3c, 0x1eaef5b14ef09e,
	0x1ec3bd07b46557, 0x1ed5f6f08799cd, 0x1ee614ae6e5689, 0x1ef46eca361ccf,
	0x1f014b76ddd4a3, 0x1f0ce313a796b6, 0x1f176369f1f779, 0x1f20f20c452571,
	0x1f29ae1951a875, 0x1f31b18fb95533, 0x1f39125157c107, 0x1f3fe2eb6e694c,
	0x1f463332d788fb, 0x1f4c10bf1d3a11, 0x1f51874c5c3323, 0x1f56a109c3ecc1,
	0x1f5b66d9099995, 0x1f5fe08210d08c, 0x1f6414dd445772, 0x1f6809f685967a,
	0x1f6bc52a2b02e7, 0x1f6f4b3d32e4f4, 0x1f72a07190f139, 0x1f75c8974d09d9,
	0x1f78c71b045cc0, 0x1f7b9f12413ff5, 0x1f7e5346079f8a, 0x1f80e63be21138,
	0x1f835a3dad9163, 0x1f85b16056b913, 0x1f87ed89b24263, 0x1f8a10759374fc,
	0x1f8c1bba3d39ad, 0x1f8e10cc45d04b, 0x1f8ff102013e16, 0x1f91bd968358e2,
	0x1f9377ac47afd7, 0x1f95204f8b64db, 0x1f96b878633894, 0x1f98410c968891,
	0x1f99bae146ba83, 0x1f9b26bc697f01, 0x1f9c85561b717a, 0x1f9dd759cfd804,
	0x1f9f1d6761a1cf, 0x1fa058140936c0, 0x1fa187eb3a333a, 0x1fa2ad6f6bc4fc,
	0x1fa3c91ace0683, 0x1fa4db5fee6aa3, 0x1fa5e4aa4d0980, 0x1fa6e55ee46783,
	0x1fa7dddca51ec5, 0x1fa8ce7ce6a876, 0x1fa9b793ce5fef, 0x1faa9970adb858,
	0x1fab745e588233, 0x1fac48a3740585, 0x1fad1682bf9feb, 0x1fadde3b5782c0,
	0x1faea008f21d6c, 0x1faf5c2418b07f, 0x1fb012c25b7a15, 0x1fb0c41681dff3,
	0x1fb17050b6f1fa, 0x1fb2179eb2963b, 0x1fb2ba2bdfa84b, 0x1fb358217f4e19,
	0x1fb3f1a6c9be0d, 0x1fb486e10cacd7, 0x1fb517f3c793fc, 0x1fb5a500c5fdaa,
	0x1fb62e2837fe5a, 0x1fb6b388c9010b, 0x1fb7353fb50798, 0x1fb7b368dc7da9,
	0x1fb82e1ed6ba0a, 0x1fb8a57b0347f6, 0x1fb919959a0f74, 0x1fb98a85ba7204,
	0x1fb9f861796f26, 0x1fba633deee287, 0x1fbacb2f41ec17, 0x1fbb3048b49146,
	0x1fbb929caea4e4, 0x1fbbf23cc8029d, 0x1fbc4f39d22996, 0x1fbca9a3e140d5,
	0x1fbd018a548fa0, 0x1fbd56fbde729c, 0x1fbdaa068bd66c, 0x1fbdfab7cb3f42,
	0x1fbe491c7364de, 0x1fbe9540c96960, 0x1fbedf3086b129, 0x1fbf26f6de6176,
	0x1fbf6c9e828ae4, 0x1fbfb031a904c6, 0x1fbff1ba0ffdb1, 0x1fc03141024589,
	0x1fc06ecf5b54b3, 0x1fc0aa6d8b1428, 0x1fc0e42399698b, 0x1fc11bf9298a65,
	0x1fc151f57d1944, 0x1fc1861f770f4b, 0x1fc1b87d9e74b4, 0x1fc1e91620ea43,
	0x1fc217eed505de, 0x1fc2450d3c8400, 0x1fc27076864fc2, 0x1fc29a2f906310,
	0x1fc2c23ce98045, 0x1fc2e8a2d2c6b5, 0x1fc30d654122ef, 0x1fc33087de9c0f,
	0x1fc3520e0b7ec9, 0x1fc371fadf66f8, 0x1fc390512a2887, 0x1fc3ad137497fa,
	0x1fc3c844013349, 0x1fc3e1e4ccab40, 0x1fc3f9f78e4da8, 0x1fc4107db85061,
	0x1fc4257877fd69, 0x1fc438e8b5bfc7, 0x1fc44acf15112a, 0x1fc45b2bf447e9,
	0x1fc469ff6c4505, 0x1fc477495001b3, 0x1fc483092bfbba, 0x1fc48d3e457ff7,
	0x1fc495e799d21b, 0x1fc49d03dd30b2, 0x1fc4a29179b434, 0x1fc4a68e8e07fc,
	0x1fc4a8f8ebfb8d, 0x1fc4a9ce16eaa0, 0x1fc4a90b41fa35, 0x1fc4a6ad4e28a1,
	0x1fc4a2b0c82e76, 0x1fc49d11e62de3, 0x1fc495cc852df4, 0x1fc48cdc265ec1,
	0x1fc4823bec237a, 0x1fc475e696dee7, 0x1fc467d6817e83, 0x1fc458059dc038,
	0x1fc4466d702e22, 0x1fc433070bcb9a, 0x1fc41dcb0d6e0e, 0x1fc406b196bbf7,
	0x1fc3edb248cb62, 0x1fc3d2c43e593d, 0x1fc3b5de0591b5, 0x1fc396f599614b,
	0x1fc376005a4592, 0x1fc352f3069372, 0x1fc32dc1b2281b, 0x1fc3065fbd7888,
	0x1fc2dcbfcbf264, 0x1fc2b0d3b99f9e, 0x1fc2828c8ffcf0, 0x1fc251da79f164,
	0x1fc21eacb6d39e, 0x1fc1e8f18c6757, 0x1fc1b09637bb3d, 0x1fc17586dccd11,
	0x1fc137ae74d6b8, 0x1fc0f6f6bb2416, 0x1fc0b348184da4, 0x1fc06c898baff1,
	0x1fc022a092f365, 0x1fbfd5710f72b8, 0x1fbf84dd294890, 0x1fbf30c52fc60d,
	0x1fbed907770cc6, 0x1fbe7d80327ddc, 0x1fbe1e094ba615, 0x1fbdba7a354408,
	0x1fbd52a7b9f826, 0x1fbce663c6201b, 0x1fbc757d2c4de5, 0x1fbbffbf63b7aa,
	0x1fbb84f23fe6a2, 0x1fbb04d9a0d18e, 0x1fba7f351a70ad, 0x1fb9f3bf92b61a,
	0x1fb9622ed4abfc, 0x1fb8ca33174a18, 0x1fb82b76765b54, 0x1fb7859c5b895d,
	0x1fb6d840d55594, 0x1fb622f7d96943, 0x1fb5654c6f37e2, 0x1fb49ebfbf69d3,
	0x1fb3cec803e747, 0x1fb2f4cf539c40, 0x1fb21032442854, 0x1fb1203e5a9605,
	0x1fb0243042e1c3, 0x1faf1b31c479a7, 0x1fae045767e106, 0x1facde9dbf2d73,
	0x1faba8e640060b, 0x1faa61f399ff29, 0x1fa908656f66a2, 0x1fa79ab3508d3d,
	0x1fa61726d1f215, 0x1fa47bd48bea00, 0x1fa2c693c5c095, 0x1fa0f4f47df316,
	0x1f9f04336bbe0b, 0x1f9cf12b79f9bd, 0x1f9ab84415abc5, 0x1f98555b782fb9,
	0x1f95c3abd03f7a, 0x1f92fda9cef1f3, 0x1f8ffcda9ae41d, 0x1f8cb99e7385f8,
	0x1f892aec479608, 0x1f8545f904db90, 0x1f80fdc336039b, 0x1f7c427839e926,
	0x1f7700a3582ace, 0x1f71200f1a241d, 0x1f6a8234b7352c, 0x1f630000a8e267,
	0x1f5a66904fe3c6, 0x1f50724ece1173, 0x1f44c7665c6fdb, 0x1f36e5a38a59a4,
	0x1f26143450340b, 0x1f113e047b0414, 0x1ef6aefa57cbe7, 0x1ed38ca188151e,
	0x1ea2a61e122db1, 0x1e5961c78b267d, 0x1dddf62bac0bb1, 0x1cdb4dd9e4e8c0,
}
var we = [256]float64{
	9.655740063209187e-16, 7.089014243955202e-18, 1.1639412496691068e-17,
	1.5243915123532025e-17, 1.8332848857237325e-17, 2.108965109464476e-17,
	2.361128077843129e-17, 2.595595772310885e-17, 2.816173554197743e-17,
	3.025504130321374e-17, 3.2255082548363667e-17, 3.417632340185019e-17,
	3.602996978734446e-17, 3.7824907768696417e-17, 3.9568321980975465e-17,
	4.1266117781759396e-17, 4.292321808442518e-17, 4.4543777432823646e-17,
	4.613133981483179e-17, 4.768895725264629e-17, 4.9219280437279567e-17,
	5.072462904503141e-17, 5.220704702792667e-17, 5.366834661718187e-17,
	5.511014372835089e-17, 5.653388673239661e-17, 5.79408800485276e-17,
	5.933230365208937e-17, 6.070922932847173e-17, 6.207263431163186e-17,
	6.342341280303069e-17, 6.476238575956133e-17, 6.609030925769398e-17,
	6.740788167872715e-17, 6.871574991183805e-17, 7.001451473403922e-17,
	7.130473549660636e-17, 7.258693422414641e-17, 7.386159921381785e-17,
	7.51291882072372e-17, 7.639013119550817e-17, 7.764483290797841e-17,
	7.889367502729783e-17, 8.013701816675448e-17, 8.137520364041755e-17,
	8.260855505210031e-17, 8.383737972539132e-17, 8.506196999385315e-17,
	8.628260436784104e-17, 8.749954859216174e-17, 8.871305660690245e-17,
	8.992337142215348e-17, 9.113072591597902e-17, 9.233534356381781e-17,
	9.35374391064912e-17, 9.473721916312942e-17, 9.593488279457989e-17,
	9.713062202221511e-17, 9.832462230649502e-17, 9.951706298915062e-17,
	1.007081177024294e-16, 1.0189795474846932e-16, 1.0308673745154211e-16,
	1.0427462448561878e-16, 1.0546177017945757e-16, 1.0664832480119141e-16,
	1.0783443482419476e-16, 1.0902024317583496e-16, 1.1020588947055772e-16,
	1.1139151022861965e-16, 1.1257723908165665e-16, 1.1376320696616837e-16,
	1.1494954230590083e-16, 1.1613637118402173e-16, 1.1732381750590448e-16,
	1.1851200315326687e-16, 1.1970104813034644e-16, 1.2089107070273848e-16,
	1.2208218752947052e-16, 1.2327451378884145e-16, 1.244681632985112e-16,
	1.256632486302898e-16, 1.2685988122003973e-16, 1.2805817147307491e-16,
	1.292582288654119e-16, 1.3046016204120286e-16, 1.3166407890665723e-16,
	1.328700867207381e-16, 1.3407829218289992e-16, 1.3528880151811752e-16,
	1.3650172055943978e-16, 1.377171548282881e-16, 1.3893520961270637e-16,
	1.4015599004375713e-16, 1.413796011702485e-16, 1.4260614803196652e-16,
	1.4383573573157902e-16, 1.4506846950536877e-16, 1.4630445479294757e-16,
	1.4754379730609514e-16, 1.4878660309686256e-16, 1.5003297862507367e-16,
	1.5128303082535392e-16, 1.5253686717381255e-16, 1.5379459575449967e-16,
	1.5505632532575771e-16, 1.5632216538658375e-16, 1.5759222624311761e-16,
	1.5886661907536842e-16, 1.6014545600429167e-16, 1.6142885015932787e-16,
	1.6271691574651303e-16, 1.6400976811727177e-16, 1.6530752383800364e-16,
	1.6661030076057416e-16, 1.6791821809382284e-16, 1.692313964762022e-16,
	1.7054995804966296e-16, 1.7187402653490314e-16, 1.732037273081008e-16,
	1.7453918747925335e-16, 1.758805359722491e-16, 1.772279036068006e-16,
	1.7858142318237321e-16, 1.7994122956424635e-16, 1.8130745977185013e-16,
	1.826802530695252e-16, 1.8405975105985876e-16, 1.8544609777975695e-16,
	1.8683943979941927e-16, 1.8823992632438918e-16, 1.8964770930086165e-16,
	1.910629435244376e-16, 1.9248578675252436e-16, 1.9391639982058992e-16,
	1.953549467624909e-16, 1.9680159493510374e-16, 1.982565151475019e-16,
	1.997198817949342e-16, 2.0119187299787347e-16, 2.0267267074641983e-16,
	2.0416246105035888e-16, 2.0566143409519179e-16, 2.071697844044737e-16,
	2.0868771100881597e-16, 2.1021541762192925e-16, 2.1175311282410757e-16,
	2.1330101025357788e-16, 2.148593288061663e-16, 2.1642829284376045e-16,
	2.1800813241207835e-16, 2.1959908346828702e-16, 2.2120138811904954e-16,
	2.22815294869618e-16, 2.2444105888463076e-16, 2.2607894226131728e-16,
	2.27729214315862e-16, 2.2939215188373104e-16, 2.310680396348213e-16,
	2.327571704043534e-16, 2.3445984554049574e-16, 2.3617637526977735e-16,
	2.379070790814276e-16, 2.396522861318623e-16, 2.4141233567062923e-16,
	2.431875774892255e-16, 2.4497837239430697e-16, 2.467850927069288e-16,
	2.486081227895851e-16, 2.504478596029556e-16, 2.523047132944216e-16,
	2.5417910782058117e-16, 2.560714816061771e-16, 2.579822882420531e-16,
	2.5991199722497464e-16, 2.618610947423924e-16, 2.6383008450549423e-16,
	2.658194886341845e-16, 2.6782984859795257e-16, 2.6986172621694894e-16,
	2.719157047279819e-16, 2.7399238992058153e-16, 2.760924113487617e-16,
	2.782164236246436e-16, 2.8036510780069835e-16, 2.825391728480253e-16,
	2.847393572388174e-16, 2.8696643064198177e-16, 2.8922119574179956e-16,
	2.9150449019052937e-16, 2.9381718870700286e-16, 2.9616020533454657e-16,
	2.9853449587300453e-16, 3.009410605012618e-16, 3.033809466085003e-16,
	3.0585525185448604e-16, 3.08365127481531e-16, 3.1091178190342663e-16,
	3.1349648459966636e-16, 3.161205703467106e-16, 3.1878544382197136e-16,
	3.214925846206798e-16, 3.242435527309452e-16, 3.270399945182241e-16,
	3.2988364927722836e-16, 3.327763564171672e-16, 3.3572006335532446e-16,
	3.387168342045505e-16, 3.417688593525637e-16, 3.4487846604534244e-16,
	3.4804813010374423e-16, 3.5128048892229794e-16, 3.545783559224792e-16,
	3.5794473666042765e-16, 3.6138284682190606e-16, 3.6489613237645425e-16,
	3.684882922095621e-16, 3.7216330360802073e-16, 3.759254510416256e-16,
	3.7977935876688744e-16, 3.8373002787892137e-16, 3.8778287856078953e-16,
	3.919437984311429e-16, 3.962191980786775e-16, 4.0061607510565417e-16,
	4.051420882956573e-16, 4.0980564389030625e-16, 4.1461599642909046e-16,
	4.195833672073399e-16, 4.247190841824385e-16, 4.3003574816674707e-16,
	4.355474314693952e-16, 4.4126991690360704e-16, 4.472209874259932e-16,
	4.534207798565834e-16, 4.598922204905932e-16, 4.666615664711476e-16,
	4.737590853262492e-16, 4.812199172829238e-16, 4.89085182739221e-16,
	4.97403423619194e-16, 5.06232507214416e-16, 5.156421828878083e-16,
	5.257175802022275e-16, 5.365640977112021e-16, 5.483144034258703e-16,
	5.611387454675159e-16, 5.752606481503331e-16, 5.909817641652101e-16,
	6.087231416180907e-16, 6.290979034877556e-16, 6.53049205356404e-16,
	6.821393079028929e-16, 7.192444966089362e-16, 7.706095350032097e-16,
	8.545517038584027e-16,
}
var fe = [256]float64{
	1, 0.9381436808621765, 0.9004699299257477, 0.8717043323812047,
	0.8477855006239905, 0.8269932966430511, 0.808421651523009,
	0.7915276369724963, 0.7759568520401162, 0.7614633888498968,
	0.7478686219851957, 0.735038092431424, 0.7228676595935725,
	0.7112747608050765, 0.7001926550827886, 0.6895664961170784,
	0.6793505722647658, 0.6695063167319252, 0.6600008410790001,
	0.6508058334145714, 0.6418967164272664, 0.6332519942143664,
	0.6248527387036662, 0.6166821809152079, 0.6087253820796223,
	0.6009689663652326, 0.5934009016917338, 0.5860103184772684,
	0.5787873586028454, 0.5717230486648262, 0.5648091929124006,
	0.5580382822625879, 0.5514034165406417, 0.5448982376724401,
	0.5385168720028622, 0.5322538802630437, 0.5261042139836201,
	0.5200631773682339, 0.5141263938147489, 0.5082897764106432,
	0.5025495018413481, 0.4969019872415499, 0.49134386959403287,
	0.48587198734188525, 0.48048336393045454, 0.4751751930373777,
	0.4699448252839603, 0.4647897562504265, 0.459707615642138,
	0.45469615747461584, 0.44975325116275533, 0.44487687341454885,
	0.4400651008423542, 0.4353161032156369, 0.43062813728845917,
	0.4259995411430347, 0.4214287289976169, 0.41691418643300326,
	0.4124544659971615, 0.40804818315203273, 0.4036940125305306,
	0.3993906844752314, 0.39513698183329043, 0.3909317369847974,
	0.38677382908413793, 0.38266218149601006, 0.37859575940958107,
	0.37457356761590244, 0.3705946484351463, 0.36665807978151443,
	0.3627629733548181, 0.35890847294875006, 0.35509375286678774,
	0.3513180164374836, 0.34758049462163726, 0.3438804447045027,
	0.34021714906678024, 0.3365899140286778, 0.3329980687618092,
	0.32944096426413655, 0.32591797239355635, 0.3224284849560893,
	0.31897191284495735, 0.315547685227129, 0.3121552487741797,
	0.30879406693456024, 0.3054636192445903, 0.3021634006756935,
	0.2988929210155818, 0.29565170428126125, 0.29243928816189263,
	0.28925522348967775, 0.2860990737370769, 0.2829704145387808,
	0.2798688332369729, 0.2767939284485174, 0.27374530965280297,
	0.27072259679906, 0.26772541993204485, 0.26475341883506226,
	0.26180624268936303, 0.2588835497490163, 0.25598500703041543,
	0.2531102900156295, 0.25025908236886235, 0.24743107566532765,
	0.24462596913189213, 0.24184346939887724, 0.23908329026244915,
	0.23634515245705962, 0.23362878343743335, 0.23093391716962744,
	0.22826029393071676, 0.22560766011668415, 0.22297576805812028,
	0.22036437584335958, 0.21777324714870058, 0.2152021510753787,
	0.21265086199297834, 0.2101191593889883, 0.20760682772422212,
	0.20511365629383782, 0.2026394390937091, 0.20018397469191135,
	0.19774706610509893, 0.19532852067956327, 0.19292814997677138,
	0.19054576966319545, 0.18818119940425432, 0.18583426276219714,
	0.1835047870977675, 0.18119260347549634, 0.17889754657247836,
	0.17661945459049494, 0.17435816917135352, 0.1721135353153201,
	0.16988540130252766, 0.1676736186172502, 0.165478041874936,
	0.16329852875190182, 0.16113493991759203, 0.1589871389693142,
	0.15685499236936523, 0.15473836938446808, 0.15263714202744288,
	0.15055118500103992, 0.14848037564386682, 0.14642459387834497,
	0.1443837221606348, 0.14235764543247223, 0.1403462510748625,
	0.1383494288635803, 0.13636707092642894, 0.1343990717022137,
	0.13244532790138763, 0.13050573846833088, 0.1285802045452283,
	0.12666862943751078, 0.12477091858083104, 0.12288697950954522,
	0.1210167218266749, 0.11916005717532775, 0.11731689921155564,
	0.1154871635786336, 0.11367076788274438, 0.11186763167005638,
	0.11007767640518545, 0.10830082545103385, 0.10653700405000172,
	0.10478613930657024, 0.1030481601712578, 0.10132299742595369,
	0.09961058367063715, 0.09791085331149221, 0.09622374255043283,
	0.09454918937605587, 0.09288713355604357, 0.09123751663104017,
	0.08960028191003284, 0.08797537446727019, 0.08636274114075689,
	0.0847623305323681, 0.08317409300963235, 0.08159798070923742,
	0.0800339475423199, 0.07848194920160644, 0.07694194317048052,
	0.0754138887340584, 0.07389774699236475, 0.07239348087570872,
	0.07090105516237181, 0.06942043649872875, 0.06795159342193662,
	0.06649449638533979, 0.06504911778675376, 0.06361543199980735,
	0.06219341540854101, 0.06078304644547963, 0.05938430563342025,
	0.05799717563120064, 0.05662164128374284, 0.05525768967669701,
	0.05390531019604605, 0.052564494593071664, 0.051235237055126254,
	0.04991753428270636, 0.04861138557337948, 0.04731679291318155,
	0.04603376107617516, 0.04476229773294327, 0.043502413568888176,
	0.04225412241331624, 0.04101744138041482, 0.03979239102337412,
	0.03857899550307485, 0.03737728277295936, 0.03618728478193143,
	0.03500903769739742, 0.033842582150874344, 0.03268796350895954,
	0.03154523217289361, 0.030414443910466608, 0.029295660224637397,
	0.028188948763978632, 0.0270943837809558, 0.02601204664513422,
	0.024942026419731787, 0.023884420511558174, 0.02283933540638524,
	0.021806887504283584, 0.020787204072578117, 0.01978042433800974,
	0.018786700744696024, 0.017806200410911355, 0.01683910682603994,
	0.015885621839973156, 0.014945968011691148, 0.014020391403181943,
	0.013109164931254991, 0.012212592426255378, 0.0113310135978346,
	0.010464810181029982, 0.009614413642502213, 0.008780314985808977,
	0.007963077438017043, 0.007163353183634991, 0.006381905937319183,
	0.005619642207205489, 0.0048776559835424, 0.0041572951208338005,
	0.003460264777836907, 0.0027887987935740783, 0.002145967743718907,
	0.0015362997803015728, 0.0009672692823271743, 0.0004541343538414966,
}
