package app

// func TestNew(t *testing.T) {
// 	server, err := New()
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}

// 	if server.storage == nil {
// 		t.Fatal("expected storage to not be nil")
// 	}

// 	if server.engine == nil {
// 		t.Fatal("expected router to not be nil")
// 	}
// }

// func TestParseFlags(t *testing.T) {

// 	tests := []struct {
// 		name    string
// 		args    []string
// 		want    Config
// 		wantErr bool
// 	}{
// 		{
// 			name:    "shortcut",
// 			args:    []string{"-a0.0.0.0:8080", "-r1.1.1.1:8383"},
// 			want:    Config{GophermartAddress: "0.0.0.0:8080", AccrualAddress: "1.1.1.1:8383"},
// 			wantErr: false,
// 		},
// 		{
// 			name:    "shortcut",
// 			args:    []string{"--gophermart-address=127.0.0.1:81", "--accrual-address=1.1.1.1:8383"},
// 			want:    Config{GophermartAddress: "127.0.0.1:81", AccrualAddress: "1.1.1.1:8383"},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			config := &Config{GophermartAddress: "default", AccrualAddress: "default"}

// 			err := parseFlags(config, "progname", tt.args)
// 			if (err != nil) != tt.wantErr {
// 				t.Fatalf("Expected no error, got %v", err)
// 			}
// 			if tt.want != *config {
// 				t.Errorf("Expected %v, got %v", tt.want, config)
// 			}
// 		})
// 	}
// }

// func TestRun(t *testing.T) {
// 	// pending: how to test lsitenAndServe? goroutine?
// }

// func TestString(t *testing.T) {
// 	// config := Config{Address: "0.0.0.0:8080"}
// 	srv, _ := New() // TODO

// 	expected := "server config: gophermart-address=0.0.0.0:8080; accrual-address=0.0.0.0:8282; database="
// 	if srv.String() != expected {
// 		t.Errorf("Expected %v, got %v", expected, srv.String())
// 	}
// }
