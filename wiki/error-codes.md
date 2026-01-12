# API Error Codes

Error codes are returned in the JSON response body when an API error occurs:

```json
{
  "error": "error message",
  "code": 1001
}
```

## Error Code Ranges

| Range | Category |
|-------|----------|
| 1xxx | Authentication errors |
| 2xxx | Festival/session errors |
| 3xxx | Request validation errors |
| 4xxx | Internal server errors |

## Authentication Errors (1xxx)

| Code | Error | Description |
|------|-------|-------------|
| 1001 | Invalid access token | The provided access token is malformed or invalid |
| 1002 | Expired access token | The access token has expired, refresh required |
| 1003 | Revoked access token | The access token has been revoked |
| 1004 | No access token | No access token was provided in the request |
| 1005 | Invalid refresh token | The provided refresh token is malformed or invalid |
| 1006 | Expired refresh token | The refresh token has expired, re-authentication required |
| 1007 | No refresh token | No refresh token was provided |
| 1008 | Revoked refresh token | The refresh token has been revoked |

## Festival/Session Errors (2xxx)

| Code | Error | Description |
|------|-------|-------------|
| 2001 | Festival not found | The festival code does not exist |
| 2002 | Invalid password | Wrong password for festival access |
| 2003 | No access | User does not have access to this festival |
| 2004 | Expired access | Festival access has expired |
| 2005 | Invalid PIN | Wrong admin PIN provided |

## Request Validation Errors (3xxx)

| Code | Error | Description |
|------|-------|-------------|
| 3001 | Invalid request | The request body could not be read |
| 3002 | Invalid JSON | The request body is not valid JSON |
| 3003 | Invalid amount | Counter increment/decrement amount must be between 1 and 100 |

## Internal Server Errors (4xxx)

| Code | Error | Description |
|------|-------|-------------|
| 4000 | Internal error | Generic internal server error |
| 4001 | Mismatched lengths | Internal data consistency error in archived events |
| 4002 | Failed encode response | Failed to encode JSON response |
| 4003 | Failed marshal | Failed to marshal JSON |
| 4004 | Failed add value | Failed to update counter in database |
| 4005 | Failed get total | Failed to retrieve counter total from database |
| 4006 | Failed hash password | Failed to hash password |
| 4007 | Failed reset festival | Failed to reset/archive festival data |

## Source Files

- Backend implementation: `backend/counter/internal/apperrors/errors.go`
- Mobile error display: `mobile/src/components/*.js` (CreateModal, JoinModal, PasswordModal, PinModal)
