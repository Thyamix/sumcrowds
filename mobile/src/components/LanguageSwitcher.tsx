import React, {useState} from 'react';
import {View, Text, TouchableOpacity, StyleSheet, ViewStyle} from 'react-native';
import {useTranslation} from 'react-i18next';
import {changeLanguage, SupportedLanguage} from '../utils/i18n';
import {colors, borderRadius, spacing, fontSize} from '../utils/theme';

interface Language {
  code: SupportedLanguage;
  label: string;
}

interface LanguageSwitcherProps {
  style?: ViewStyle;
}

const languages: Language[] = [
  {code: 'en', label: 'EN'},
  {code: 'fr', label: 'FR'},
];

export const LanguageSwitcher: React.FC<LanguageSwitcherProps> = ({style}) => {
  const {i18n} = useTranslation();
  const [isOpen, setIsOpen] = useState(false);

  const currentLang = i18n.language || 'en';

  const handleSelect = async (code: SupportedLanguage): Promise<void> => {
    await changeLanguage(code);
    setIsOpen(false);
  };

  return (
    <View style={[styles.container, style]}>
      <TouchableOpacity
        style={styles.button}
        onPress={() => setIsOpen(!isOpen)}>
        <Text style={styles.globe}>🌐</Text>
        <Text style={styles.buttonText}>{currentLang.toUpperCase()}</Text>
      </TouchableOpacity>

      {isOpen && (
        <View style={styles.dropdown}>
          {languages.map(lang => (
            <TouchableOpacity
              key={lang.code}
              style={[
                styles.option,
                currentLang === lang.code && styles.optionActive,
              ]}
              onPress={() => handleSelect(lang.code)}>
              <Text
                style={[
                  styles.optionText,
                  currentLang === lang.code && styles.optionTextActive,
                ]}>
                {lang.label}
              </Text>
            </TouchableOpacity>
          ))}
        </View>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    position: 'relative',
    zIndex: 100,
  },
  button: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: colors.white,
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
    borderRadius: borderRadius.md,
    borderWidth: 1,
    borderColor: colors.cardBorder,
  },
  globe: {
    fontSize: fontSize.md,
    marginRight: spacing.xs,
  },
  buttonText: {
    fontSize: fontSize.sm,
    color: colors.text,
    fontWeight: '500',
  },
  dropdown: {
    position: 'absolute',
    top: '100%',
    right: 0,
    marginTop: spacing.xs,
    backgroundColor: colors.white,
    borderRadius: borderRadius.md,
    borderWidth: 1,
    borderColor: colors.cardBorder,
    shadowColor: colors.black,
    shadowOffset: {width: 0, height: 2},
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
    minWidth: 80,
  },
  option: {
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
  },
  optionActive: {
    backgroundColor: colors.primary,
  },
  optionText: {
    fontSize: fontSize.sm,
    color: colors.text,
    textAlign: 'center',
  },
  optionTextActive: {
    color: colors.white,
    fontWeight: '600',
  },
});

export default LanguageSwitcher;
