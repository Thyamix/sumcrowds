export async function auth(): Promise<void> {
	const APIURL = import.meta.env.VITE_APIURL as string;

	const response = await fetch(APIURL + "v1/auth/refreshaccess", { credentials: "include" })
	if (response.ok) {
		return
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

