package block

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"

	. "github.com/nknorg/nkn/common"
	"github.com/nknorg/nkn/common/serialization"
	. "github.com/nknorg/nkn/pb"
	"github.com/nknorg/nkn/vm/signature"
)

type Header struct {
	BlockHeader
	hash Uint256
}

func (h *Header) Marshal() (dAtA []byte, err error) {
	return h.BlockHeader.Marshal()
}

func (h *Header) Unmarshal(dAtA []byte) error {
	return h.BlockHeader.Unmarshal(dAtA)
}

//Serialize the blockheader
func (h *Header) Serialize(w io.Writer) error {
	return nil
}

func (h *Header) Deserialize(r io.Reader) error {
	return nil
}

//Serialize the blockheader data without program
func (h *Header) SerializeUnsigned(w io.Writer) error {
	serialization.WriteUint32(w, h.UnsignedHeader.Version)
	serialization.WriteVarBytes(w, h.UnsignedHeader.PrevBlockHash)
	serialization.WriteVarBytes(w, h.UnsignedHeader.TransactionsRoot)
	serialization.WriteVarBytes(w, h.UnsignedHeader.StateRoot)
	serialization.WriteUint64(w, uint64(h.UnsignedHeader.Timestamp))
	serialization.WriteUint32(w, h.UnsignedHeader.Height)
	serialization.WriteUint32(w, uint32(h.UnsignedHeader.WinnerType))
	serialization.WriteVarBytes(w, h.UnsignedHeader.Signer)
	serialization.WriteVarBytes(w, h.UnsignedHeader.ChordID)

	return nil
}

func (h *Header) DeserializeUnsigned(r io.Reader) error {
	h.UnsignedHeader.Version, _ = serialization.ReadUint32(r)
	h.UnsignedHeader.PrevBlockHash, _ = serialization.ReadVarBytes(r)
	h.UnsignedHeader.TransactionsRoot, _ = serialization.ReadVarBytes(r)
	h.UnsignedHeader.StateRoot, _ = serialization.ReadVarBytes(r)
	timestamp, _ := serialization.ReadUint64(r)
	h.UnsignedHeader.Timestamp = int64(timestamp)
	h.UnsignedHeader.Height, _ = serialization.ReadUint32(r)
	winnerType, _ := serialization.ReadUint32(r)
	h.UnsignedHeader.WinnerType = WinnerType(winnerType)
	h.UnsignedHeader.Signer, _ = serialization.ReadVarBytes(r)
	h.UnsignedHeader.ChordID, _ = serialization.ReadVarBytes(r)

	return nil
}

func (h *Header) GetProgramHashes() ([]Uint160, error) {
	programHashes := []Uint160{}

	pg := *h.Program
	outputHashes, err := ToCodeHash(pg.Code)
	if err != nil {
		return nil, fmt.Errorf("[Header], GetProgramHashes failed: %v", err)
	}
	programHashes = append(programHashes, outputHashes)
	return programHashes, nil
}

func (h *Header) SetPrograms(programs []*Program) {
	if len(programs) != 1 {
		return
	}

	h.Program = programs[0]
}

func (h *Header) GetPrograms() []*Program {
	return []*Program{h.Program}
}

func (h *Header) Hash() Uint256 {
	d := signature.GetHashData(h)
	temp := sha256.Sum256([]byte(d))
	f := sha256.Sum256(temp[:])
	hash := Uint256(f)
	return hash
}

func (h *Header) GetMessage() []byte {
	return signature.GetHashData(h)
}

func (h *Header) ToArray() []byte {
	dt, _ := h.Marshal()
	return dt
}

func (h *Header) GetInfo() ([]byte, error) {
	type programInfo struct {
		Code      string `json:"code"`
		Parameter string `json:"parameter"`
	}

	type headerInfo struct {
		Version          uint32      `json:"version"`
		PrevBlockHash    string      `json:"prevBlockHash"`
		TransactionsRoot string      `json:"transactionsRoot"`
		StateRoot        string      `json:"stateRoot"`
		Timestamp        int64       `json:"timestamp"`
		Height           uint32      `json:"height"`
		ConsensusData    uint64      `json:"consensusData"`
		NextBookKeeper   string      `json:"nextBookKeeper"`
		WinnerHash       string      `json:"winnerHash"`
		WinnerType       string      `json:"winnerType"`
		Signer           string      `json:"signer"`
		ChordID          string      `json:"chordID"`
		Signature        string      `json:"signature"`
		Program          programInfo `json:"program"`
		Hash             string      `json:"hash"`
	}

	hash := h.Hash()
	info := &headerInfo{
		Version:          h.UnsignedHeader.Version,
		PrevBlockHash:    BytesToHexString(h.UnsignedHeader.PrevBlockHash),
		TransactionsRoot: BytesToHexString(h.UnsignedHeader.TransactionsRoot),
		StateRoot:        BytesToHexString(h.UnsignedHeader.StateRoot),
		Timestamp:        h.UnsignedHeader.Timestamp,
		Height:           h.UnsignedHeader.Height,
		ConsensusData:    h.UnsignedHeader.ConsensusData,
		NextBookKeeper:   BytesToHexString(h.UnsignedHeader.NextBookKeeper),
		WinnerHash:       BytesToHexString(h.UnsignedHeader.WinnerHash),
		WinnerType:       h.UnsignedHeader.WinnerType.String(),
		Signer:           BytesToHexString(h.UnsignedHeader.Signer),
		ChordID:          BytesToHexString(h.UnsignedHeader.ChordID),
		Signature:        BytesToHexString(h.Signature),
		Program: programInfo{
			Code:      BytesToHexString(h.Program.Code),
			Parameter: BytesToHexString(h.Program.Parameter),
		},
		Hash: hash.ToHexString(),
	}

	marshaledInfo, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	return marshaledInfo, nil
}
