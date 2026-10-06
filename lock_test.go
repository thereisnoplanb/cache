package cache

import "testing"

func TestMapKeyPanicDoesNotLeaveCacheLocked(t *testing.T) {
	tests := []struct {
		name string
		run  func(*Cache[any, string])
	}{
		{
			name: "Add",
			run: func(cache *Cache[any, string]) {
				_ = cache.Add([]int{1}, "bad key")
			},
		},
		{
			name: "Replace",
			run: func(cache *Cache[any, string]) {
				_ = cache.Replace([]int{1}, "bad key")
			},
		},
		{
			name: "AddOrReplace",
			run: func(cache *Cache[any, string]) {
				_ = cache.AddOrReplace([]int{1}, "bad key")
			},
		},
		{
			name: "Remove",
			run: func(cache *Cache[any, string]) {
				cache.Remove([]int{1})
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache, err := New[any, string](NeverExpire)
			if err != nil {
				t.Fatal(err)
			}
			if err := cache.Add("valid", "value"); err != nil {
				t.Fatal(err)
			}

			func() {
				defer func() {
					if recover() == nil {
						t.Fatal("expected map key panic")
					}
				}()
				test.run(cache)
			}()

			if value, found := cache.Get("valid"); !found || value != "value" {
				t.Fatalf("cache unusable after key panic: got (%q, %t)", value, found)
			}
		})
	}
}
