import React from 'react';
import {View, StyleSheet} from 'react-native';
import {colors, borderRadius, spacing} from '../../utils/theme';

export const Card = ({children, style, ...props}) => {
  return (
    <View style={[styles.card, style]} {...props}>
      {children}
    </View>
  );
};

export const CardContent = ({children, style, ...props}) => {
  return (
    <View style={[styles.cardContent, style]} {...props}>
      {children}
    </View>
  );
};

const styles = StyleSheet.create({
  card: {
    backgroundColor: colors.card,
    borderRadius: borderRadius.lg,
    borderWidth: 1,
    borderColor: colors.cardBorder,
    shadowColor: colors.black,
    shadowOffset: {width: 0, height: 2},
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  cardContent: {
    padding: spacing.lg,
  },
});

export default Card;
