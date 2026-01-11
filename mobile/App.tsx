/**
 * Sumcrowds React Native App
 * Crowd counting application for events and festivals
 */

import React, {useEffect, useState} from 'react';
import {View, Text, StyleSheet, ActivityIndicator} from 'react-native';
import {SafeAreaProvider} from 'react-native-safe-area-context';
import {AppNavigator} from './src/navigation';
import {initI18n} from './src/utils/i18n';
import {auth} from './src/utils/auth';
import {colors} from './src/utils/theme';

function App() {
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    const initialize = async () => {
      try {
        // Initialize i18n
        await initI18n();
        // Initialize authentication
        await auth();
      } catch (error) {
        console.error('Initialization error:', error);
      } finally {
        setIsReady(true);
      }
    };

    initialize();
  }, []);

  if (!isReady) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color={colors.primary} />
        <Text style={styles.loadingText}>Loading...</Text>
      </View>
    );
  }

  return (
    <SafeAreaProvider>
      <AppNavigator />
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: colors.background,
  },
  loadingText: {
    marginTop: 16,
    fontSize: 16,
    color: colors.textSecondary,
  },
});

export default App;
