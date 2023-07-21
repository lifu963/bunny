package bilibili

import (
	"app/extractors"
	"app/test"
	"testing"
)

func TestBilibili(t *testing.T) {
	tests := []struct {
		name     string
		args     test.Args
		playlist bool
	}{
		{
			name: "normal test",
			args: test.Args{
				URL:   "https://www.bilibili.com/video/av39522458",
				Title: "果然还是旧版的《起风了》比较好听",
			},
			playlist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				data []*extractors.Data
				err  error
			)
			if tt.playlist {
				// for playlist, we don't check the data
				_, err = New().Extract(tt.args.URL, extractors.Options{
					Playlist:     true,
					ThreadNumber: 9,
				})
				test.CheckError(t, err)
			} else {
				data, err = New().Extract(tt.args.URL, extractors.Options{})
				test.CheckError(t, err)
				test.Check(t, tt.args, data[0])
			}
		})
	}
}
