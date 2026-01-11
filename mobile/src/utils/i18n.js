import i18n from 'i18next';
import {initReactI18next} from 'react-i18next';
import AsyncStorage from '@react-native-async-storage/async-storage';

import en from '../locales/en/translation.json';
import fr from '../locales/fr/translation.json';

const LANGUAGE_KEY = 'lang';

const resources = {
  en: {translation: en},
  fr: {translation: fr},
};

const initI18n = async () => {
  let savedLanguage = 'en';
  try {
    const stored = await AsyncStorage.getItem(LANGUAGE_KEY);
    if (stored) {
      savedLanguage = stored;
    }
  } catch (error) {
    console.log('Error loading language:', error);
  }

  await i18n.use(initReactI18next).init({
    resources,
    lng: savedLanguage,
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false,
    },
    react: {
      useSuspense: false,
    },
  });

  return i18n;
};

export const changeLanguage = async lng => {
  try {
    await AsyncStorage.setItem(LANGUAGE_KEY, lng);
    await i18n.changeLanguage(lng);
  } catch (error) {
    console.log('Error saving language:', error);
  }
};

export {initI18n};
export default i18n;
