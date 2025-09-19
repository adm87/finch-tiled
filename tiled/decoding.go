package tiled

import (
	"fmt"
	"strconv"
	"strings"
)

type DecodingFunc func(data string) ([]uint32, error)

var decodingFunctions = map[TMXEncoding]DecodingFunc{
	TMXEncodingCSV:    parse_csv_layer_data,
	TMXEncodingBase64: parse_base64_layer_data,
}

func parse_csv_layer_data(data string) ([]uint32, error) {
	var tileIndices []uint32
	for _, s := range strings.Split(data, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		tileIndex, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid CSV layer data: %w", err)
		}
		tileIndices = append(tileIndices, uint32(tileIndex))
	}
	return tileIndices, nil
}

func parse_base64_layer_data(data string) ([]uint32, error) {
	// TASK: Implement base64 decoding
	return nil, fmt.Errorf("base64 decoding not implemented")
}

func DecodeData(data string, encoding TMXEncoding) ([]uint32, error) {
	if decodeFunc, ok := decodingFunctions[encoding]; ok {
		return decodeFunc(data)
	}
	panic(fmt.Sprintf("unsupported TMX encoding: %s", encoding))
}

func DecodeTile(tileIndex uint32) (gid uint32, hFlip, vFlip, dFlip, hRot bool) {
	gid = tileIndex & 0x1FFFFFFF
	hFlip = (tileIndex & 0x80000000) != 0
	vFlip = (tileIndex & 0x40000000) != 0
	dFlip = (tileIndex & 0x20000000) != 0
	hRot = (tileIndex & 0x10000000) != 0
	return
}
