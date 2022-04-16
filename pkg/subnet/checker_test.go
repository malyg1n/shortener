package subnet

import "testing"

func TestCheckSubnet(t *testing.T) {
	tests := []struct {
		name   string
		ip     string
		subnet string
		want   bool
	}{
		{
			name:   "valid",
			ip:     "127.0.0.1",
			subnet: "127.0.0.0/16",
			want:   true,
		},
		{
			name:   "invalid",
			ip:     "127.1.0.1",
			subnet: "127.0.0.0/16",
			want:   false,
		},
		{
			name:   "empty ip",
			ip:     "",
			subnet: "127.0.0.0/16",
			want:   false,
		},
		{
			name:   "empty snet",
			ip:     "127.0.0.1",
			subnet: "",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSubnet(tt.ip, tt.subnet); got != tt.want {
				t.Errorf("CheckSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}
