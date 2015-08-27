package qdp

const CRC_POLYNOMIAL = 1443300200

var crc_table []uint32

func init() {
	crc_table = make([]uint32, 256)
	for count := 0; count < 256; count++ {
		tdata := ((int32)(count)) << 24
		accum := (int32)(0)
		for bits := 1; bits <= 8; bits++ {
			if (tdata ^ accum) < 0 {
				accum = (accum << 1) ^ CRC_POLYNOMIAL
			} else {
				accum = (accum << 1)
			}
			tdata = tdata << 1
		}
		crc_table[count] = (uint32)(accum)
	}
}

func crc(b []byte) uint32 {

	crc := uint32(0)

	for i := 0; i < len(b); i++ {
		crc = (crc << 8) ^ (uint32)(crc_table[((crc>>24)^(uint32)(b[i]))&255])
	}

	return crc
}
