# pocket-go

Pocket, go implementation

## Token management

The caller owns the access token. The client does not cache one and does not log
in on your behalf: call `Login`, install the result with `SetToken`, and every
call carries it.

```go
client := pocket.New("merchant", "client-id", "client-secret", "production", terminalID)

token, err := client.Login(ctx)
if err != nil {
	return err
}
client.SetToken(token)

invoice, err := client.CreateInvoice(input)
```

`Token.ExpiresAt` is anchored to when the response arrived, because Keycloak
reports `expires_in` as a duration. `Close()` clears the installed token; it
reaches Pocket for nothing, since the token may still be in use by whoever owns
it.

Two errors are worth matching with `errors.Is`:

- `ErrNoToken` — a call was made before `SetToken`.
- `ErrUnauthorized` — Pocket refused the token (or, from `Login`, the
  credentials). Log in again and retry; the refused request was never processed.
