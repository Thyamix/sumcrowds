import { useTranslation } from "react-i18next"
import { Globe } from "lucide-react"

const SUPPORTED_LANGUAGES = ["en", "fr"]

export function LanguageSwitcher() {
  const { i18n } = useTranslation()

  function handleChange(e) {
    const lang = e.target.value;
    i18n.changeLanguage(lang);
    localStorage.setItem('lang', lang);
  };

  return (
    <div className="relative inline-flex items-center">
      <Globe className="w-4 h-4 absolute left-2.5 text-white/70 pointer-events-none" />
      <select
        value={i18n.language}
        onChange={handleChange}
        className="h-9 pl-8 pr-3 rounded-lg bg-white/10 border border-white/20 text-white text-sm font-medium cursor-pointer hover:bg-white/20 transition-colors focus:outline-none focus:ring-2 focus:ring-white/30 appearance-none"
      >
        {SUPPORTED_LANGUAGES.map((lang) => (
          <option key={lang} value={lang} className="bg-gray-800 text-white">
            {lang.toUpperCase()}
          </option>
        ))}
      </select>
    </div>
  );
}
