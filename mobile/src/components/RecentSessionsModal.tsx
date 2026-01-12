import React, {useState, useEffect} from 'react';
import {View, Text, StyleSheet, TouchableOpacity, ActivityIndicator} from 'react-native';
import {useTranslation} from 'react-i18next';
import {Modal, Button, Alert} from './ui';
import {fetchWithAuth} from '../utils/auth';
import {colors, spacing, fontSize, fontWeight, borderRadius} from '../utils/theme';

interface Session {
  code: string;
  last_used_at: number;
}

interface RecentSessionsModalProps {
  visible: boolean;
  onClose: () => void;
  onSelect: (code: string) => void;
}

export const RecentSessionsModal: React.FC<RecentSessionsModalProps> = ({
  visible,
  onClose,
  onSelect,
}) => {
  const {t} = useTranslation();
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  useEffect(() => {
    if (visible) {
      loadSessions();
    }
  }, [visible]);

  const loadSessions = async (): Promise<void> => {
    setLoading(true);
    setError('');

    try {
      const response = await fetchWithAuth('v1/user/sessions');
      if (response.ok) {
        const data = await response.json();
        setSessions(data || []);
      } else {
        setError(t('error_generic'));
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Network error';
      setError(`${t('error_generic')} (${message})`);
    } finally {
      setLoading(false);
    }
  };

  const formatDate = (timestamp: number): string => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleDateString(undefined, {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const handleSelect = (code: string): void => {
    onSelect(code);
    onClose();
  };

  return (
    <Modal visible={visible} onClose={onClose} title={t('recent_sessions_title')}>
      <View style={styles.content}>
        {error ? (
          <Alert message={error} onDismiss={() => setError('')} />
        ) : null}

        {loading ? (
          <View style={styles.loadingContainer}>
            <ActivityIndicator size="large" color={colors.primary} />
          </View>
        ) : sessions.length === 0 ? (
          <Text style={styles.emptyText}>{t('recent_sessions_empty')}</Text>
        ) : (
          <View style={styles.sessionList}>
            {sessions.map((session) => (
              <TouchableOpacity
                key={session.code}
                style={styles.sessionItem}
                onPress={() => handleSelect(session.code)}>
                <View style={styles.sessionInfo}>
                  <Text style={styles.sessionCode}>{session.code}</Text>
                  <Text style={styles.sessionDate}>
                    {formatDate(session.last_used_at)}
                  </Text>
                </View>
                <Text style={styles.arrow}>→</Text>
              </TouchableOpacity>
            ))}
          </View>
        )}

        <Button
          onPress={onClose}
          variant="outline"
          style={styles.closeButton}>
          {t('recent_sessions_close')}
        </Button>
      </View>
    </Modal>
  );
};

const styles = StyleSheet.create({
  content: {
    gap: spacing.md,
  },
  loadingContainer: {
    paddingVertical: spacing.xl,
    alignItems: 'center',
  },
  emptyText: {
    fontSize: fontSize.md,
    color: colors.textSecondary,
    textAlign: 'center',
    paddingVertical: spacing.lg,
  },
  sessionList: {
    gap: spacing.sm,
  },
  sessionItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: spacing.md,
    backgroundColor: colors.secondary,
    borderRadius: borderRadius.md,
  },
  sessionInfo: {
    flex: 1,
  },
  sessionCode: {
    fontSize: fontSize.lg,
    fontWeight: fontWeight.semibold,
    color: colors.text,
    fontFamily: 'monospace',
  },
  sessionDate: {
    fontSize: fontSize.sm,
    color: colors.textSecondary,
    marginTop: spacing.xs,
  },
  arrow: {
    fontSize: fontSize.lg,
    color: colors.primary,
    marginLeft: spacing.md,
  },
  closeButton: {
    marginTop: spacing.md,
  },
});

export default RecentSessionsModal;
