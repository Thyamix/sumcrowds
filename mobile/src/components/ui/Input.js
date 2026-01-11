import React from 'react';
import {View, TextInput, Text, StyleSheet} from 'react-native';
import {colors, borderRadius, spacing, fontSize} from '../../utils/theme';

export const Input = ({
  label,
  value,
  onChangeText,
  placeholder,
  secureTextEntry = false,
  keyboardType = 'default',
  maxLength,
  error,
  style,
  inputStyle,
  ...props
}) => {
  return (
    <View style={[styles.container, style]}>
      {label && <Text style={styles.label}>{label}</Text>}
      <TextInput
        style={[
          styles.input,
          error && styles.inputError,
          inputStyle,
        ]}
        value={value}
        onChangeText={onChangeText}
        placeholder={placeholder}
        placeholderTextColor={colors.textMuted}
        secureTextEntry={secureTextEntry}
        keyboardType={keyboardType}
        maxLength={maxLength}
        autoCapitalize="none"
        {...props}
      />
      {error && <Text style={styles.errorText}>{error}</Text>}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    marginBottom: spacing.md,
  },
  label: {
    fontSize: fontSize.sm,
    color: colors.textSecondary,
    marginBottom: spacing.xs,
    fontWeight: '500',
  },
  input: {
    backgroundColor: colors.white,
    borderWidth: 1,
    borderColor: colors.cardBorder,
    borderRadius: borderRadius.md,
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.md,
    fontSize: fontSize.md,
    color: colors.text,
  },
  inputError: {
    borderColor: colors.destructive,
  },
  errorText: {
    fontSize: fontSize.sm,
    color: colors.destructive,
    marginTop: spacing.xs,
  },
});

export default Input;
