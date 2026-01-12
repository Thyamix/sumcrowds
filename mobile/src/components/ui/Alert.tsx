import React, {useEffect} from 'react';
import {View, Text, StyleSheet, TouchableOpacity, ViewStyle} from 'react-native';
import {colors, borderRadius, spacing, fontSize} from '../../utils/theme';

type AlertVariant = 'destructive' | 'default';

interface AlertProps {
  message: string;
  variant?: AlertVariant;
  onDismiss?: () => void;
  autoDismiss?: boolean;
  dismissTimeout?: number;
  style?: ViewStyle;
}

export const Alert: React.FC<AlertProps> = ({
  message,
  variant = 'destructive',
  onDismiss,
  autoDismiss = true,
  dismissTimeout = 7500,
  style,
}) => {
  useEffect(() => {
    if (autoDismiss && onDismiss) {
      const timer = setTimeout(() => {
        onDismiss();
      }, dismissTimeout);
      return () => clearTimeout(timer);
    }
  }, [autoDismiss, dismissTimeout, onDismiss]);

  if (!message) return null;

  const backgroundColor =
    variant === 'destructive' ? colors.destructive : colors.primary;

  return (
    <View style={[styles.container, {backgroundColor}, style]}>
      <Text style={styles.message}>{message}</Text>
      {onDismiss && (
        <TouchableOpacity onPress={onDismiss} style={styles.closeButton}>
          <Text style={styles.closeText}>×</Text>
        </TouchableOpacity>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    borderRadius: borderRadius.md,
    padding: spacing.md,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  message: {
    color: colors.white,
    fontSize: fontSize.sm,
    flex: 1,
  },
  closeButton: {
    marginLeft: spacing.sm,
    padding: spacing.xs,
  },
  closeText: {
    color: colors.white,
    fontSize: fontSize.lg,
    fontWeight: 'bold',
  },
});

export default Alert;
