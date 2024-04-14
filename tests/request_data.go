package tests

var (
	userRegister = map[string]interface{}{
		"NickName":        "user",
		"Email":           "user@gmail.com",
		"Password":        "password",
		"PasswordConfirm": "password",
		"Role":            "user",
	}

	adminRegister = map[string]interface{}{
		"NickName":        "admin",
		"Email":           "admin@gmail.com",
		"Password":        "password",
		"PasswordConfirm": "password",
		"Role":            "admin",
	}

	userLogin = map[string]interface{}{
		"NickName": "user",
		"Password": "password",
	}

	adminLogin = map[string]interface{}{
		"NickName": "admin",
		"Password": "password",
	}

	bannerTest1 = map[string]interface{}{
		"tag_ids":    []int{0, 1, 2},
		"feature_id": 1,
		"content": map[string]string{
			"title": "Test banner1",
			"text":  "Test text1",
			"url":   "https://test1.url",
		},
		"is_active": true,
	}

	bannerNotActive = map[string]interface{}{
		"tag_ids":    []int{0, 2},
		"feature_id": 2,
		"content": map[string]string{
			"title": "Test banner2",
			"text":  "Test text2",
			"url":   "https://test2.url",
		},
		"is_active": false,
	}

	bannerTest2 = map[string]interface{}{
		"tag_ids":    []int{0},
		"feature_id": 3,
		"content": map[string]string{
			"title": "Test banner3",
			"text":  "Test text3",
			"url":   "https://test3.url",
		},
		"is_active": true,
	}

	inactiveBannerSearch = map[string]interface{}{
		"tag_id":     2,
		"feature_id": 2,
	}

	activeBannerSearch = map[string]interface{}{
		"tag_id":     0,
		"feature_id": 3,
	}
)
