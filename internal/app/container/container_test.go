package container

import (
	"HappyKod/service-api-gofermart/internal/models"
	"testing"

	"go.uber.org/zap"
)

func TestBuildContainer(t *testing.T) {
	type args struct {
		cfg    models.Config
		logger *zap.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "сборка контейнера",
			args: args{cfg: models.Config{}, logger: &zap.Logger{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BuildContainer(tt.args.cfg, tt.args.logger); (err != nil) != tt.wantErr {
				t.Errorf("BuildContainer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
