import { useTranslation } from "react-i18next"


const SUPPORTED_LANGUAGES = ["en", "fr"]

export default function LanguageSwitcher() {
  const { i18n } = useTranslation()
  function handleChange(e) {
    const lang = e.target.value;
    i18n.changeLanguage(lang);
    localStorage.setItem('lang', lang);
  };

  return (
    <div style={{ display: "flex", justifyContent: "flex-end" }}>
      <select value={i18n.language} onChange={handleChange}>
        {SUPPORTED_LANGUAGES.map((lang) => (
          <option key={lang} value={lang}>
            {lang === 'en' ? 'English' : 'Français'}
          </option>
        ))}
      </select>
    </div >
  );
}
