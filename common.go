package mdt

import (
	. "github.com/Gurvan/melee-data-tools/lib"
)

const (
	Top uint8 = iota
	Trans
	XRot
	YRot
	Hip
	Waist
	LLegA
	LLeg
	LKnee
	LFootA
	LFoot
	RLegA
	RLeg
	RKnee
	RFootA
	RFoot
	WaistB
	Bust
	LClavicle
	LShoulderA
	LShoulder
	LArm
	LHand
	LIndex1
	LIndex2
	LMiddle1
	LMiddle2
	LRing1
	LRing2
	LPinky1
	LPinky2
	LHave
	LThumb1
	LThumb2
	NeckN
	HeadN
	RClavicle
	RShoulderA
	RShoulder
	RArm
	RHand
	RIndex1
	RIndex2
	RMiddle1
	RMiddle2
	RRing1
	RRing2
	RPinky1
	RPinky2
	RHave
	RThumb1
	RThumb2
	Throw
	Extra
)

type BoneLookupTable struct {
	_         [4]byte
	Bones     Ptr[[54]int8]
	BoneCount int32
}

type CommonData struct {
	_               [16]byte
	BoneLookupTable Ptr[[32]Ptr[BoneLookupTable]]
	_               [64]byte
}

type CommonFile = File[CommonData]
