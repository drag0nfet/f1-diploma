export function getUnMdText(text) {
    if (!text) return "";

    const firstLine = text.split('\n')[0];

    return firstLine
        .replace(/^#+\s*/, '')
        .replace(/[_*~`>]/g, '')
        .trim();
}