package utils

import (
	"strings"
	"testing"
)

func TestRsa2048Decrypt(t *testing.T) {
	tests := []struct {
		name       string
		encrypted  string
		wantErr    bool
		wantPrefix string
	}{
		{
			name:      "空字符串",
			encrypted: "",
			wantErr:   true,
		},
		{
			name:      "无效的 Base64",
			encrypted: "invalid-base64",
			wantErr:   true,
		},
		{
			name:       "正确的加密数据",
			encrypted:  "KSJDMhPIaSrxjkIEUf242bqtRYk4PdXn832Q4G7uNkA9+ZcnjsCuKjt4f/shMWbZzA+oEL3OnWKoF6yGQ1ek/jzH+Fh4USjTS2Fs62nYq5YEdgrzD3Wkogifcedi7cveAi94rRJjtfaX5aCZSicDmXGmGDMKUMXQQIhfvEfbsvwjzD65ZxXGMoN97somfOnBXX/RXwM/ZduylzBrQnvWOydH/gJnLT+Ec1Pskk8yKeYJ0Y3/fqK9BUFUwuXV8gfdnJd9SLjNy/J3r8B3v7/zUbXp1GckRL1FDDBCIEyGphMvAOf8LP/0o71gWEwU5LIAdDg5HhWlLaZMvDeQRp5RIw==",
			wantErr:    false,
			wantPrefix: "1111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decrypted, err := Rsa2048Decrypt(tt.encrypted)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rsa2048Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantPrefix != "" && !strings.HasPrefix(decrypted, tt.wantPrefix) {
				t.Errorf("Rsa2048Decrypt() = %v, want prefix %v", decrypted, tt.wantPrefix)
			}
		})
	}
}
