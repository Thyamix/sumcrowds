import React, {useState} from 'react';
import {View, Text, StyleSheet, StatusBar, Platform} from 'react-native';
import {useTranslation} from 'react-i18next';
import {SafeAreaView} from 'react-native-safe-area-context';
import type {NativeStackNavigationProp} from '@react-navigation/native-stack';
import {Button, Card, CardContent} from '../components/ui';
import {LanguageSwitcher, JoinModal, CreateModal} from '../components';
import {colors, spacing, fontSize, fontWeight} from '../utils/theme';
import type {RootStackParamList} from '../navigation';

const STATUSBAR_HEIGHT: number = Platform.OS === 'android' ? StatusBar.currentHeight || 24 : 0;

interface HomeScreenProps {
  navigation: NativeStackNavigationProp<RootStackParamList, 'Home'>;
}

export const HomeScreen: React.FC<HomeScreenProps> = ({navigation}) => {
  const {t} = useTranslation();
  const [joinOpen, setJoinOpen] = useState(false);
  const [createOpen, setCreateOpen] = useState(false);

  const handleJoin = (code: string): void => {
    navigation.navigate('Counter', {festivalCode: code});
  };

  const handleCreated = (code: string): void => {
    navigation.navigate('Counter', {festivalCode: code});
  };

  return (
    <SafeAreaView style={styles.container} edges={['left', 'right', 'bottom']}>
      <StatusBar barStyle="dark-content" backgroundColor={colors.background} />

      <View style={styles.header}>
        <Text style={styles.headerTitle}>{t('home_home')}</Text>
        <LanguageSwitcher />
      </View>

      <View style={styles.content}>
        <Card style={styles.card}>
          <CardContent>
            <Text style={styles.welcome}>{t('home_welcome')}</Text>
            <Text style={styles.subtitle}>{t('home_select_option')}</Text>

            <View style={styles.buttons}>
              <Button
                onPress={() => setJoinOpen(true)}
                size="lg"
                variant="default"
                style={styles.button}>
                {t('home_join_button')}
              </Button>

              <Button
                onPress={() => setCreateOpen(true)}
                size="lg"
                variant="success"
                style={styles.button}>
                {t('home_create_button')}
              </Button>
            </View>
          </CardContent>
        </Card>
      </View>

      <JoinModal
        visible={joinOpen}
        onClose={() => setJoinOpen(false)}
        onJoin={handleJoin}
      />

      <CreateModal
        visible={createOpen}
        onClose={() => setCreateOpen(false)}
        onCreated={handleCreated}
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
    paddingTop: STATUSBAR_HEIGHT + spacing.md,
  },
  headerTitle: {
    fontSize: fontSize.xl,
    fontWeight: fontWeight.bold,
    color: colors.text,
  },
  content: {
    flex: 1,
    justifyContent: 'center',
    padding: spacing.lg,
  },
  card: {
    marginHorizontal: spacing.md,
  },
  welcome: {
    fontSize: fontSize.xxl,
    fontWeight: fontWeight.bold,
    color: colors.primary,
    textAlign: 'center',
    marginBottom: spacing.sm,
  },
  subtitle: {
    fontSize: fontSize.md,
    color: colors.textSecondary,
    textAlign: 'center',
    marginBottom: spacing.xl,
  },
  buttons: {
    gap: spacing.md,
  },
  button: {
    marginBottom: spacing.sm,
  },
});

export default HomeScreen;
