import React, {useState} from 'react';
import {View, StyleSheet, TouchableOpacity, Text} from 'react-native';
import {useTranslation} from 'react-i18next';
import {Modal, Input, Button, Alert} from './ui';
import {fetchWithAuth} from '../utils/auth';
import {colors, spacing, fontSize} from '../utils/theme';

export const PasswordModal = ({visible, onClose, festivalCode, onSuccess}) => {
  const {t} = useTranslation();
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    if (!password) {
      setError(t('pwpopup_alert'));
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await fetchWithAuth(`v1/festival/${festivalCode}/getaccess`, {
        method: 'POST',
        body: JSON.stringify({password}),
      });

      if (response.ok) {
        setPassword('');
        onSuccess();
      } else if (response.status === 403) {
        setError(t('pwpopup_alert'));
      } else {
        try {
          const errorData = await response.json();
          const errorCode = errorData.code ? ` (${errorData.code})` : '';
          setError(t('error_generic') + errorCode);
        } catch {
          setError(t('error_generic'));
        }
      }
    } catch (err) {
      setError(t('error_generic'));
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setPassword('');
    setError('');
    onClose();
  };

  return (
    <Modal
      visible={visible}
      onClose={handleClose}
      title={t('pwpopup_header')}>
      <View style={styles.content}>
        <Text style={styles.festivalText}>
          {t('pwpopup_festival')}: {festivalCode}
        </Text>

        {error ? (
          <Alert message={error} onDismiss={() => setError('')} />
        ) : null}

        <Input
          label={t('pwpopup_password')}
          value={password}
          onChangeText={setPassword}
          placeholder={t('pwpopup_password')}
          secureTextEntry={!showPassword}
        />

        <TouchableOpacity
          style={styles.showPassword}
          onPress={() => setShowPassword(!showPassword)}>
          <Text style={styles.showPasswordText}>
            {showPassword ? '✓ ' : '○ '}
            {t('pwpopup_show_password')}
          </Text>
        </TouchableOpacity>

        <Button
          onPress={handleSubmit}
          loading={loading}
          disabled={!password}
          style={styles.button}
          variant="success">
          {t('pwpopup_confirm')}
        </Button>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  content: {
    gap: spacing.md,
  },
  festivalText: {
    fontSize: fontSize.sm,
    color: colors.textSecondary,
    textAlign: 'center',
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

export default PasswordModal;
