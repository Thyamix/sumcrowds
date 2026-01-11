# Sumcrowds Mobile App

React Native Android app for Sumcrowds - a crowd counting application for events and festivals.

## Features

- **Join Session**: Enter a 6-character festival code to join an existing session
- **Create Session**: Create a new festival with an admin PIN and optional password
- **Real-time Counter**: WebSocket-based live updates showing current crowd count
- **Visual Status**: Color-coded capacity indicators (green/orange/red)
- **Admin Panel**: Set capacity, archive events, and download CSV reports
- **Internationalization**: English and French language support
- **JWT Authentication**: Secure token-based authentication

## Prerequisites

- Node.js >= 20
- Android Studio with Android SDK
- Java Development Kit (JDK) 17+
- An Android device or emulator

## Installation

```bash
# Navigate to the mobile directory
cd mobile

# Install dependencies
npm install
```

## Configuration

Update the API and WebSocket URLs in `src/config.js`:

```javascript
// For production
export const API_URL = 'https://api.sumcrowds.com/';
export const WS_URL = 'wss://ws.sumcrowds.com/';

// For development with Android emulator
export const API_URL = 'http://10.0.2.2:8080/';
export const WS_URL = 'ws://10.0.2.2:8080/ws/';
```

Note: `10.0.2.2` is the special IP that Android emulator uses to access the host machine's localhost.

## Running the App

### Development

```bash
# Start Metro bundler
npm start

# In another terminal, run on Android
npm run android
```

### Building a Release APK

```bash
cd android
./gradlew assembleRelease
```

The APK will be at `android/app/build/outputs/apk/release/app-release.apk`

## Project Structure

```
mobile/
├── src/
│   ├── screens/           # Main app screens
│   │   ├── HomeScreen.js      # Landing page
│   │   ├── CounterScreen.js   # Real-time counter
│   │   └── AdminScreen.js     # Admin panel
│   ├── components/        # Reusable components
│   │   ├── ui/               # Base UI components
│   │   ├── JoinModal.js      # Join session modal
│   │   ├── CreateModal.js    # Create session modal
│   │   ├── PasswordModal.js  # Password entry modal
│   │   ├── PinModal.js       # Admin PIN modal
│   │   └── LanguageSwitcher.js
│   ├── navigation/        # React Navigation setup
│   ├── utils/             # Utilities
│   │   ├── auth.js           # JWT authentication
│   │   ├── i18n.js           # Internationalization
│   │   └── theme.js          # Colors and styling
│   ├── locales/           # Translation files
│   │   ├── en/
│   │   └── fr/
│   └── config.js          # API configuration
├── android/               # Android native code
└── App.tsx               # App entry point
```

## Authentication

The app uses JWT-based authentication instead of cookies (which are used in the web version). Tokens are stored securely using AsyncStorage.

### Token Flow

1. On app start, attempt to refresh access token
2. If refresh fails, initialize new access tokens
3. All API requests include `Authorization: Bearer <token>` header
4. On 401 response, automatically refresh and retry

## Troubleshooting

### Metro bundler issues
```bash
npm start -- --reset-cache
```

### Android build issues
```bash
cd android
./gradlew clean
cd ..
npm run android
```

### WebSocket not connecting
- Ensure the WebSocket URL is correct in `src/config.js`
- For local development, use `10.0.2.2` instead of `localhost`
- Check that the backend WebSocket server is running
