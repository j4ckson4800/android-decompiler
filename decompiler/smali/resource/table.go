package resource

import (
	"errors"
	"fmt"
	"io"

	"github.com/j4ckson4800/android-decompiler/decompiler/smali/resource/internal"
)

var (
	ErrInvalidType = errors.New("invalid type")
)

type Table struct {
	Strings       StringPool
	StringsByID   map[uint32]string
	StringsByName map[string]string
}

func NewTable(parser internal.Parser) (Table, error) {
	table := Table{
		StringsByID:   make(map[uint32]string, 128),
		StringsByName: make(map[string]string, 128),
	}

	hdr := internal.ResChunkHeader{}
	if err := parser.ReadStruct(&hdr); err != nil {
		return table, fmt.Errorf("read header: %w", err)
	}

	if hdr.Type != internal.ResTableType {
		return table, ErrInvalidType
	}

	tableHeader := internal.ResTableHeader{}
	if err := parser.ReadStruct(&tableHeader); err != nil {
		return table, fmt.Errorf("read table header: %w", err)
	}

	for {
		chunkOffset := parser.Pos()
		if err := parser.ReadStruct(&hdr); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return table, fmt.Errorf("read header: %w", err)
		}

		switch hdr.Type {
		case internal.ResStringPoolType:
			stringPool, err := NewStringPool(parser)
			if err != nil {
				return table, fmt.Errorf("read string pool: %w", err)
			}
			table.Strings = stringPool
		case internal.ResTablePackageType:
			packageHeader := internal.ResTable{}
			if err := parser.ReadStruct(&packageHeader); err != nil {
				return table, fmt.Errorf("read package header: %w", err)
			}
			if err := table.parsePackageTable(parser, chunkOffset, packageHeader); err != nil {
				return table, fmt.Errorf("parse package table: %w", err)
			}
		default:
		}

		if err := parser.SetCursorTo(chunkOffset + int64(hdr.Size)); err != nil {
			return table, fmt.Errorf("set cursor: %w", err)
		}
	}

	return table, nil
}

func (t *Table) parseResTableType(parser internal.Parser, pkg internal.ResTable, typeStrings, keyStrings StringPool, resTypeChunk ResTypeChunk, chunkEnd int64) error {
	if tstr, err := typeStrings.GetString(parser, uint32(resTypeChunk.RawChunk.ID-1)); err != nil || tstr != "string" {
		return fmt.Errorf("get string: %w", err)
	}

	for i := range resTypeChunk.Entries {
		entry, err := resTypeChunk.GetEntry(parser, i, chunkEnd)
		if err != nil {
			return fmt.Errorf("get entry: %w", err)
		}
		if entry.Key == -1 {
			continue
		}
		str, err := keyStrings.GetString(parser, uint32(entry.Key))
		if err != nil {
			return fmt.Errorf("get string: %w", err)
		}

		if str == "" {
			continue
		}

		entryValue, err := t.Strings.GetString(parser, entry.Data)
		if err != nil {
			return fmt.Errorf("get string: %w", err)
		}

		t.StringsByName[str] = entryValue
		t.StringsByID[pkg.PackageID<<24|uint32(resTypeChunk.RawChunk.ID)<<16|uint32(resTypeChunk.Entries[i].ID)] = entryValue
	}

	return nil
}

func (t *Table) parsePackageTable(parser internal.Parser, chunkOffset int64, pkg internal.ResTable) error {
	var keyStrings StringPool
	var typeStrings StringPool

	if pkg.TypeStringOffset != 0 {
		if err := parser.SetCursorTo(chunkOffset + int64(pkg.TypeStringOffset)); err != nil {
			return fmt.Errorf("set cursor: %w", err)
		}
		hdr := internal.ResChunkHeader{}
		if err := parser.ReadStruct(&hdr); err != nil {
			return fmt.Errorf("read header: %w", err)
		}
		sp, err := NewStringPool(parser)
		if err != nil {
			return fmt.Errorf("new string pool: %w", err)
		}
		if err := parser.SetCursorTo(chunkOffset + int64(pkg.TypeStringOffset) + int64(hdr.Size)); err != nil {
			return fmt.Errorf("set cursor: %w", err)
		}
		typeStrings = sp
	}

	if pkg.KeyStringOffset != 0 {
		if err := parser.SetCursorTo(chunkOffset + int64(pkg.KeyStringOffset)); err != nil {
			return fmt.Errorf("set cursor: %w", err)
		}
		hdr := internal.ResChunkHeader{}
		if err := parser.ReadStruct(&hdr); err != nil {
			return fmt.Errorf("read header: %w", err)
		}
		sp, err := NewStringPool(parser)
		if err != nil {
			return fmt.Errorf("new string pool: %w", err)
		}
		if err := parser.SetCursorTo(chunkOffset + int64(pkg.KeyStringOffset) + int64(hdr.Size)); err != nil {
			return fmt.Errorf("set cursor: %w", err)
		}
		keyStrings = sp
	}

	hdr := internal.ResChunkHeader{}
	for {
		chunkOffset = parser.Pos()
		if err := parser.ReadStruct(&hdr); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return fmt.Errorf("read header: %w", err)
		}

		switch hdr.Type {
		case internal.ResTableTypeSpecType:
			if err := parser.SetCursorTo(chunkOffset + int64(hdr.Size)); err != nil {
				return fmt.Errorf("set cursor: %w", err)
			}
			continue
		case internal.ResTableTypeType:
			chunkEnd := chunkOffset + int64(hdr.Size)
			resTypeChunk, err := NewResTypeChunk(parser, chunkEnd)
			if err != nil {
				return fmt.Errorf("new type chunk: %w", err)
			}

			if err := t.parseResTableType(parser, pkg, typeStrings, keyStrings, resTypeChunk, chunkEnd); err != nil {
				return fmt.Errorf("parse res table type: %w", err)
			}
			if err := parser.SetCursorTo(chunkEnd); err != nil {
				return fmt.Errorf("set cursor: %w", err)
			}
		default:
			if err := parser.SetCursorTo(chunkOffset + int64(hdr.Size)); err != nil {
				return fmt.Errorf("set cursor: %w", err)
			}
		}
	}
}
