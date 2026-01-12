import React, {useState} from 'react';
import {View, StyleSheet, TouchableOpacity, Text} from 'react-native';
import {useTranslation} from 'react-i18next';
import {Modal, Input, Button, Alert} from './ui';
import {fetchWithAuth} from '../utils/auth';
import {colors, spacing, fontSize} from '../utils/theme';

interface PinModalProps {
  visible: boolean;
  onClose: () => void;
  festivalCode: string;
  onSuccess: () => void;
}

export const PinModal: React.FC<PinModalProps> = ({
  visible,
  onClose,
  festivalCode,
  onSuccess,
}) => {
  const {t} = useTranslation();
  const [pin, setPin] = useState('');
  const [showPin, setShowPin] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handlePinChange = (text: string): void => {
    const filtered = text.replace(/[^0-9]/g, '');
    if (filtered.length <= 4) {
      setPin(filtered);
    }
  };

  const handleSubmit = async (): Promise<void> => {
    if (pin.length !== 4) {
      setError(t('pinpopup_alert'));
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await fetchWithAuth(`v1/festival/${festivalCode}/admin/access`, {
        headers: {
          'admin-pin': pin,
        },
      });

      if (response.ok) {
        setPin('');
        onSuccess();
      } else if (response.status === 403 || response.status === 422) {
        setError(t('pinpopup_alert'));
      } else {
        try {
          const errorData = await response.json();
          const errorCode = errorData.code ? ` (${errorData.code})` : '';
          setError(t('error_generic') + errorCode);
        } catch {
          setError(t('error_generic'));
        }
      }
    } catch {
      setError(t('error_generic'));
    } finally {
      setLoading(false);
    }
  };

  const handleClose = (): void => {
    setPin('');
    setError('');
    onClose();
  };

  return (
    <Modal
      visible={visible}
      onClose={handleClose}
      title={t('pinpopup_admin_access')}>
      <View style={styles.content}>
        <Text style={styles.label}>{t('pinpopup_label')}</Text>

        {error ? (
          <Alert message={error} onDismiss={() => setError('')} />
        ) : null}

        <Input
          value={pin}
          onChangeText={handlePinChange}
          placeholder={t('pinpopup_pin')}
          keyboardType="number-pad"
          maxLength={4}
          secureTextEntry={!showPin}
        />

        <TouchableOpacity
          style={styles.showPin}
          onPress={() => setShowPin(!showPin)}>
          <Text style={styles.showPinText}>
            {showPin ? '✓ ' : '○ '}
            {t('pinpopup_show_pin')}
          </Text>
        </TouchableOpacity>

        <Button
          onPress={handleSubmit}
          loading={loading}
          disabled={pin.length !== 4}
          style={styles.button}
          variant="accent">
          {t('pinpopup_confirm')}
        </Button>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  content: {
    gap: spacing.md,
  },
  label: {
    fontSize: fontSize.sm,
    color: colors.textSecondary,
    textAlign: 'center',
  },
  showPin: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  showPinText: {
    fontSize: fontSize.sm,
    color: colors.textSecondary,
  },
  button: {
    marginTop: spacing.md,
  },
});

export default PinModal;
