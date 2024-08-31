package common

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestFilenameFromPath(t *testing.T) {
	g := NewGomegaWithT(t)

	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "only filename with extension",
			args: args{file: "myfile.bww"},
			want: "myfile",
		},
		{
			name: "complete filepath",
			args: args{file: "/home/macloud/myfile.bww"},
			want: "myfile",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			got := FilenameFromPath(tt.args.file)
			g.Expect(got).Should(Equal(tt.want))
		})
	}
}
