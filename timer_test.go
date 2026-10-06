package cache

import (
	"testing"
	"time"
)

func TestStaleExpirationDoesNotRemoveReplacedItem(t *testing.T) {
	instance, err := New[string, string](NeverExpire)
	if err != nil {
		t.Fatal(err)
	}

	staleToken := &expirationToken{}
	currentToken := &expirationToken{}
	instance.items["key"] = item[string]{
		Value:           "new value",
		Timer:           time.NewTimer(time.Hour),
		expirationToken: currentToken,
	}

	instance.expire("key", staleToken)

	got, found := instance.Get("key")
	if !found || got != "new value" {
		t.Fatalf("stale timer changed current item: got (%q, %t)", got, found)
	}

	instance.items["key"].Timer.Stop()
}

func TestExpirationKeyPanicDoesNotLeaveCacheLocked(t *testing.T) {
	instance, err := New[any, string](NeverExpire)
	if err != nil {
		t.Fatal(err)
	}
	if err := instance.Add("valid", "value"); err != nil {
		t.Fatal(err)
	}

	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("expected map key panic")
			}
		}()
		instance.expire([]int{1}, &expirationToken{})
	}()

	if value, found := instance.Get("valid"); !found || value != "value" {
		t.Fatalf("cache unusable after expiration key panic: got (%q, %t)", value, found)
	}
}
