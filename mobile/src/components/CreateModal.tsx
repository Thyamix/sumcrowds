import React, {useState} from 'react';
import {View, StyleSheet, TouchableOpacity, Text} from 'react-native';
import {useTranslation} from 'react-i18next';
import {Modal, Input, Button, Alert} from './ui';
import {fetchWithAuth} from '../utils/auth';
import {colors, spacing, fontSize} from '../utils/theme';

interface CreateModalProps {
  visible: boolean;
  onClose: () => void;
  onCreated: (code: string) => void;
}

export const CreateModal: React.FC<CreateModalProps> = ({visible, onClose, onCreated}) => {
  const {t} = useTranslation();
  const [pin, setPin] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handlePinChange = (text: string): void => {
    const filtered = text.replace(/[^0-9]/g, '');
    if (filtered.length <= 4) {
      setPin(filtered);
    }
  };

  const handleCreate = async (): Promise<void> => {
    if (pin.length !== 4) {
      setError(t('pinpopup_alert'));
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await fetchWithAuth('v1/create/festival', {
        method: 'POST',
        body: JSON.stringify({
          password: password || '',
          pin: pin,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        if (data.type === 'success' && data.content) {
          onCreated(data.content);
          setPin('');
          setPassword('');
          onClose();
        } else {
          setError(`${t('error_generic')} (unexpected response)`);
        }
      } else {
        try {
          const errorData = await response.json();
          const errorCode = errorData.code ? ` (${errorData.code})` : '';
          setError(t('error_generic') + errorCode);
        } catch {
          setError(`${t('error_generic')} (HTTP ${response.status})`);
        }
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Network error';
      setError(`${t('error_generic')} (${message})`);
    } finally {
      setLoading(false);
    }
  };

  const handleClose = (): void => {
    setPin('');
    setPassword('');
    setError('');
    onClose();
  };

  return (
    <Modal
      visible={visible}
      onClose={handleClose}
      title={t('createpopup_header')}>
      <View style={styles.content}>
        {error ? (
          <Alert message={error} onDismiss={() => setError('')} />
        ) : null}

        <Input
          label={t('createpopup_pin_label')}
          value={pin}
          onChangeText={handlePinChange}
          placeholder={t('createpopup_admin_pin')}
          keyboardType="number-pad"
          maxLength={4}
          secureTextEntry
        />

        <Input
          label={t('createpopup_password_label')}
          value={password}
          onChangeText={setPassword}
          placeholder={t('createpopup_password')}
          secureTextEntry={!showPassword}
        />

        <TouchableOpacity
          style={styles.showPassword}
          onPress={() => setShowPassword(!showPassword)}>
          <Text style={styles.showPasswordText}>
            {showPassword ? '✓ ' : '○ '}
            {t('createpopup_show_password')}
          </Text>
        </TouchableOpacity>

        <Button
          onPress={handleCreate}
          loading={loading}
          disabled={pin.length !== 4}
          style={styles.button}
          variant="success">
          {t('createpopup_confirm')}
        </Button>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  content: {
    gap: spacing.md,
  },
  showPassword: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  showPasswordText: {
    fontSize: fontSize.sm,
    color: colors.textSecondary,
  },
  button: {
    marginTop: spacing.md,
  },
});

export default CreateModal;
