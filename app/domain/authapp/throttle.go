package authapi

import (
	"strings"
	"sync"
	"time"
)

// loginThrottle tracks consecutive login failures and locks a key out for a
// configured duration once the failure limit is reached. The key combines the
// email (primary) and the client IP (secondary). State is held in memory, so
// it is per-pod and resets on restart; that is acceptable while the auth
// service runs as a single instance, and is listed as a follow-up if the
// service ever scales out.
type loginThrottle struct {
	mu       sync.Mutex
	maxFails int
	lockout  time.Duration
	now      func() time.Time
	attempts map[string]loginAttempt
}

type loginAttempt struct {
	fails       int
	lockedUntil time.Time
}

func newLoginThrottle(maxFails int, lockout time.Duration) *loginThrottle {
	return &loginThrottle{
		maxFails: maxFails,
		lockout:  lockout,
		now:      time.Now,
		attempts: make(map[string]loginAttempt),
	}
}

// throttleKey combines the normalized email with the client IP.
func throttleKey(email, ip string) string {
	return strings.ToLower(strings.TrimSpace(email)) + "\x1f" + ip
}

// locked reports whether attempts for the key are currently locked out. An
// expired lock clears the entry so the next attempt starts fresh.
func (t *loginThrottle) locked(email, ip string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	key := throttleKey(email, ip)

	att, exists := t.attempts[key]
	if !exists || att.lockedUntil.IsZero() {
		return false
	}

	if t.now().Before(att.lockedUntil) {
		return true
	}

	delete(t.attempts, key)

	return false
}

// recordFailure registers a failed attempt and locks the key once the
// consecutive failure limit is reached.
func (t *loginThrottle) recordFailure(email, ip string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.evictExpired()

	key := throttleKey(email, ip)

	att := t.attempts[key]
	att.fails++

	if att.fails >= t.maxFails {
		att.lockedUntil = t.now().Add(t.lockout)
	}

	t.attempts[key] = att
}

// recordSuccess clears any failure state for the key.
func (t *loginThrottle) recordSuccess(email, ip string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.attempts, throttleKey(email, ip))
}

// evictExpired bounds the map size by dropping entries whose lock has expired
// once the map grows past a ceiling. Must be called with mu held.
func (t *loginThrottle) evictExpired() {
	const ceiling = 10000

	if len(t.attempts) < ceiling {
		return
	}

	now := t.now()
	for key, att := range t.attempts {
		if att.lockedUntil.IsZero() || now.After(att.lockedUntil) {
			delete(t.attempts, key)
		}
	}
}
