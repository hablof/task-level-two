package main

import (
	"errors"
	"testing"

	mysort "github.com/hablof/task-level-two/develop/dev03/sort"
	"github.com/stretchr/testify/assert"
)

func Test_parseFlags(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name         string
		progname     string
		args         []string
		wantConfig   *mysort.SortOptions
		wantFilename string
		wantOutput   bool
		wantErr      error
	}{
		{
			name:         "sort_type_conflict",
			progname:     "sort",
			args:         []string{"-nrubchM", "lines.txt"},
			wantConfig:   nil,
			wantFilename: "",
			wantOutput:   false,
			wantErr: ErrSortTypeConflict{
				type1: mysort.Numeric,
				type2: mysort.ByMonth,
			},
		},
		{
			name:     "valid_simple_config",
			progname: "sort",
			args:     []string{"lines.txt"},
			wantConfig: &mysort.SortOptions{
				SortType:     mysort.Alphabetical,
				ColumnNumber: -1,
				Delim:        " ",
			},
			wantFilename: "lines.txt",
			wantOutput:   false,
			wantErr:      nil,
		},
		{
			name:     "valid_sort_numeric",
			progname: "sort",
			args:     []string{"-n", "lines.txt"},
			wantConfig: &mysort.SortOptions{
				SortType:     mysort.Numeric,
				ColumnNumber: -1,
				Delim:        " ",
			},
			wantFilename: "lines.txt",
			wantOutput:   false,
			wantErr:      nil,
		},
		{
			name:         "unsecified filename",
			progname:     "sort",
			args:         []string{},
			wantConfig:   nil,
			wantFilename: "",
			wantOutput:   false,
			wantErr:      ErrFileUnspecified,
		},
		{
			name:         "unsecified filename with flags",
			progname:     "sort",
			args:         []string{"-nr"},
			wantConfig:   nil,
			wantFilename: "",
			wantOutput:   false,
			wantErr:      ErrFileUnspecified,
		},
		{
			name:         "valid_sort_numeric",
			progname:     "sort",
			args:         []string{"--help"},
			wantConfig:   nil,
			wantFilename: "",
			wantOutput:   true,
			wantErr:      errors.New("pflag: help requested"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConfig, gotFilename, gotOutput, err := parseFlags(tt.progname, tt.args)
			assert.Equal(t, tt.wantConfig, gotConfig)
			assert.Equal(t, tt.wantFilename, gotFilename)
			if tt.wantOutput {
				assert.NotEqual(t, "", gotOutput)
			} else {
				assert.Equal(t, "", gotOutput)
			}
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

// Не работает :(
//
// func Example() {
// 	os.Args = []string{"sort", "testlines.txt"}
// 	main()
// 	// Output:
// 	// As the hours pass
// 	// As the hours pass
// 	// As your voice consoles me
// 	// Before I'm alone
// 	// Before I'm alone
// 	// How it feels to rest
// 	// How it feels to rest
// 	// How pleasant, this feeling
// 	// I missed you, I'm sorry
// 	// I see you, you see me
// 	// I showed you I'm growing
// 	// I will let you know
// 	// I will let you know
// 	// I'm so glad to know
// 	// I'm so glad to know
// 	// I've given what I have
// 	// On your patient lips
// 	// On your patient lips
// 	// That I need to ask
// 	// That I need to ask
// 	// The ashes fall slowly
// 	// The moment you hold me
// 	// To eternal bliss
// 	// To eternal bliss
// }
