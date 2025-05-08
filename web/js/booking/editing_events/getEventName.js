export function getEventName(description) {
    if (!description) return "";

    const firstLine = description.split('\n')[0];

    return firstLine
        .replace(/^#+\s*/, '')
        .replace(/[_*~`>]/g, '')
        .trim();
}