package api

import (
	"github.com/tomvodi/limepipes/internal/apigen/apimodel"
	"testing"
)

func TestModelValidator_ValidateUpdateTune(t *testing.T) {
	type args struct {
		tuneUpd apimodel.UpdateTune
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{

			name: "tune has no title",
			args: args{
				tuneUpd: apimodel.UpdateTune{
					Title: "",
				},
			},
			wantErr: true,
		},
		{

			name: "tune has title",
			args: args{
				tuneUpd: apimodel.UpdateTune{
					Title: "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewAPIModelValidator(NewGinValidator())
			if err := v.ValidateUpdateTune(tt.args.tuneUpd); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdateTune() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModelValidator_ValidateUpdateSet(t *testing.T) {
	type args struct {
		setUpd apimodel.UpdateSet
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{

			name: "set has no title",
			args: args{
				setUpd: apimodel.UpdateSet{
					Title: "",
				},
			},
			wantErr: true,
		},
		{

			name: "set has title",
			args: args{
				setUpd: apimodel.UpdateSet{
					Title: "test",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewAPIModelValidator(NewGinValidator())
			if err := v.ValidateUpdateSet(tt.args.setUpd); (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdateSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
