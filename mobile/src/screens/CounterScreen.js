import React, {useState, useEffect, useRef, useCallback} from 'react';
import {
  View,
  Text,
  StyleSheet,
  StatusBar,
  Platform,
  TouchableOpacity,
} from 'react-native';
import {useTranslation} from 'react-i18next';
import {SafeAreaView} from 'react-native-safe-area-context';
import {Button} from '../components/ui';
import {
  LanguageSwitcher,
  PasswordModal,
  LeaveConfirmModal,
} from '../components';
import {fetchWithAuth, auth} from '../utils/auth';
import {WS_URL} from '../config';
import {colors, spacing, fontSize, fontWeight, borderRadius} from '../utils/theme';

const STATUSBAR_HEIGHT = Platform.OS === 'android' ? StatusBar.currentHeight || 24 : 0;

export const CounterScreen = ({route, navigation}) => {
  const {festivalCode} = route.params;
  const {t} = useTranslation();

  const [total, setTotal] = useState('...');
  const [maxJauge, setMaxJauge] = useState(0);
  const [status, setStatus] = useState('normal');
  const [isConnected, setIsConnected] = useState(false);
  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [showLeaveModal, setShowLeaveModal] = useState(false);
  const [hasAccess, setHasAccess] = useState(false);

  const socketRef = useRef(null);
  const heartbeatRef = useRef(null);

  // Calculate status based on total and max
  const calculateStatus = useCallback((currentTotal, max) => {
    if (max === 0) return 'normal';
    const percentage = (currentTotal / max) * 100;
    if (percentage >= 100) return 'danger';
    if (percentage >= 75) return 'warning';
    return 'normal';
  }, []);

  // Get status color
  const getStatusColor = () => {
    switch (status) {
      case 'danger':
        return colors.statusDanger;
      case 'warning':
        return colors.statusWarning;
      default:
        return colors.statusNormal;
    }
  };

  // WebSocket connection
  const connectWebSocket = useCallback(() => {
    if (socketRef.current) {
      socketRef.current.close();
    }
    if (heartbeatRef.current) {
      clearInterval(heartbeatRef.current);
    }

    const ws = new WebSocket(WS_URL + festivalCode);

    ws.onopen = () => {
      setIsConnected(true);
      ws.send(JSON.stringify({type: 'getTotal'}));

      // Start heartbeat
      heartbeatRef.current = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify({type: 'ping'}));
        }
      }, 10000);
    };

    ws.onmessage = event => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'pong') return;

        if (data.total !== undefined) {
          setTotal(data.total);
          if (data.jauge !== undefined) {
            setMaxJauge(data.jauge);
            setStatus(calculateStatus(data.total, data.jauge));
          }
        }
      } catch (err) {
        console.error('WebSocket message error:', err);
      }
    };

    ws.onclose = () => {
      setIsConnected(false);
      if (heartbeatRef.current) {
        clearInterval(heartbeatRef.current);
      }
    };

    ws.onerror = error => {
      console.error('WebSocket error:', error);
    };

    socketRef.current = ws;
  }, [festivalCode, calculateStatus]);

  // Check access on mount
  useEffect(() => {
    const checkAccess = async () => {
      await auth();
      try {
        const response = await fetchWithAuth(`v1/festival/${festivalCode}/access`);
        if (response.status === 403) {
          setShowPasswordModal(true);
        } else if (response.ok) {
          setHasAccess(true);
        } else {
          navigation.replace('Home');
        }
      } catch (err) {
        console.error('Access check error:', err);
        navigation.replace('Home');
      }
    };

    checkAccess();
  }, [festivalCode, navigation]);

  // Connect WebSocket when access is granted
  useEffect(() => {
    if (hasAccess) {
      connectWebSocket();
    }

    return () => {
      if (socketRef.current) {
        socketRef.current.close();
      }
      if (heartbeatRef.current) {
        clearInterval(heartbeatRef.current);
      }
    };
  }, [hasAccess, connectWebSocket]);

  // Increment/Decrement handlers
  const handleIncrement = async amount => {
    try {
      await fetchWithAuth(`v1/festival/${festivalCode}/inc`, {
        method: 'POST',
        body: JSON.stringify({amount}),
      });
    } catch (err) {
      console.error('Increment error:', err);
    }
  };

  const handleDecrement = async amount => {
    try {
      await fetchWithAuth(`v1/festival/${festivalCode}/dec`, {
        method: 'POST',
        body: JSON.stringify({amount}),
      });
    } catch (err) {
      console.error('Decrement error:', err);
    }
  };

  const handlePasswordSuccess = () => {
    setShowPasswordModal(false);
    setHasAccess(true);
  };

  const handleLeave = () => {
    navigation.replace('Home');
  };

  const handleAdminPress = () => {
    navigation.navigate('Admin', {festivalCode});
  };

  return (
    <SafeAreaView style={[styles.container, {backgroundColor: getStatusColor()}]} edges={['left', 'right', 'bottom']}>
      <StatusBar barStyle="light-content" backgroundColor={getStatusColor()} translucent />

      <View style={styles.header}>
        <TouchableOpacity onPress={() => setShowLeaveModal(true)}>
          <Text style={styles.homeButton}>{t('home_home')}</Text>
        </TouchableOpacity>

        <View style={styles.headerRight}>
          <TouchableOpacity onPress={handleAdminPress} style={styles.settingsButton}>
            <Text style={styles.settingsIcon}>⚙️</Text>
          </TouchableOpacity>
          <LanguageSwitcher />
        </View>
      </View>

      <View style={styles.codeContainer}>
        <Text style={styles.codeLabel}>{t('counter_code')}</Text>
        <Text style={styles.codeValue}>{festivalCode}</Text>
      </View>

      <View style={styles.counterContainer}>
        {!isConnected ? (
          <View style={styles.disconnectedContainer}>
            <Text style={styles.disconnectedText}>{t('counter_disconnected')}</Text>
            <Button
              onPress={connectWebSocket}
              variant="outline"
              style={styles.reconnectButton}
              textStyle={styles.reconnectButtonText}>
              {t('counter_reconnect')}
            </Button>
          </View>
        ) : (
          <>
            <Text style={styles.counterValue}>{total}</Text>
            {maxJauge > 0 && (
              <View style={styles.gaugeContainer}>
                <Text style={styles.gaugeLabel}>{t('counter_gauge')}</Text>
                <Text style={styles.gaugeValue}>{maxJauge}</Text>
              </View>
            )}
          </>
        )}
      </View>

      <View style={styles.controls}>
        {/* Exit column */}
        <View style={styles.controlColumn}>
          <Text style={[styles.controlLabel, {color: colors.destructive}]}>{t('counter_exit')}</Text>
          <View style={styles.smallButtonRow}>
            <Button
              onPress={() => handleDecrement(2)}
              variant="destructive"
              size="counter"
              style={styles.smallButton}>
              -2
            </Button>
            <Button
              onPress={() => handleDecrement(3)}
              variant="destructive"
              size="counter"
              style={styles.smallButton}>
              -3
            </Button>
          </View>
          <Button
            onPress={() => handleDecrement(1)}
            variant="destructive"
            size="counterLg"
            style={styles.largeButton}>
            -1
          </Button>
        </View>

        {/* Enter column */}
        <View style={styles.controlColumn}>
          <Text style={[styles.controlLabel, {color: colors.success}]}>{t('counter_enter')}</Text>
          <View style={styles.smallButtonRow}>
            <Button
              onPress={() => handleIncrement(2)}
              variant="success"
              size="counter"
              style={styles.smallButton}>
              +2
            </Button>
            <Button
              onPress={() => handleIncrement(3)}
              variant="success"
              size="counter"
              style={styles.smallButton}>
              +3
            </Button>
          </View>
          <Button
            onPress={() => handleIncrement(1)}
            variant="success"
            size="counterLg"
            style={styles.largeButton}>
            +1
          </Button>
        </View>
      </View>

      <PasswordModal
        visible={showPasswordModal}
        onClose={() => navigation.replace('Home')}
        festivalCode={festivalCode}
        onSuccess={handlePasswordSuccess}
      />

      <LeaveConfirmModal
        visible={showLeaveModal}
        onClose={() => setShowLeaveModal(false)}
        onConfirm={handleLeave}
      />
    </SafeAreaView>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    padding: spacing.lg,
    paddingTop: STATUSBAR_HEIGHT + spacing.md,
  },
  homeButton: {
    fontSize: fontSize.lg,
    fontWeight: fontWeight.semibold,
    color: colors.white,
  },
  headerRight: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: spacing.md,
  },
  settingsButton: {
    padding: spacing.sm,
  },
  settingsIcon: {
    fontSize: fontSize.xl,
  },
  codeContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: spacing.sm,
    paddingVertical: spacing.md,
  },
  codeLabel: {
    fontSize: fontSize.md,
    color: colors.white,
    opacity: 0.8,
  },
  codeValue: {
    fontSize: fontSize.lg,
    fontWeight: fontWeight.bold,
    color: colors.white,
    backgroundColor: 'rgba(255,255,255,0.2)',
    paddingHorizontal: spacing.md,
    paddingVertical: spacing.xs,
    borderRadius: borderRadius.md,
  },
  counterContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  disconnectedContainer: {
    alignItems: 'center',
    gap: spacing.md,
  },
  disconnectedText: {
    fontSize: fontSize.xl,
    color: colors.white,
    fontWeight: fontWeight.semibold,
  },
  reconnectButton: {
    borderColor: colors.white,
  },
  reconnectButtonText: {
    color: colors.white,
  },
  counterValue: {
    fontSize: fontSize.massive,
    fontWeight: fontWeight.bold,
    color: colors.white,
    fontFamily: 'monospace',
  },
  gaugeContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: spacing.md,
    gap: spacing.sm,
  },
  gaugeLabel: {
    fontSize: fontSize.md,
    color: colors.white,
    opacity: 0.8,
  },
  gaugeValue: {
    fontSize: fontSize.xl,
    fontWeight: fontWeight.bold,
    color: colors.white,
  },
  controls: {
    flexDirection: 'row',
    padding: spacing.lg,
    paddingBottom: spacing.xxl,
    gap: spacing.md,
  },
  controlColumn: {
    flex: 1,
    alignItems: 'center',
    gap: spacing.sm,
    padding: spacing.md,
    borderRadius: borderRadius.lg,
  },
  controlLabel: {
    fontSize: fontSize.sm,
    fontWeight: fontWeight.semibold,
    textTransform: 'uppercase',
    letterSpacing: 1,
  },
  smallButtonRow: {
    flexDirection: 'row',
    gap: spacing.sm,
  },
  smallButton: {
    minWidth: 60,
    minHeight: 50,
  },
  largeButton: {
    minWidth: 130,
    minHeight: 60,
  },
});

export default CounterScreen;
