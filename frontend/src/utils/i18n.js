import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import Backend from 'i18next-http-backend';

const currentLang = localStorage.getItem("lang") || "en"

export function initI18n() {
	return i18n
		.use(initReactI18next)
		.use(Backend)
		.init({
			fallbackLng: 'en',
			debug: false,
			interpolation: {
				escapeValue: false,
			},
			lng: currentLang,
			backend: {
				loadPath: '/locales/{{lng}}/translation.json',
			},
			react: {
				useSuspense: false,
			},
		});
}

