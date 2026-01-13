# Sumcrowds

## What is Sumcrowds?

[Sumcrowds](https://sumcrowds.com) is a real-time event attendance counter. It was designed to help the volunteers of a local music festival keep track of attendance and ensure it did not exceed the legal capacity of the venue. Unfortunately, the festival ended up getting canceled.

## How does it work?

The app uses a session code to sync users to events. When a session is created, a random 6-character code is generated, which others can use to join. Each session is password protected, with an additional PIN required for access to the admin panel.

## Features

- Real-time tracking of festival attendance
- Unique 6-character session codes to join events
- Password-protected sessions for security
- Additional PIN protection for admin panel access
- Easy-to-use interface for volunteers and organizers
- Supports multiple simultaneous sessions
- Recent sessions history for quick access to previously joined events
- Exports data to csv
- Support for both English and French

## Configuration

The project uses a centralized configuration system with separate files for each environment.

### Config Files

Non-secret configuration (endpoints, CORS, ports) is stored in TOML files:

- `config.dev.toml` - Development environment
- `config.staging.toml` - Staging environment
- `config.prod.toml` - Production environment

### Environment Files

Secrets (database passwords, API keys) are stored in `.env` files:

- `.env.dev` - Development secrets
- `.env.staging` - Staging secrets
- `.env.prod` - Production secrets

Copy `.env.example` and fill in the required values for your environment.

### Mobile Configuration

The mobile app generates its config at build time:

```bash
cd mobile
npm run generate-config:dev   # For development
npm run generate-config:prod  # For production
```

This reads from the root config files and generates `mobile/src/config.ts`.
