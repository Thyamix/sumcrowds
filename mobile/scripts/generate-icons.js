const sharp = require('sharp');
const fs = require('fs');
const path = require('path');

// SVG content from the website's plus.svg
const svgContent = `<svg width="1000" height="1000" viewBox="0 0 1000 1000" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <mask id="plus">
      <rect x="0" y="0" width="1000" height="1000" rx="0" fill="white"/>
      <rect x="146" y="392" width="709" height="216" rx="120" fill="black"/>
      <rect x="392" y="146" width="216" height="709" rx="120" fill="black"/>
    </mask>
  </defs>
  <circle cx="500" cy="500" r="450" fill="#eef3fe"/>
  <circle cx="500" cy="500" r="500" fill="#514fee" mask="url(#plus)"/>
</svg>`;

// Android icon sizes
const androidIcons = [
  { size: 48, folder: 'mipmap-mdpi' },
  { size: 72, folder: 'mipmap-hdpi' },
  { size: 96, folder: 'mipmap-xhdpi' },
  { size: 144, folder: 'mipmap-xxhdpi' },
  { size: 192, folder: 'mipmap-xxxhdpi' },
];

const androidResPath = path.join(__dirname, '..', 'android', 'app', 'src', 'main', 'res');

async function generateIcons() {
  console.log('Generating Android app icons...');

  for (const icon of androidIcons) {
    const outputDir = path.join(androidResPath, icon.folder);

    // Generate regular icon
    await sharp(Buffer.from(svgContent))
      .resize(icon.size, icon.size)
      .png()
      .toFile(path.join(outputDir, 'ic_launcher.png'));

    // Generate round icon
    await sharp(Buffer.from(svgContent))
      .resize(icon.size, icon.size)
      .png()
      .toFile(path.join(outputDir, 'ic_launcher_round.png'));

    console.log(`Generated ${icon.size}x${icon.size} icons in ${icon.folder}`);
  }

  console.log('Done!');
}

generateIcons().catch(console.error);
