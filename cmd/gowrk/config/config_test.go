package config_test

import (
	"errors"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"radiohead/gowrk/cmd/gowrk/config"
	"radiohead/gowrk/test/mocks"
)

func TestNewFromFlagSet(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		parseErr  error
		want      config.Config
		wantUsage string
		wantErr   bool
	}{
		{
			name:     "when config values cannot be parsed",
			parseErr: errors.New("parsing error"),
			wantErr:  true,
		},
		{
			name:    "when the URL is missing",
			wantErr: true,
		},
		{
			name: "when the flags are well-formed",
			args: []string{"-url", "localhost", "-verbose", "true", "rate", "5", "poolsize", "10", "maxconns", "10"},
			want: config.Config{
				URL: url.URL{
					Path: "localhost",
				},
				Rate:      5,
				PoolSize:  10,
				MaxConns:  10,
				IsVerbose: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			flagSet := mocks.NewMockflagSet(ctrl)
			urlString := tt.want.URL.String()
			poolSize := uint(tt.want.PoolSize)
			maxConns := uint(tt.want.MaxConns)

			flagSet.EXPECT().String("url", "", gomock.Any()).Return(&urlString)
			flagSet.EXPECT().Bool("verbose", false, gomock.Any()).Return(&tt.want.IsVerbose)
			flagSet.EXPECT().Uint(gomock.Eq("rate"), uint(10), gomock.Any()).Return(&tt.want.Rate)
			flagSet.EXPECT().Uint(gomock.Eq("poolsize"), uint(1), gomock.Any()).Return(&poolSize)
			flagSet.EXPECT().Uint(gomock.Eq("maxconns"), uint(1), gomock.Any()).Return(&maxConns)

			flagSet.EXPECT().SetOutput(gomock.Any()).AnyTimes()
			flagSet.EXPECT().Parse(tt.args).Return(tt.parseErr)

			if tt.wantErr {
				flagSet.EXPECT().PrintDefaults()
			}

			cfg, usage, err := config.NewFromFlagSet(flagSet, tt.args)

			assert.Equal(t, tt.want, cfg)
			assert.Equal(t, tt.wantUsage, usage)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
