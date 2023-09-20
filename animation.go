package mdt

import (
	"io"

	"github.com/Gurvan/melee-data-tools/animation"
	"github.com/Gurvan/melee-data-tools/binread"
	. "github.com/Gurvan/melee-data-tools/lib"
	"github.com/Gurvan/melee-data-tools/descriptor"
)

type TracksCounts []uint8

func (t *TracksCounts) BinRead(r *binread.Reader, args ...Args) error {

	// Follow Ptr
	var addr Addr
	err := r.Decode(&addr)
	if err != nil {
		return err
	}

	if addr == 0x20 {
		*t = make([]uint8, 0)
		return nil
	}

	before := r.CurrentPosition()

	_, err = r.Seek(addr.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	var counts []uint8 = make([]uint8, 0)
	var c uint8
	for {
		err = r.Decode(&c)
		if err != nil {
			return err
		}
		if c == 0xFF {
			break
		}
		counts = append(counts, c)
	}

	*t = counts
	_, err = r.Seek(before, io.SeekStart)
	return err
}

type AnimationData struct {
	Type         uint32
	Unused       [4]byte
	FrameCount   float32
	TracksCounts TracksCounts
	Tracks       Ptr[animation.Tracks]
}

func (a *AnimationData) AfterParse(r *binread.Reader, _ ...Args) error {
	before := r.CurrentPosition()
	_, err := r.Seek(a.Tracks.Offset.ToSeek(), io.SeekStart)
	if err != nil {
		return err
	}

	nodes := make([][]animation.Track, 0)
	for _, c := range a.TracksCounts {
		// if c == 0 {
		//         continue
		// }
		tracks := make([]animation.Track, c)
		err = r.Decode(&tracks)
		if err != nil {
			return err
		}
		nodes = append(nodes, tracks)
	}

	a.Tracks.SetValue(nodes)
	_, err = r.Seek(before, io.SeekStart)
	return err
}

type AnimationFile = File[AnimationData]

func padTo0x20(n int64) int64 {
	if n%0x20 == 0 {
		return n
	}
	return 0x20 * (n/0x20 + 1)
}

func SplitAnimationFile(r binread.Reader) ([][]byte, []int64, error) {
	filesOffsets := []int64{0}
	var fs int64 = 0
	h := descriptor.Header{}

	for {
		err := r.Decode(&h)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		fs += padTo0x20(int64(h.FileSize))
		filesOffsets = append(filesOffsets, fs)

		_, err = r.Seek(fs, io.SeekStart)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
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

	return filesData, filesOffsets, nil
}
