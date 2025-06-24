/** @type {string} */
const APIURL = import.meta.env.VITE_APIURL;


export async function auth() {
	const response = await fetch(APIURL + "v1/auth/refreshaccess", { credentials: "include" })
	if (response.ok) {
		return
	}
	await fetch(APIURL + "v1/auth/initaccess", { credentials: "include" })
	return
}

/**
 * Make a fetch request with auto reauthenticate and retry
 * @param {string} url 
 * @param {{}} [options={}] 
 * @returns {Promise<Response>}
*/
export async function fetchWithAuth(url, options = {}) {
	let response = await fetch(url, {
		...options,
		credentials: 'include'
	})

	if (response.status === 403 || response.status === 401) {
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

