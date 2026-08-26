import type { FormResponseSummary, FormSubmission, FormSubmissionExportResponse } from '../../forms/types';

type ExportExtension = 'json' | 'pdf' | 'xls';

const pdfPageWidth = 595;
const pdfPageHeight = 842;
const pdfMargin = 48;
const pdfLineHeight = 14;

export function downloadResponsesJSON(payload: FormSubmissionExportResponse) {
  downloadBlob(JSON.stringify(payload, null, 2), 'application/json', responseExportFilename(payload.form.title, 'json'));
}

export function downloadResponsesExcel(form: FormResponseSummary, responses: FormSubmission[]) {
  const headers = ['Enviada em', 'Ciência LGPD', ...form.fields.map((field) => field.label)];
  const rows = responses.map((response) => [
    formatDate(response.submittedAt),
    formatDate(response.privacyAcknowledgedAt),
    ...form.fields.map((field) => formatAnswer(response, field.id)),
  ]);
  const table = [headers, ...rows]
    .map((row) => `<tr>${row.map((cell) => `<td>${escapeHTML(cell)}</td>`).join('')}</tr>`)
    .join('');
  const html = `<!doctype html><html><head><meta charset="utf-8"></head><body><table>${table}</table></body></html>`;

  downloadBlob(html, 'application/vnd.ms-excel;charset=utf-8', responseExportFilename(form.title, 'xls'));
}

export function downloadResponsesPDF(form: FormResponseSummary, responses: FormSubmission[]) {
  const lines = [
    `Respostas - ${form.title}`,
    `Exportado em: ${formatDate(new Date().toISOString())}`,
    '',
    ...responses.flatMap((response, index) => [
      `Resposta ${index + 1}`,
      `Enviada em: ${formatDate(response.submittedAt)}`,
      `Ciência LGPD: ${formatDate(response.privacyAcknowledgedAt)}`,
      ...form.fields.map((field) => `${field.label}: ${formatAnswer(response, field.id)}`),
      '',
    ]),
  ];

  const pdfBytes = createPDF(lines);
  downloadBlob(pdfBytes, 'application/pdf', responseExportFilename(form.title, 'pdf'));
}

export function formatAnswer(response: FormSubmission, fieldId: string) {
  const value = response.answers[fieldId];
  if (value === undefined || value === null || value === '') {
    return '-';
  }

  if (typeof value === 'boolean') {
    return value ? 'Sim' : 'Não';
  }

  return String(value);
}

export function formatDate(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    dateStyle: 'short',
    timeStyle: 'short',
  }).format(new Date(value));
}

function responseExportFilename(title: string, extension: ExportExtension) {
  const slug = title
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)/g, '');

  return `respostas-${slug || 'formulario'}.${extension}`;
}

function downloadBlob(content: BlobPart, type: string, filename: string) {
  const blob = new Blob([content], { type });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
}

function escapeHTML(value: string) {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function createPDF(lines: string[]) {
  const pages = paginatePDFLines(lines.flatMap((line) => wrapPDFLine(line)));
  const objectBodies = new Map<number, string>();
  const pageObjectNumbers = pages.map((_, index) => 4 + index * 2);
  const contentObjectNumbers = pages.map((_, index) => 5 + index * 2);
  const pagesObjectNumber = 2;
  const fontObjectNumber = 3;

  objectBodies.set(1, `<< /Type /Catalog /Pages ${pagesObjectNumber} 0 R >>`);
  objectBodies.set(
    pagesObjectNumber,
    `<< /Type /Pages /Kids [${pageObjectNumbers.map((objectNumber) => `${objectNumber} 0 R`).join(' ')}] /Count ${pages.length} >>`,
  );
  objectBodies.set(fontObjectNumber, '<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>');

  pages.forEach((pageLines, index) => {
    const pageObjectNumber = pageObjectNumbers[index];
    const contentObjectNumber = contentObjectNumbers[index];
    const stream = createPDFContentStream(pageLines);

    objectBodies.set(
      pageObjectNumber,
      `<< /Type /Page /Parent ${pagesObjectNumber} 0 R /MediaBox [0 0 ${pdfPageWidth} ${pdfPageHeight}] /Resources << /Font << /F1 ${fontObjectNumber} 0 R >> >> /Contents ${contentObjectNumber} 0 R >>`,
    );
    objectBodies.set(contentObjectNumber, `<< /Length ${toLatin1Bytes(stream).length} >>\nstream\n${stream}endstream`);
  });

  return buildPDFBytes(objectBodies);
}

function wrapPDFLine(line: string) {
  const maxLength = 92;
  if (line.length <= maxLength) {
    return [line];
  }

  const wrapped: string[] = [];
  let remaining = line;
  while (remaining.length > maxLength) {
    const breakpoint = Math.max(remaining.lastIndexOf(' ', maxLength), 1);
    wrapped.push(remaining.slice(0, breakpoint).trim());
    remaining = remaining.slice(breakpoint).trim();
  }
  wrapped.push(remaining);

  return wrapped;
}

function paginatePDFLines(lines: string[]) {
  const maxLinesPerPage = Math.floor((pdfPageHeight - pdfMargin * 2) / pdfLineHeight);
  const pages: string[][] = [];
  for (let index = 0; index < lines.length; index += maxLinesPerPage) {
    pages.push(lines.slice(index, index + maxLinesPerPage));
  }

  return pages.length > 0 ? pages : [[]];
}

function createPDFContentStream(lines: string[]) {
  return lines
    .map((line, index) => {
      const y = pdfPageHeight - pdfMargin - index * pdfLineHeight;
      const size = index === 0 ? 14 : 10;
      return `BT /F1 ${size} Tf ${pdfMargin} ${y} Td (${escapePDFString(line)}) Tj ET\n`;
    })
    .join('');
}

function buildPDFBytes(objectBodies: Map<number, string>) {
  const maxObjectNumber = Math.max(...objectBodies.keys());
  const bytes: number[] = [];
  const offsets = Array(maxObjectNumber + 1).fill(0);

  pushLatin1(bytes, '%PDF-1.4\n%\xD0\xD0\xD0\xD0\n');
  for (let objectNumber = 1; objectNumber <= maxObjectNumber; objectNumber += 1) {
    const body = objectBodies.get(objectNumber);
    if (!body) {
      continue;
    }

    offsets[objectNumber] = bytes.length;
    pushLatin1(bytes, `${objectNumber} 0 obj\n${body}\nendobj\n`);
  }

  const xrefOffset = bytes.length;
  pushLatin1(bytes, `xref\n0 ${maxObjectNumber + 1}\n0000000000 65535 f \n`);
  for (let objectNumber = 1; objectNumber <= maxObjectNumber; objectNumber += 1) {
    pushLatin1(bytes, `${String(offsets[objectNumber]).padStart(10, '0')} 00000 n \n`);
  }
  pushLatin1(bytes, `trailer\n<< /Size ${maxObjectNumber + 1} /Root 1 0 R >>\nstartxref\n${xrefOffset}\n%%EOF`);

  return new Uint8Array(bytes);
}

function escapePDFString(value: string) {
  return value
    .replace(/\\/g, '\\\\')
    .replace(/\(/g, '\\(')
    .replace(/\)/g, '\\)')
    .replace(/\r?\n/g, ' ');
}

function toLatin1Bytes(value: string) {
  return Array.from(value, (character) => {
    const code = character.charCodeAt(0);
    return code <= 255 ? code : 63;
  });
}

function pushLatin1(bytes: number[], value: string) {
  bytes.push(...toLatin1Bytes(value));
}
