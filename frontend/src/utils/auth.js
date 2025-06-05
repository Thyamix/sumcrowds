/** @type {string} */
const APIURL = import.meta.env.VITE_APIURL;


export async function auth() {
	let response = await fetch(APIURL + "v1/auth/validateaccess", { credentials: "include" })
	if (response.ok) {
		return
	}
	response = await fetch(APIURL + "v1/auth/refreshaccess", { credentials: "include" })
	if (response.ok) {
		return
	}
	response = await fetch(APIURL + "v1/auth/initaccess", { credentials: "include" })
	return
}

/**
 * Make a fetch request with auto reauthenticate and retry
 * @param {string} url 
 * @param {{}} [options={}] 
 * @returns {Promise<Response>}
*/
export async function fetchWithAuth(url, options = {}) {
	const response = await fetch(url, {
		...options,
		credentials: 'include'
	})

	if (response.status === 403 || response.status === 401) {
		await auth()
	}

	const retryResponse = await fetch(
		url, {
		...options,
		credentials: 'include'
	})

	return retryResponse

}

