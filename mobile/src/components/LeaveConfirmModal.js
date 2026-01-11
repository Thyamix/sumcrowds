import React from 'react';
import {View, Text, StyleSheet} from 'react-native';
import {useTranslation} from 'react-i18next';
import {Modal, Button} from './ui';
import {colors, spacing, fontSize} from '../utils/theme';

export const LeaveConfirmModal = ({visible, onClose, onConfirm}) => {
  const {t} = useTranslation();

  return (
    <Modal
      visible={visible}
      onClose={onClose}
      title={t('leave_confirm_title')}
      showCloseButton={false}>
      <View style={styles.content}>
        <Text style={styles.message}>{t('leave_confirm_message')}</Text>

        <View style={styles.buttons}>
          <Button
            onPress={onClose}
            variant="secondary"
            style={styles.button}>
            {t('leave_confirm_cancel')}
          </Button>
          <Button
            onPress={onConfirm}
            variant="destructive"
            style={styles.button}>
            {t('leave_confirm_leave')}
          </Button>
        </View>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  content: {
    gap: spacing.lg,
  },
  message: {
    fontSize: fontSize.md,
    color: colors.textSecondary,
    textAlign: 'center',
    lineHeight: 24,
  },
  buttons: {
    flexDirection: 'row',
    gap: spacing.md,
  },
  button: {
    flex: 1,
  },
});

export default LeaveConfirmModal;
