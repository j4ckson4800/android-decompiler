package resource

import (
	"fmt"
	"unsafe"

	"github.com/j4ckson4800/android-decompiler/decompiler/smali/resource/internal"
)

type ResChunkFlag byte

const (
	FlagSparse   ResChunkFlag = 0x01
	FlagOffset16 ResChunkFlag = 0x02
)

type ResTypeEntryFlag uint16

const (
	FlagComplex ResTypeEntryFlag = 0x0001
	FlagCompact ResTypeEntryFlag = 0x0008
)

type RawResTypeChunk struct {
	ID            byte
	Flags         byte
	Reserved      uint16
	EntryCount    uint32
	EntriesOffset uint32
	ConfigSize    uint32
}

func (r *RawResTypeChunk) IsSparse() bool {
	return r.Flags&byte(FlagSparse) != 0
}

func (r *RawResTypeChunk) IsOffset16() bool {
	return r.Flags&byte(FlagOffset16) != 0
}

type TypeEntryOffset struct {
	ID     int
	Offset int32
}

type RawTypeEntry struct {
	Size  uint16
	Flags uint16
}

func (e *RawTypeEntry) IsComplex() bool {
	return e.Flags&uint16(FlagComplex) != 0
}

func (e *RawTypeEntry) IsCompact() bool {
	return e.Flags&uint16(FlagCompact) != 0
}

type TypeEntry struct {
	rawEntry RawTypeEntry
	DataType int
	Data     uint32
	Key      int
}

type ResTypeChunk struct {
	RawChunk    RawResTypeChunk
	entryOffset int64
	Entries     []TypeEntryOffset
}

func NewResTypeChunk(p internal.Parser, chunkEnd int64) (ResTypeChunk, error) {
	const uint16Size = 2
	entriesOffset := p.Pos() - int64(unsafe.Sizeof(internal.ResChunkHeader{}))
	chunk := ResTypeChunk{}
	if err := p.ReadStruct(&chunk.RawChunk); err != nil {
		return chunk, fmt.Errorf("read header: %w", err)
	}

	configOffset := p.Pos() + int64(chunk.RawChunk.ConfigSize-4)
	if err := p.SetCursorTo(configOffset); err != nil {
		return chunk, fmt.Errorf("set cursor: %w", err)
	}

	chunk.entryOffset = entriesOffset + int64(chunk.RawChunk.EntriesOffset)
	chunk.Entries = make([]TypeEntryOffset, chunk.RawChunk.EntryCount)
	for i := range int(chunk.RawChunk.EntryCount) {
		var offset int32
		idx := i

		if p.Pos() >= chunkEnd-uint16Size || offset >= int32(chunk.RawChunk.EntryCount) {
			// TODO: make custom error type, and decide if this is even needed
			return chunk, fmt.Errorf("invalid offset, probability of obfuscation: %d", offset) //nolint:err113 // I don't care rn
		}

		switch {
		case chunk.RawChunk.IsOffset16():
			off, err := p.ReadUint16()
			if err != nil {
				return chunk, fmt.Errorf("read offset16: %w", err)
			}
			offset = int32(off)
		case chunk.RawChunk.IsSparse():
			index, err := p.ReadUint16()
			if err != nil {
				return chunk, fmt.Errorf("read sparse index: %w", err)
			}
			off, err := p.ReadUint16()
			if err != nil {
				return chunk, fmt.Errorf("read offset16: %w", err)
			}
			offset = int32(off * 4)
			idx = int(index)
		default:
			off, err := p.ReadUint32()
			if err != nil {
				return chunk, fmt.Errorf("read offset32: %w", err)
			}
			offset = int32(off)
		}

		chunk.Entries[idx] = TypeEntryOffset{
			ID:     idx,
			Offset: offset,
		}
	}

	return chunk, nil
}

func (c *ResTypeChunk) GetEntry(p internal.Parser, index int, chunkEnd int64) (TypeEntry, error) {
	entryOffset := &c.Entries[index]
	entry := TypeEntry{
		Key: -1,
	}
	entryStartOffset := int64(entryOffset.Offset) + c.entryOffset
	if err := p.SetCursorTo(entryStartOffset); err != nil {
		return entry, fmt.Errorf("set cursor: %w", err)
	}
	if err := p.ReadStruct(&entry.rawEntry); err != nil {
		return entry, fmt.Errorf("read entry: %w", err)
	}
	if p.Pos() >= chunkEnd {
		return entry, nil
	}

	key := int(entry.rawEntry.Size)
	if !entry.rawEntry.IsCompact() {
		k, err := p.ReadUint32()
		if err != nil {
			return entry, fmt.Errorf("read key: %w", err)
		}
		key = int(k)
	}

	if key == -1 {
		return entry, nil
	}

	switch {
	case entry.rawEntry.IsComplex():
		entry.Key = key
		return entry, nil
	case entry.rawEntry.IsCompact():
		dataType := entry.rawEntry.Flags >> 8
		data, err := p.ReadUint32()
		if err != nil {
			return entry, fmt.Errorf("read data: %w", err)
		}
		entry.DataType = int(dataType)
		entry.Data = data
		entry.Key = key
		return entry, nil
	default:
		if _, err := p.ReadUint16(); err != nil {
			return entry, fmt.Errorf("read uint16: %w", err)
		}
		if _, err := p.ReadByte(); err != nil {
			return entry, fmt.Errorf("read byte: %w", err)
		}
		dataType, err := p.ReadByte()
		if err != nil {
			return entry, fmt.Errorf("read byte: %w", err)
		}
		data, err := p.ReadUint32()
		if err != nil {
			return entry, fmt.Errorf("read uint32: %w", err)
		}
		entry.DataType = int(dataType)
		entry.Data = data
		entry.Key = key
		return entry, nil
	}
}
