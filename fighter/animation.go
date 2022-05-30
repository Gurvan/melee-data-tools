package fighter

import (
	"io"

	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/common"
	"github.com/Gurvan/melee-data-tools/descriptor"
)

type TrackType uint8

const (
	TRACK_NONE TrackType = iota
	TRACK_ROTX
	TRACK_ROTY
	TRACK_ROTZ
	TRACK_PATH
	TRACK_TRAX
	TRACK_TRAY
	TRACK_TRAZ
	TRACK_SCAX
	TRACK_SCAY
	TRACK_SCAZ
	TRACK_NODE
)

type InterpolationType uint8

const (
	INTERPOLATION_NONE InterpolationType = iota
	INTERPOLATION_CON
	INTERPOLATION_LIN
	INTERPOLATION_SPL0
	INTERPOLATION_SPL
	INTERPOLATION_SLP
	INTERPOLATION_KEY
)

type AnimationDataFormat uint8

const (
	FORMAT_FLOAT AnimationDataFormat = 0x0
	FORMAT_S16   AnimationDataFormat = 0x20
	FORMAT_U16   AnimationDataFormat = 0x40
	FORMAT_S8    AnimationDataFormat = 0x60
	FORMAT_U8    AnimationDataFormat = 0x80
)

type Animation struct {
	Type        uint32
	_           [4]byte
	FrameCount  float32
	TracksCount Ptr[uint32]
	Tracks      Ptr[[][]Track]
}

type Track struct {
	DataLength uint16
	StartFrame uint16
	Type       TrackType
	// ValueFlags uint8
	ValueFormat AnimationDataFormat
	ValueScale  uint32
	// TangentFlags uint8
	TangentFormat AnimationDataFormat
	TangetnScale  uint32
	// Data Ptr[[]byte]
	Keys []TrackKey
}

type TrackKey struct {
	Frame         float32
	Value         float32
	Tangent       float32
	Interpolation InterpolationType
}

func padTo0x20(n int64) int64 {
	if n%0x20 == 0 {
		return n
	}
	return 0x20 * (n/0x20 + 1)
}

func SplitAnimationFile(r binread.Reader) ([][]byte, error) {
	filesOffsets := []int64{0}
	var fs int64 = 0
	h := descriptor.Header{}

	for {
		err := r.Decode(&h)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		fs += padTo0x20(int64(h.FileSize))
		filesOffsets = append(filesOffsets, fs)

		_, err = r.Seek(fs, io.SeekStart)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	filesData := make([][]byte, 0)
	r.Seek(0, io.SeekStart)
	var offset int64 = 0
	for _, i := range filesOffsets[1:] {
		// fmt.Println(i, offset, i-offset)
		b := make([]byte, i-offset)
		io.ReadFull(r, b)
		filesData = append(filesData, b)
		offset = i
	}

	return filesData, nil
}
