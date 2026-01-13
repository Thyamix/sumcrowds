#!/usr/bin/env node

/**
 * Mobile Config Generator
 *
 * Reads the root config.{env}.toml file and generates mobile/src/config.ts
 *
 * Usage:
 *   node scripts/generate-mobile-config.js [env]
 *
 *   env: dev, staging, or prod (default: prod)
 */

const fs = require('fs');
const path = require('path');

// Simple TOML parser for our needs (avoiding external dependencies)
function parseTOML(content) {
    const result = {};
    let currentSection = null;

    const lines = content.split('\n');
    for (const line of lines) {
        const trimmed = line.trim();

        // Skip comments and empty lines
        if (trimmed.startsWith('#') || trimmed === '') {
            continue;
        }

        // Section header
        const sectionMatch = trimmed.match(/^\[([^\]]+)\]$/);
        if (sectionMatch) {
            currentSection = sectionMatch[1];
            result[currentSection] = {};
            continue;
        }

        // Key-value pair
        const kvMatch = trimmed.match(/^([^=]+)=(.*)$/);
        if (kvMatch && currentSection) {
            const key = kvMatch[1].trim();
            let value = kvMatch[2].trim();

            // Parse value
            if (value.startsWith('"') && value.endsWith('"')) {
                // String
                value = value.slice(1, -1);
            } else if (value.startsWith('[') && value.endsWith(']')) {
                // Array
                value = JSON.parse(value.replace(/'/g, '"'));
            } else if (value === 'true') {
                value = true;
            } else if (value === 'false') {
                value = false;
            } else if (!isNaN(value)) {
                value = Number(value);
            }

            result[currentSection][key] = value;
        }
    }

    return result;
}

function main() {
    const env = process.argv[2] || 'prod';
    const validEnvs = ['dev', 'staging', 'prod'];

    if (!validEnvs.includes(env)) {
        console.error(`Error: Invalid environment "${env}". Must be one of: ${validEnvs.join(', ')}`);
        process.exit(1);
    }

    // Find project root (where this script is in scripts/)
    const scriptDir = __dirname;
    const projectRoot = path.dirname(scriptDir);

    const configPath = path.join(projectRoot, `config.${env}.toml`);
    const outputPath = path.join(projectRoot, 'mobile', 'src', 'config.ts');

    // Read and parse config
    if (!fs.existsSync(configPath)) {
        console.error(`Error: Config file not found: ${configPath}`);
        process.exit(1);
    }

    const configContent = fs.readFileSync(configPath, 'utf8');
    const config = parseTOML(configContent);

    // Extract endpoints
    const apiUrl = config.endpoints?.api_base || '';
    const wsUrl = config.endpoints?.ws_base || '';

    if (!apiUrl || !wsUrl) {
        console.error('Error: Missing api_base or ws_base in config');
        process.exit(1);
    }

    // Generate TypeScript config
    const tsContent = `// Auto-generated from config.${env}.toml
// Do not edit manually - run: node scripts/generate-mobile-config.js ${env}

// API Configuration
export const API_URL: string = '${apiUrl}';
export const WS_URL: string = '${wsUrl}';

// Environment helper
declare const __DEV__: boolean;
export const isDev: boolean = __DEV__;

// Config metadata
export const CONFIG_ENV: string = '${env}';
`;

    // Write output
    fs.writeFileSync(outputPath, tsContent);
    console.log(`Generated ${outputPath} from config.${env}.toml`);
    console.log(`  API_URL: ${apiUrl}`);
    console.log(`  WS_URL: ${wsUrl}`);
}

main();
