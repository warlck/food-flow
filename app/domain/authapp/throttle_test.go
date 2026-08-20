package authapi

import (
	"fmt"
	"testing"
	"time"
)

func Test_LoginThrottle(t *testing.T) {
	t.Parallel()

	th := newLoginThrottle(3, 15*time.Minute)

	now := time.Now()
	th.now = func() time.Time { return now }

	email := "operator@example.com"
	ip := "203.0.113.10"

	// Failures below the limit do not lock.
	th.recordFailure(email, ip)
	th.recordFailure(email, ip)

	if th.locked(email, ip) {
		t.Fatal("should not be locked below the failure limit")
	}

	// Reaching the limit locks.
	th.recordFailure(email, ip)

	if !th.locked(email, ip) {
		t.Fatal("should be locked once the failure limit is reached")
	}

	// A different email or IP is not affected.
	if th.locked("other@example.com", ip) {
		t.Error("different email should not be locked")
	}

	if th.locked(email, "198.51.100.20") {
		t.Error("different IP should not be locked")
	}

	// The email is matched case-insensitively.
	if !th.locked("Operator@Example.com", ip) {
		t.Error("email matching should be case-insensitive")
	}

	// Still locked just before the lockout expires.
	now = now.Add(14 * time.Minute)

	if !th.locked(email, ip) {
		t.Error("should still be locked before the lockout expires")
	}

	// The lock expires after the lockout duration.
	now = now.Add(2 * time.Minute)

	if th.locked(email, ip) {
		t.Error("lock should expire after the lockout duration")
	}

	// A success resets the failure count.
	th.recordFailure(email, ip)
	th.recordFailure(email, ip)
	th.recordSuccess(email, ip)
	th.recordFailure(email, ip)
	th.recordFailure(email, ip)

	if th.locked(email, ip) {
		t.Error("success should reset the consecutive failure count")
	}
}

func Test_LoginThrottleEviction(t *testing.T) {
	t.Parallel()

	th := newLoginThrottle(1, time.Minute)

	now := time.Now()
	th.now = func() time.Time { return now }

	// Fill the map past the eviction ceiling with entries whose locks have
	// all expired.
	for i := range 10001 {
		key := throttleKey(fmt.Sprintf("user%d@example.com", i), "203.0.113.1")
		th.attempts[key] = loginAttempt{fails: 1, lockedUntil: now.Add(-time.Minute)}
	}

	th.recordFailure("new@example.com", "203.0.113.1")

	if len(th.attempts) != 1 {
		t.Fatalf("expected expired entries to be evicted, map size = %d", len(th.attempts))
	}
}
