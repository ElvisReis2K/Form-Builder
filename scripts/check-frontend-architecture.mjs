import { readdirSync, readFileSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';

const root = process.cwd().endsWith('frontend') ? join(process.cwd(), '..') : process.cwd();
const srcDir = join(root, 'frontend', 'src');
const violations = [];

function walk(dir) {
  for (const entry of readdirSync(dir)) {
    const fullPath = join(dir, entry);
    const stats = statSync(fullPath);

    if (stats.isDirectory()) {
      walk(fullPath);
      continue;
    }

    if (!entry.endsWith('.tsx')) {
      continue;
    }

    const content = readFileSync(fullPath, 'utf8');
    const checks = [
      { pattern: /sx=\{\{/g, message: 'Move inline MUI sx objects to a *.styles.ts file.' },
      { pattern: /style=\{\{/g, message: 'Move inline style objects to a *.styles.ts file.' },
    ];

    for (const check of checks) {
      for (const match of content.matchAll(check.pattern)) {
        const line = content.slice(0, match.index).split('\n').length;
        violations.push(`${relative(root, fullPath)}:${line} - ${check.message}`);
      }
    }
  }
}

walk(srcDir);

if (violations.length > 0) {
  console.error('Frontend architecture violations found:\n');
  console.error(violations.join('\n'));
  process.exit(1);
}

console.log('frontend architecture ok');
