import React, {useState, useEffect} from 'react';
import {
  View,
  Text,
  StyleSheet,
  SafeAreaView,
  StatusBar,
  ScrollView,
  TouchableOpacity,
  Linking,
} from 'react-native';
import {useTranslation} from 'react-i18next';
import {Button, Input, Card, CardContent, Alert} from '../components/ui';
import {LanguageSwitcher, PinModal} from '../components';
import {fetchWithAuth, auth, getAccessToken} from '../utils/auth';
import {API_URL} from '../config';
import {colors, spacing, fontSize, fontWeight, borderRadius} from '../utils/theme';

export const AdminScreen = ({route, navigation}) => {
  const {festivalCode} = route.params;
  const {t} = useTranslation();

  const [alert, setAlert] = useState('');
  const [inputValue, setInputValue] = useState('');
  const [archives, setArchives] = useState([]);
  const [showPinModal, setShowPinModal] = useState(false);
  const [hasAccess, setHasAccess] = useState(false);
  const [loading, setLoading] = useState(false);

  // Check admin access on mount
  useEffect(() => {
    const checkAccess = async () => {
      await auth();
      try {
        const response = await fetchWithAuth(`v1/festival/${festivalCode}/admin/access`);
        if (response.status === 422 || response.status === 403) {
          setShowPinModal(true);
        } else if (response.ok) {
          setHasAccess(true);
        } else {
          navigation.goBack();
        }
      } catch (err) {
        console.error('Admin access check error:', err);
        navigation.goBack();
      }
    };

    checkAccess();
  }, [festivalCode, navigation]);

  // Load archives when access is granted
  useEffect(() => {
    if (hasAccess) {
      loadArchives();
    }
  }, [hasAccess]);

  const loadArchives = async () => {
    try {
      const response = await fetchWithAuth(
        `v1/festival/${festivalCode}/admin/getarchivedevents`,
      );
      if (response.ok) {
        const data = await response.json();
        setArchives(data || []);
      }
    } catch (err) {
      console.error('Load archives error:', err);
    }
  };

  const handleInputChange = text => {
    // Only allow numbers
    const filtered = text.replace(/[^0-9]/g, '');
    setInputValue(filtered);
  };

  const handleSetGauge = async () => {
    if (!inputValue) {
      setAlert(t('admin_numbers_only'));
      return;
    }

    setLoading(true);
    try {
      const response = await fetchWithAuth(
        `v1/festival/${festivalCode}/admin/setgauge`,
        {
          method: 'POST',
          body: JSON.stringify({max: parseInt(inputValue, 10)}),
        },
      );

      if (response.ok) {
        setInputValue('');
      } else {
        setAlert(t('error_generic'));
      }
    } catch (err) {
      setAlert(t('error_generic'));
    } finally {
      setLoading(false);
    }
  };

  const handleArchive = async () => {
    setLoading(true);
    try {
      const response = await fetchWithAuth(
        `v1/festival/${festivalCode}/admin/archivecurrentevent`,
      );

      if (response.ok) {
        await loadArchives();
      } else {
        setAlert(t('error_generic'));
      }
    } catch (err) {
      setAlert(t('error_generic'));
    } finally {
      setLoading(false);
    }
  };

  const handleDownloadCurrent = async () => {
    try {
      const token = await getAccessToken();
      const url = `${API_URL}v1/festival/${festivalCode}/admin/download/activecsv`;
      Linking.openURL(url + (token ? `?token=${token}` : ''));
    } catch (err) {
      setAlert(t('error_generic'));
    }
  };

  const handleDownloadArchive = async id => {
    try {
      const token = await getAccessToken();
      const url = `${API_URL}v1/festival/${festivalCode}/admin/download/archivedcsv/${id}`;
      Linking.openURL(url + (token ? `?token=${token}` : ''));
    } catch (err) {
      setAlert(t('error_generic'));
    }
  };

  const handlePinSuccess = () => {
    setShowPinModal(false);
    setHasAccess(true);
  };

  if (!hasAccess && !showPinModal) {
    return null;
  }

  return (
    <SafeAreaView style={styles.container}>
      <StatusBar barStyle="dark-content" backgroundColor={colors.background} />

      <View style={styles.header}>
        <TouchableOpacity onPress={() => navigation.goBack()}>
          <Text style={styles.backButton}>← {t('counter_code')}: {festivalCode}</Text>
        </TouchableOpacity>
        <LanguageSwitcher />
      </View>

      <ScrollView style={styles.content} contentContainerStyle={styles.contentContainer}>
        <Text style={styles.title}>{t('admin_title')}</Text>

        {alert ? (
          <Alert message={alert} onDismiss={() => setAlert('')} />
        ) : null}

        {/* Set Capacity */}
        <Card style={styles.card}>
          <CardContent>
            <Text style={styles.sectionTitle}>{t('admin_set_capacity')}</Text>
            <View style={styles.row}>
              <Input
                value={inputValue}
                onChangeText={handleInputChange}
                placeholder={t('admin_max_gauge')}
                keyboardType="number-pad"
                style={styles.input}
              />
              <Button
                onPress={handleSetGauge}
                loading={loading}
                variant="accent"
                style={styles.setButton}>
                {t('admin_set_gauge')}
              </Button>
            </View>
          </CardContent>
        </Card>

        {/* Current Event */}
        <Card style={styles.card}>
          <CardContent>
            <Text style={styles.sectionTitle}>{t('admin_current_event')}</Text>
            <View style={styles.buttonRow}>
              <Button
                onPress={handleArchive}
                loading={loading}
                variant="secondary"
                style={styles.actionButton}>
                {t('admin_archive')}
              </Button>
              <Button
                onPress={handleDownloadCurrent}
                variant="accent"
                style={styles.actionButton}>
                {t('admin_get_csv')}
              </Button>
            </View>
          </CardContent>
        </Card>

        {/* Archived Data */}
        <Card style={styles.card}>
          <CardContent>
            <Text style={styles.sectionTitle}>{t('admin_archived_data')}</Text>
            {archives.length === 0 ? (
              <Text style={styles.noArchives}>{t('admin_no_archives')}</Text>
            ) : (
              <View style={styles.archiveList}>
                {archives.map((archive, index) => (
                  <View key={archive.id || index} style={styles.archiveItem}>
                    <Text style={styles.archiveText}>
                      #{index + 1} - ID: {archive.id}
                    </Text>
                    <Button
                      onPress={() => handleDownloadArchive(archive.id)}
                      variant="outline"
                      size="sm">
                      {t('admin_get_csv')}
                    </Button>
                  </View>
                ))}
              </View>
            )}
          </CardContent>
        </Card>
      </ScrollView>

      <PinModal
        visible={showPinModal}
        onClose={() => navigation.goBack()}
        festivalCode={festivalCode}
        onSuccess={handlePinSuccess}
      />
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: colors.background,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: spacing.lg,
    paddingTop: spacing.xl,
  },
  backButton: {
    fontSize: fontSize.md,
    fontWeight: fontWeight.semibold,
    color: colors.accent,
  },
  content: {
    flex: 1,
  },
  contentContainer: {
    padding: spacing.lg,
    paddingBottom: spacing.xxl,
  },
  title: {
    fontSize: fontSize.xxl,
    fontWeight: fontWeight.bold,
    color: colors.accent,
    textAlign: 'center',
    marginBottom: spacing.xl,
  },
  card: {
    marginBottom: spacing.lg,
  },
  sectionTitle: {
    fontSize: fontSize.lg,
    fontWeight: fontWeight.semibold,
    color: colors.text,
    marginBottom: spacing.md,
  },
  row: {
    flexDirection: 'row',
    gap: spacing.md,
    alignItems: 'flex-end',
  },
  input: {
    flex: 1,
    marginBottom: 0,
  },
  setButton: {
    minWidth: 100,
  },
  buttonRow: {
    flexDirection: 'row',
    gap: spacing.md,
  },
  actionButton: {
    flex: 1,
  },
  noArchives: {
    fontSize: fontSize.md,
    color: colors.textSecondary,
    textAlign: 'center',
    paddingVertical: spacing.lg,
  },
  archiveList: {
    gap: spacing.sm,
  },
  archiveItem: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: spacing.md,
    backgroundColor: colors.secondary,
    borderRadius: borderRadius.md,
  },
  archiveText: {
    fontSize: fontSize.sm,
    color: colors.text,
  },
});

export default AdminScreen;
