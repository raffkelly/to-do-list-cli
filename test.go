package main

func maintest() {
	testConfig := &Config{
		CurrentUser: UserData{
			JWT:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			RefreshToken: "rt_abc123def456",
			Username:     "current_user",
		},
		Users: map[string]UserData{
			"alice": {
				JWT:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				RefreshToken: "rt_alice789ghi101",
				Username:     "alice",
			},
			"bob": {
				JWT:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				RefreshToken: "rt_bob202jkl303",
				Username:     "bob",
			},
			"charlie": {
				JWT:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
				RefreshToken: "rt_charlie404mno505",
				Username:     "charlie",
			},
		},
	}

	// Test the function
	testConfig.handleUserAuthentication()
}
