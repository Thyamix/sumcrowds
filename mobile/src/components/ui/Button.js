import React from 'react';
import {
  TouchableOpacity,
  Text,
  StyleSheet,
  ActivityIndicator,
} from 'react-native';
import {colors, borderRadius, spacing, fontSize, fontWeight} from '../../utils/theme';

const variants = {
  default: {
    backgroundColor: colors.primary,
    textColor: colors.white,
    borderColor: colors.primary,
  },
  secondary: {
    backgroundColor: colors.secondary,
    textColor: colors.text,
    borderColor: colors.secondary,
  },
  destructive: {
    backgroundColor: colors.destructive,
    textColor: colors.white,
    borderColor: colors.destructive,
  },
  success: {
    backgroundColor: colors.success,
    textColor: colors.white,
    borderColor: colors.success,
  },
  outline: {
    backgroundColor: 'transparent',
    textColor: colors.primary,
    borderColor: colors.primary,
  },
  ghost: {
    backgroundColor: 'transparent',
    textColor: colors.text,
    borderColor: 'transparent',
  },
  accent: {
    backgroundColor: colors.accent,
    textColor: colors.white,
    borderColor: colors.accent,
  },
};

const sizes = {
  sm: {
    paddingVertical: spacing.sm,
    paddingHorizontal: spacing.md,
    fontSize: fontSize.sm,
  },
  md: {
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.lg,
    fontSize: fontSize.md,
  },
  lg: {
    paddingVertical: spacing.lg,
    paddingHorizontal: spacing.xl,
    fontSize: fontSize.lg,
  },
  xl: {
    paddingVertical: spacing.xl,
    paddingHorizontal: spacing.xxl,
    fontSize: fontSize.xl,
  },
  counter: {
    paddingVertical: spacing.md,
    paddingHorizontal: spacing.md,
    fontSize: fontSize.xxl,
    minWidth: 80,
    minHeight: 80,
  },
};

export const Button = ({
  children,
  variant = 'default',
  size = 'md',
  disabled = false,
  loading = false,
  onPress,
  style,
  textStyle,
  ...props
}) => {
  const variantStyles = variants[variant] || variants.default;
  const sizeStyles = sizes[size] || sizes.md;

  return (
    <TouchableOpacity
      style={[
        styles.button,
        {
          backgroundColor: variantStyles.backgroundColor,
          borderColor: variantStyles.borderColor,
          paddingVertical: sizeStyles.paddingVertical,
          paddingHorizontal: sizeStyles.paddingHorizontal,
          minWidth: sizeStyles.minWidth,
          minHeight: sizeStyles.minHeight,
        },
        disabled && styles.disabled,
        style,
      ]}
      onPress={onPress}
      disabled={disabled || loading}
      activeOpacity={0.7}
      {...props}>
      {loading ? (
        <ActivityIndicator color={variantStyles.textColor} />
      ) : (
        <Text
          style={[
            styles.text,
            {
              color: variantStyles.textColor,
              fontSize: sizeStyles.fontSize,
            },
            textStyle,
          ]}>
          {children}
        </Text>
      )}
    </TouchableOpacity>
  );
};

const styles = StyleSheet.create({
  button: {
    borderRadius: borderRadius.lg,
    borderWidth: 1,
    alignItems: 'center',
    justifyContent: 'center',
    flexDirection: 'row',
  },
  text: {
    fontWeight: fontWeight.semibold,
    textAlign: 'center',
  },
  disabled: {
    opacity: 0.5,
  },
});

export default Button;
