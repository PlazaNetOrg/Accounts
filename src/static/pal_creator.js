const layerOrder = ['head', 'hair', 'eyebrows', 'eyes', 'nose', 'mustache', 'mouth', 'beard', 'accessory'];

let pal = {
    base: [58, 159, 241],
    head: ['Teardrop', [255, 224, 189]],
    hair: [null, [0, 0, 0]],
    eyebrows: [null, [0, 0, 0]],
    eyes: ['Round', [0, 0, 0]],
    nose: ['Curve', [0, 0, 0]],
    mustache: [null, [0, 0, 0]],
    mouth: ['Neutral', [0, 0, 0]],
    beard: [null, [0, 0, 0]],
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
    // Show JSON
    document.getElementById('pal-json').textContent = JSON.stringify(pal, null, 2);
    // Show SVG preview
    fetch('/api/pal_preview', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pal })
    })
    .then(r => r.ok ? r.text() : Promise.reject('Preview error'))
    .then(svg => {
        let preview = document.getElementById('pal-svg-preview');
        if (!preview) {
            preview = document.createElement('div');
            preview.id = 'pal-svg-preview';
            preview.style.width = '256px';
            preview.style.height = '256px';
            preview.style.display = 'flex';
            preview.style.alignItems = 'center';
            preview.style.justifyContent = 'center';
            document.querySelector('.pal-preview').prepend(preview);
        }
        preview.innerHTML = svg;
    })
    .catch(() => {
        let preview = document.getElementById('pal-svg-preview');
        if (preview) preview.innerHTML = '<span style="color:red">Preview error</span>';
    });
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
    console.info(JSON.stringify({ pal }));
    fetch('/api/save_pal', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pal }),
        credentials: 'include'
    }).then(r => r.json()).then(data => {
        if (data.success) alert('Pal saved!');
        else alert('Error saving pal: ' + (data.error || 'unknown'));
    });
}

document.addEventListener('DOMContentLoaded', function() {
    // Tab switching
    const tabBtns = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-content');
    tabBtns.forEach(btn => {
        btn.addEventListener('click', function() {
            tabBtns.forEach(b => b.classList.remove('active'));
            tabContents.forEach(tc => tc.style.display = 'none');
            btn.classList.add('active');
            document.getElementById('tab-' + btn.dataset.tab).style.display = 'block';
        });
    });
    if (tabBtns.length) {
        tabBtns[0].classList.add('active');
        tabContents[0].style.display = 'block';
    }

    // Populate part pickers
    for (const layer of layerOrder) {
        const previews = document.getElementById(layer + '-previews');
        if (previews) {
            let noneBtn = document.createElement('div');
            noneBtn.className = 'part-preview' + (!pal[layer][0] ? ' selected' : '');
            noneBtn.title = 'None';
            noneBtn.innerHTML = `<svg width="40" height="40" fill="none" xmlns="http://www.w3.org/2000/svg"></svg>`;
            noneBtn.addEventListener('click', function() {
                setPart(layer, '');
                previews.querySelectorAll('.part-preview').forEach(d => d.classList.remove('selected'));
                noneBtn.classList.add('selected');
            });
            previews.innerHTML = '';
            previews.appendChild(noneBtn);

            (pal_parts[layer] || []).forEach(f => {
                let div = document.createElement('div');
                div.className = 'part-preview' + (pal[layer] && pal[layer][0] === f.replace('.svg','') ? ' selected' : '');
                div.dataset.fn = f.replace('.svg','');
                let img = document.createElement('img');
                img.src = `/static/assets/pal_parts/${layer}/${f}`;
                img.alt = f;
                img.style.borderRadius = '8px';
                div.appendChild(img);
                div.addEventListener('click', function() {
                    setPart(layer, div.dataset.fn);
                    previews.querySelectorAll('.part-preview').forEach(d => d.classList.remove('selected'));
                    div.classList.add('selected');
                });
                previews.appendChild(div);
            });
        }
    }
    // Color pickers
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
