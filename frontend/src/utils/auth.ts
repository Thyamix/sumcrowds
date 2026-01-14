// Error codes that indicate token issues (should trigger initaccess)
const TOKEN_ERROR_CODES = [1005, 1006, 1007, 1008]; // Invalid/expired/no/revoked refresh token

export async function auth(): Promise<void> {
	const APIURL = import.meta.env.VITE_APIURL as string;

	const response = await fetch(APIURL + "v1/auth/refreshaccess", { credentials: "include" })
	if (response.ok) {
		return
	}

	// Only call initaccess for token-specific errors, NOT for server errors (500/503)
	// This prevents creating new tokens when the database is temporarily offline
	if (response.status >= 500) {
		console.warn("Server error during auth refresh, not requesting new tokens")
		return
	}

	// Check for token-specific error codes
	try {
		const errorData = await response.json()
		if (!TOKEN_ERROR_CODES.includes(errorData.code)) {
			console.warn("Non-token error during auth refresh:", errorData)
			return
		}
	} catch {
		// If we can't parse the response, only proceed if it was a 401
		if (response.status !== 401) {
			return
		}
	}

	await fetch(APIURL + "v1/auth/initaccess", { credentials: "include" })
	return
}

/**
 * Make a fetch request with auto reauthenticate and retry
 */
export async function fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
	let response = await fetch(url, {
		...options,
		credentials: 'include'
	})

	if (response.status === 401) {
		await auth()
	} else {
		return response
	}

	response = await fetch(
		url, {
		...options,
		credentials: 'include'
	})

	return response

}

