// Color theme matching the web frontend
export const colors = {
  primary: '#3366CC', // Blue - Main actions
  primaryLight: '#5588EE',
  primaryDark: '#2255BB',

  secondary: '#E8EBF0', // Light gray
  secondaryDark: '#C8CBD0',

  accent: '#8833CC', // Purple - Admin/secondary actions
  accentLight: '#AA55EE',

  destructive: '#CC3333', // Red - Exit/danger
  destructiveLight: '#EE5555',

  success: '#33AA66', // Green - Enter/positive
  successLight: '#55CC88',

  warning: '#F5A623', // Orange - Warnings

  background: '#F5F7FA',
  foreground: '#1A1D24',

  card: '#FFFFFF',
  cardBorder: '#E2E8F0',

  text: '#1A1D24',
  textSecondary: '#64748B',
  textMuted: '#94A3B8',

  white: '#FFFFFF',
  black: '#000000',

  // Status colors for counter
  statusNormal: '#33AA66',
  statusWarning: '#F5A623',
  statusDanger: '#CC3333',
};

export const spacing = {
  xs: 4,
  sm: 8,
  md: 16,
  lg: 24,
  xl: 32,
  xxl: 48,
};

export const borderRadius = {
  sm: 4,
  md: 8,
  lg: 12,
  xl: 16,
  full: 9999,
};

export const fontSize = {
  xs: 12,
  sm: 14,
  md: 16,
  lg: 18,
  xl: 20,
  xxl: 24,
  huge: 48,
  massive: 72,
};

export const fontWeight = {
  normal: '400' as const,
  medium: '500' as const,
  semibold: '600' as const,
  bold: '700' as const,
};
