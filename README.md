<img width="1496" height="792" alt="image" src="https://github.com/user-attachments/assets/2e2905f4-2feb-455b-9202-47513ed02c1d" />
## JWT Authentication – Best Practices & Design Decisions

### Best Practices Followed

* Rate Limiting
* Access & Refresh Tokens
* Exponential Backoff
* Token Blacklisting / Rotation Strategy
* Hashed Passwords (Argon2)

---

## Layered Protection Strategy

* **Layer 1: IP-based limiting**
  Prevents brute-force attacks from a single source

* **Layer 2: Username / Email-based limiting**
  Protects specific user accounts from targeted attacks

* **Layer 3: Progressive penalties (Exponential Backoff)**
  Slows down repeated failed login attempts

---

## Rate Limiting Strategy

### IP-based Rate Limiting using Token Bucket

**Data Structure**

```text
Key: IP Address  
Value: { tokens, last_refill_timestamp }
```

**Configuration**

* Max tokens: 10
* Refill rate: 1 token every 6 seconds

### Flow

**New user (no existing entry):**

```text
IP → { tokens: 10, last_refill_timestamp: 1:00 }
```

**User makes a request:**

```text
IP → { tokens: 9, last_refill_timestamp: 1:00 }
```

**User exhausts all 10 tokens within 3 seconds**
For the 11th request at `1:03`:

* `current_time - last_refill_timestamp = 3s`
* New tokens added = `3 / 6 = 0`
* Total tokens = `0`
* Response: **429 – Too Many Requests**

**At 7 seconds**, when the user makes another request:

* Tokens are recalculated
* Bucket is refilled accordingly
* `last_refill_timestamp` is updated

---

## Access Token & Refresh Token Flow

* **Access Token** → Short-lived
* **Refresh Token** → Long-lived

### Login

* The Auth Service generates both access and refresh tokens on successful login
* Tokens are signed using **RS256 (asymmetric encryption)**
* The **private key** is stored only in the Auth Service and is used to sign tokens
* The **public key** is shared with other microservices for token verification

This ensures that **only the Auth Service can issue tokens**, while all other microservices can independently validate tokens without making network calls to the Auth Service.

---

### API Requests

* Client sends the **access token** with every request
* Backend validates the access token
* If access token is expired/invalid → return **401 Unauthorized**

### Token Refresh

* Client sends the **refresh token**
* If refresh token is valid → issue a new access token
* If refresh token is invalid:

  * Clear local storage
  * Cancel ongoing requests
  * Redirect user to login
  * Show “Login again” message

---

## Exponential Backoff (Login Only)

Exponential backoff is applied **only when login credentials are invalid**.

### Backoff Timing

* 2 attempts → 1 second
* 3 attempts → 2 seconds
* 4 attempts → 4 seconds
* 5 attempts → 8 seconds

### Flow

1. Request passes `rateLimiterMiddleware`
2. Auth Service checks if backoff is active for the user
3. If active → return **429 – Too Many Failed Attempts**
4. If not active:

   * Validate credentials
   * On success → reset backoff counter
   * On failure → increment failure counter

### Algorithm

```text
delay = baseDelay * 2^failCount
nextAllowedTime = currentTime + delay
```

### Redis Storage

```text
login_backoff:<user> = {
    failCount: 2,
    nextAllowed: 1733842205
}
```

* Redis keys expire automatically using TTL
* Data is cleared either:

  * After successful login, or
  * When TTL expires
* `TTL = nextAllowed - currentTime`

---

## Token Versioning (Instead of Token Blacklisting)

### Problem

* User logs in → gets access + refresh token
* Attacker steals refresh token
* User logs out
* Attacker can still use stolen refresh token until it expires (e.g., 7 days)

### Solution: Token Versioning

* Store `tokenVersion` in user data
* Include `tokenVersion` in access and refresh token payloads
* Increment `tokenVersion` on:

  * Logout
  * Password change
  * Token refresh
---

### Refresh Flow

1. Validate the refresh token
2. Check token version:

```text
if (user.tokenVersion !== payload.tokenVersion) {
    return "Token revoked"
}
```

3. If valid:

   * Increment `tokenVersion` in the database
   * Issue new access and refresh tokens with the updated `tokenVersion`

This ensures stolen refresh tokens become unusable immediately.

---


