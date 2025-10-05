const layerOrder = ['head', 'hair', 'eyes', 'nose', 'mouth', 'accessory'];

let pal = {
    base: [255, 223, 186],
    head: [null, [255, 224, 189]],
    hair: [null, [0, 0, 0]],
    nose: [null, [0, 0, 0]],
    eyes: [null, [0, 0, 0]],
    mouth: [null, [255, 0, 0]],
    accessory: [null, [0, 0, 0]]
};
if (typeof pal_json !== 'undefined' && pal_json) {
    try {
        const loaded = typeof pal_json === 'string' ? JSON.parse(pal_json) : pal_json;
        if (loaded && typeof loaded === 'object') {
            pal = Object.assign(pal, loaded);
        }
    } catch (e) {}
}

function rgbToHex(rgb) {
    return '#' + rgb.map(x => x.toString(16).padStart(2, '0')).join('');
}
function hexToRgb(hex) {
    hex = hex.replace('#', '');
    return [0, 1, 2].map(i => parseInt(hex.substr(i * 2, 2), 16));
}

function updatePreview() {
    // For now, just show the JSON. Later it will request a preview from the backend.
    document.getElementById('pal-json').textContent = JSON.stringify(pal, null, 2);
}

function setPart(layer, filename) {
    pal[layer][0] = filename;
    updatePreview();
}
function setColor(layer, color) {
    pal[layer][1] = hexToRgb(color);
    updatePreview();
}
function setBaseColor(color) {
    pal.base = hexToRgb(color);
    updatePreview();
}

function savePal() {
    fetch('/api/save_pal', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pal })
    }).then(r => r.json()).then(data => {
        if (data.success) alert('Pal saved!');
        else alert('Error saving pal: ' + (data.error || 'unknown'));
    });
}

document.addEventListener('DOMContentLoaded', function() {
    for (const layer of layerOrder) {
        const select = document.getElementById(layer + '-select');
        if (!select) continue;
        select.innerHTML = '<option value="">None</option>' +
            (pal_parts[layer] || []).map(f => `<option value="${f.replace('.svg','')}">${f.replace('.svg','')}</option>`).join('');
        if (pal[layer] && pal[layer][0]) {
            select.value = pal[layer][0];
        }
        select.addEventListener('change', e => setPart(layer, e.target.value));
    }
    for (const layer of layerOrder) {
        const colorInput = document.getElementById(layer + '-color');
        if (!colorInput) continue;
        colorInput.value = rgbToHex(pal[layer][1]);
        colorInput.addEventListener('input', e => setColor(layer, e.target.value));
    }
    const baseColorInput = document.getElementById('base-color');
    if (baseColorInput) {
        baseColorInput.value = rgbToHex(pal.base);
        baseColorInput.addEventListener('input', e => setBaseColor(e.target.value));
    }
    document.getElementById('save-pal-btn').addEventListener('click', savePal);
    updatePreview();
});
