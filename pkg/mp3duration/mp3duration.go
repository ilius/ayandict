package mp3duration

import (
	"errors"
	"io"
	"os"
	"time"
)

var (
	versions = []string{"2.5", "x", "2", "1"}
	layers   = []string{"x", "3", "2", "1"}
	bitRates = map[string][]int{
		"V1Lx": {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"V1L1": {0, 32, 64, 96, 128, 160, 192, 224, 256, 288, 320, 352, 384, 416, 448},
		"V1L2": {0, 32, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 384},
		"V1L3": {0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320},
		"V2Lx": {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"V2L1": {0, 32, 48, 56, 64, 80, 96, 112, 128, 144, 160, 176, 192, 224, 256},
		"V2L2": {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},
		"V2L3": {0, 8, 16, 24, 32, 40, 48, 56, 64, 80, 96, 112, 128, 144, 160},
		"VxLx": {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"VxL1": {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"VxL2": {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		"VxL3": {0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	sampleRates = map[string][]int{
		"x":   {0, 0, 0},
		"1":   {44100, 48000, 32000},
		"2":   {22050, 24000, 16000},
		"2.5": {11025, 12000, 8000},
	}
	samples = map[string]map[string]int{
		"x": {
			"x": 0,
			"1": 0,
			"2": 0,
			"3": 0,
		},
		"1": { // MPEGv1,     Layers 1,2,3
			"x": 0,
			"1": 384,
			"2": 1152,
			"3": 1152,
		},
		"2": { // MPEGv2/2.5, Layers 1,2,3
			"x": 0,
			"1": 384,
			"2": 1152,
			"3": 576,
		},
	}
)

type frame struct {
	bitRate    int
	sampleRate int
	frameSize  int
	sample     int
}

func skipID3(buffer []byte) int {
	var id3v2Flags, z0, z1, z2, z3 byte
	var tagSize, footerSize int

	// http://id3.org/d3v2.3.0
	if buffer[0] == 0x49 && buffer[1] == 0x44 && buffer[2] == 0x33 { //'ID3'
		id3v2Flags = buffer[5]
		if (id3v2Flags & 0x10) != 0 {
			footerSize = 10
		} else {
			footerSize = 0
		}

		// ID3 size encoding is crazy (7 bits in each of 4 bytes)
		z0 = buffer[6]
		z1 = buffer[7]
		z2 = buffer[8]
		z3 = buffer[9]
		if ((z0 & 0x80) == 0) && ((z1 & 0x80) == 0) && ((z2 & 0x80) == 0) && ((z3 & 0x80) == 0) {
			tagSize = (((int)(z0 & 0x7f)) * 2097152) +
				(((int)(z1 & 0x7f)) * 16384) +
				(((int)(z2 & 0x7f)) * 128) +
				((int)(z3 & 0x7f))
			return 10 + tagSize + footerSize
		}
	}
	return 0
}

func frameSize(samples int, layer string, bitRate, sampleRate, paddingBit int) int {
	if sampleRate == 0 {
		return 0
	} else if layer == "1" {
		return ((samples * bitRate * 125 / sampleRate) + paddingBit*4)
	} else { // layer 2, 3
		return (((samples * bitRate * 125) / sampleRate) + paddingBit)
	}
}

func parseFrameHeader(header []byte) *frame {
	b1 := header[1]
	b2 := header[2]

	versionBits := (b1 & 0x18) >> 3
	version := versions[versionBits]

	var simpleVersion string
	if version == "2.5" {
		simpleVersion = "2"
	} else {
		simpleVersion = version
	}

	layerBits := (b1 & 0x06) >> 1
	layer := layers[layerBits]

	bitRateKey := "V" + simpleVersion + "L" + layer
	bitRateIndex := (b2 & 0xf0) >> 4

	var bitRate int
	frameBitRates, exists := bitRates[bitRateKey]
	if exists && int(bitRateIndex) < len(frameBitRates) {
		bitRate = frameBitRates[bitRateIndex]
	} else {
		bitRate = 0
	}

	sampleRateIdx := (b2 & 0x0c) >> 2
	var sampleRate int
	frameSampleRates, exists := sampleRates[version]
	if exists && int(sampleRateIdx) < len(frameSampleRates) {
		sampleRate = frameSampleRates[sampleRateIdx]
	} else {
		sampleRate = 0
	}
	sample := samples[simpleVersion][layer]

	paddingBit := (b2 & 0x02) >> 1
	return &frame{
		bitRate:    bitRate,
		sampleRate: sampleRate,
		frameSize:  frameSize(sample, layer, bitRate, sampleRate, int(paddingBit)),
		sample:     sample,
	}
}

// Calculate returns the duration of an mp3 file
func Calculate(filename string) (time.Duration, error) {
	var f *os.File
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	stats, statsErr := f.Stat()
	if statsErr != nil {
		return 0, statsErr
	}
	size := stats.Size()

	buffer := make([]byte, 100)
	var bytesRead int
	bytesRead, err = f.Read(buffer)
	if err != nil {
		return 0, err
	}
	if bytesRead < 100 {
		err = errors.New("Corrupt file")
		return 0, err
	}
	offset := int64(skipID3(buffer))

	buffer = make([]byte, 10)
	duration := 0.0
	for offset < size {
		bytesRead, e := f.ReadAt(buffer, offset)
		if e != nil && e != io.EOF {
			return 0, e
		}
		if bytesRead < 10 {
			return time.Duration(duration*1000.0) * time.Millisecond, nil
		}

		if buffer[0] == 0xff && (buffer[1]&0xe0) == 0xe0 {
			info := parseFrameHeader(buffer)
			if info.frameSize > 0 && info.sample > 0 {
				offset += int64(info.frameSize)
				duration += (float64(info.sample) / float64(info.sampleRate))
			} else {
				offset++ // Corrupt file?
			}
		} else if buffer[0] == 0x54 && buffer[1] == 0x41 && buffer[2] == 0x47 { //'TAG'
			offset += 128 // Skip over id3v1 tag size
		} else {
			offset++ // Corrupt file?
		}
	}

	return time.Duration(duration*1000.0) * time.Millisecond, nil
}
