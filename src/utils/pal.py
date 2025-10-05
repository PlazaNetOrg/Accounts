from PIL import Image
import os
import json

def compose_pal(pal_json, parts_dir, size=(256, 256)):
    if isinstance(pal_json, str):
        try:
            pal = json.loads(pal_json)
        except Exception:
            pal = {}
    else:
        pal = pal_json or {}

    background_color = tuple(pal.get("base", [255, 255, 255]))

    base = Image.new('RGB', size, background_color)

    layer_order = ['face', 'hair', 'eyes', 'mouth', 'accessory']
    for layer in layer_order:
        choice = pal.get(layer)
        if not choice:
            continue
        path = os.path.join(parts_dir, layer, f"{choice}.png")
        if os.path.exists(path):
            try:
                img = Image.open(path).convert('RGBA')
                img = img.resize(size, Image.NEAREST)
                base.paste(img, (0, 0), img)
            except Exception:
                continue
    return base