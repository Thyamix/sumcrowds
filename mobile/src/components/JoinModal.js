import React, {useState} from 'react';
import {View, StyleSheet} from 'react-native';
import {useTranslation} from 'react-i18next';
import {Modal, Input, Button, Alert} from './ui';
import {fetchWithAuth} from '../utils/auth';
import {spacing} from '../utils/theme';

export const JoinModal = ({visible, onClose, onJoin}) => {
  const {t} = useTranslation();
  const [code, setCode] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleCodeChange = text => {
    // Only allow alphanumeric, max 6 characters
    const filtered = text.replace(/[^a-zA-Z0-9]/g, '').toUpperCase();
    if (filtered.length <= 6) {
      setCode(filtered);
    }
  };

  const handleJoin = async () => {
    if (code.length !== 6) {
      setError(t('joinpopup_alert'));
      return;
    }

    setLoading(true);
    setError('');

    try {
      const response = await fetchWithAuth(`v1/festival/${code}/exists`);
      if (response.status === 404) {
        setError(t('joinpopup_alert'));
      } else if (response.ok) {
        onJoin(code);
        setCode('');
        onClose();
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
    setCode('');
    setError('');
    onClose();
  };

  return (
    <Modal
      visible={visible}
      onClose={handleClose}
      title={t('joinpopup_header')}>
      <View style={styles.content}>
        {error ? (
          <Alert message={error} onDismiss={() => setError('')} />
        ) : null}

        <Input
          label={t('joinpopup_code_label')}
          value={code}
          onChangeText={handleCodeChange}
          placeholder={t('joinpopup_enter_code')}
          maxLength={6}
          autoCapitalize="characters"
        />

        <Button
          onPress={handleJoin}
          loading={loading}
          disabled={code.length !== 6}
          style={styles.button}
          variant="success">
          {t('joinpopup_confirm')}
        </Button>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  content: {
    gap: spacing.md,
  },
  button: {
    marginTop: spacing.md,
  },
});

export default JoinModal;
