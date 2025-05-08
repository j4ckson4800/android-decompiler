package smali_test

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/j4ckson4800/android-decompiler/decompiler/smali"
	"github.com/j4ckson4800/android-decompiler/decompiler/smali/internal/defs"
	"github.com/stretchr/testify/require"
)

func TestParser_ReadULEB128(t *testing.T) {
	r := require.New(t)
	tests := []struct {
		input    []byte
		want     uint64
		wantSize int
		wantErr  bool
	}{
		{
			input:    []byte{0x00},
			want:     0,
			wantSize: 1,
			wantErr:  false,
		},
		{
			input:    []byte{0x2A},
			want:     42,
			wantSize: 1,
			wantErr:  false,
		},
		{
			input:    []byte{0xE5, 0x8E, 0x26},
			want:     624485,
			wantSize: 3,
			wantErr:  false,
		},
		{
			input:    []byte{0x7F},
			want:     127,
			wantSize: 1,
			wantErr:  false,
		},
		{
			input:    []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02},
			want:     0,
			wantSize: 0,
			wantErr:  true,
		},
		{
			input:    []byte{0x80},
			want:     0,
			wantSize: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		parser := smali.NewParser(bytes.NewReader(tt.input))

		got, err := parser.ReadULEB128()

		if tt.wantErr {
			r.Error(err)
			continue
		}

		r.NoError(err)
		r.Equal(tt.want, got)
	}
}

func TestParser_ReadSLEB128(t *testing.T) {
	r := require.New(t)

	tests := []struct {
		input    []byte
		want     int64
		wantSize int
		wantErr  bool
	}{
		{
			input:    []byte{0x00},
			want:     0,
			wantSize: 1,
			wantErr:  false,
		},
		{
			input:    []byte{0x2A},
			want:     42,
			wantSize: 1,
			wantErr:  false,
		},
		{
			input:    []byte{0x7E},
			want:     -2,
			wantSize: 1,
			wantErr:  false,
		},
		{
			input:    []byte{0xE5, 0x00},
			want:     101,
			wantSize: 2,
			wantErr:  false,
		},
		{
			input:    []byte{0x9B, 0x7F},
			want:     -101,
			wantSize: 2,
			wantErr:  false,
		},
		{
			input:    []byte{0x3F},
			want:     63,
			wantSize: 1,
			wantErr:  false,
		},
		{
			input:    []byte{0x40},
			want:     -64,
			wantSize: 1,
			wantErr:  false,
		},
		{
			input:    []byte{0x80},
			want:     0,
			wantSize: 0,
			wantErr:  true,
		},
		{
			input:    []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02},
			want:     0,
			wantSize: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		parser := smali.NewParser(bytes.NewReader(tt.input))

		got, err := parser.ReadSLEB128()

		if tt.wantErr {
			r.Error(err)
			continue
		}

		r.NoError(err)
		r.Equal(tt.want, got)
	}
}

func TestParser_ReadStruct(t *testing.T) {
	r := require.New(t)

	dexHeader, err := hex.DecodeString(
		strings.Join(
			[]string{
				`6465780A30333500A5712C260F97ACA2`,
				`8F7EF7539794E4F31065394C97FC957F`,
				`00A58600700000007856341200000000`,
				`0000000024A486002FE1000070000000`,
				`812500002C85030003380000301B0400`,
				`8C7C000054BB0600C2FF0000B49F0A00`,
				`DA1E0000C49D1200FC2B700004791600`,
			}, "",
		),
	)
	r.NoError(err)

	hdr := defs.DexHeader{}

	parser := smali.NewParser(bytes.NewReader(dexHeader))
	r.NoError(parser.ReadStruct(&hdr))

	r.Equal(uint64(defs.Magic), hdr.Magic&uint64(defs.Magic))
	r.Equal(uint32(defs.LEConstant), hdr.EndianTag)
	r.Equal(uint32(defs.DexHeaderSize), hdr.HeaderSize)
}
