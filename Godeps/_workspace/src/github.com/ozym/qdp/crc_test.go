package qdp

import (
	"testing"
)

func TestCRC_Table(t *testing.T) {

	table := []uint32{
		0, 1443300200, 2886600400, 4194895288, 236654280, 1478233504, 2719287320, 4094823280,
		473308560, 1244733176, 2956467008, 3862894632, 304943960, 1143607344, 3189970312, 3894679264,
		946617120, 1852520520, 2489466352, 3261415064, 913782248, 1617966720, 2591634232, 3430821968,
		609887920, 1918707160, 2287214688, 3729990408, 708913272, 2084973328, 2253336232, 3494391232,
		1893234240, 652178728, 3705041040, 2328982520, 2126739592, 683965408, 3536682584, 2227862832,
		1827564496, 988375224, 3235933440, 2531749480, 1660249368, 888301168, 3472581576, 2566676640,
		1219775840, 515067400, 3837414320, 2998749400, 1185891240, 279462080, 3936437624, 3165013520,
		1417826544, 42292120, 4169946656, 2928366920, 1520000568, 211705168, 4137113832, 2693815168,
		3786468480, 3082285032, 1304357456, 465168696, 4021019208, 3115114784, 1134945432, 362997744,
		4253479184, 2877420152, 1367930816, 126874792, 4087218136, 2778397872, 1603533064, 160758368,
		3655128992, 2413548744, 1976750448, 601215512, 3620198760, 2176899584, 2076827576, 768531664,
		3320498736, 2481834328, 1776602336, 1071894408, 3421619448, 2650195856, 1744814632, 838385984,
		2439551680, 3345979816, 1030134800, 1801559928, 2675151880, 3379861344, 863867608, 1702531504,
		2371782480, 3680076856, 558924160, 2002223848, 2202372504, 3577907952, 793481032, 2035059744,
		2835653088, 4278428296, 84584240, 1393402968, 2803871528, 4044926016, 185707000, 1561766544,
		3040001136, 3811950360, 423410336, 1329314248, 3140072120, 3979260368, 388478056, 1092663040,
		2506545768, 3277969664, 963173560, 1869602768, 2608714912, 3447379912, 930337392, 1635045656,
		2303772664, 3747071120, 626966824, 1935262272, 2269890864, 3511470680, 725995488, 2101529736,
		2903171400, 4211991072, 17098648, 1459873008, 2735861632, 4111920360, 253749584, 1494805048,
		2973564120, 3879468976, 489880072, 1261828448, 3207066128, 3911250296, 321516736, 1160705960,
		3854478376, 3015290688, 1236314872, 532130192, 3953500896, 3181552008, 1202431024, 296527704,
		4186485176, 2945430224, 1434892136, 58831872, 4153655152, 2710879256, 1537063328, 228244168,
		3721565960, 2346030176, 1910280664, 668701360, 3553204672, 2244909736, 2143788816, 700488824,
		3252980376, 2548271600, 1844087880, 1005424416, 3489629264, 2583201592, 1676771968, 905347560,
		1960195816, 584136064, 3638046776, 2396992336, 2060269600, 751450952, 3603119856, 2160344472,
		1759521656, 1055336464, 3303943592, 2464755392, 1727735216, 821831384, 3405063008, 2633113608,
		1287261640, 448597664, 3769895704, 3065186416, 1117848320, 346423400, 4004447696, 3098019512,
		1351356504, 109777712, 4236383880, 2860848608, 1586962064, 143662584, 4070119488, 2761825064,
		68042920, 1376338880, 2818590328, 4261889296, 169168480, 1544703240, 2786805936, 4028386264,
		406347064, 1312775760, 3023461352, 3794884736, 371414000, 1076121752, 3123533088, 3962197576,
		1013087112, 1785034976, 2423029080, 3328933424, 846820672, 1686009384, 2658628496, 3362812152,
		542402072, 1985176944, 2354733256, 3663553440, 776956112, 2018012088, 2185326080, 3561385320,
	}

	for a, b := range table {
		if b != crc_table[a] {
			t.Errorf("CRC TABLE[%d]: %u != %u", a, b, crc_table[a])
		}

	}
}

func TestCRC_Test1(t *testing.T) {
	test1 := []byte{56, 2, 4, 0, 1, 0, 0, 0, 4, 0, 0, 0}
	if crc(test1) != 1481917136 {
		t.Errorf("CRC TEST1: %u != %u", crc(test1), 1481917136)
	}
}

func TestCRC_Test2(t *testing.T) {
	test2 := []byte{56, 2, 0, 4, 0, 1, 0, 0, 0, 4, 0, 0}
	if crc(test2) != 2625915688 {
		t.Errorf("CRC TEST2: %u != %u", crc(test2), 2625915688)
	}
}
