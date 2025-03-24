function localizeDate(elementId, timestamp) {
    const element = document.getElementById(elementId);
    if (!element) {
        console.warn(`Element with ID "${elementId}" not found.`);
        return;
    }
    const date = new Date(timestamp);
    element.innerText = date.toLocaleDateString(undefined, {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
    });
}

window.localizeDate = localizeDate;