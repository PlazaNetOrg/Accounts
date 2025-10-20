import os
import json
import re


def compose_pal(pal_json, parts_dir, size=(300, 300)):
    if isinstance(pal_json, str):
        try: pal = json.loads(pal_json)
        except: pal = {}
    else: pal = pal_json or {}
    w, h = size
    bg = pal.get("base", [255,255,255])
    bg_hex = '#%02x%02x%02x' % tuple(bg)
    svg = [f'<svg xmlns="http://www.w3.org/2000/svg" width="{w}" height="{h}" viewBox="0 0 {w} {h}">']
    svg.append(f'<circle cx="{w//2}" cy="{h//2}" r="{min(w,h)//2}" fill="{bg_hex}"/>')
    layers = ['head', 'hair', 'eyebrows', 'eyes', 'nose', 'mustache', 'mouth', 'beard', 'accessory']
    from flask import current_app
    static = current_app.static_folder
    parts_dir = os.path.join(static, os.path.relpath(parts_dir, 'static'))
    for layer in layers:
        part = pal.get(layer)
        if not part or not part[0]: continue
        fn = part[0]
        color = part[1] if len(part)>1 and isinstance(part[1], list) else None
        path = os.path.join(parts_dir, layer, f"{fn}.svg")
        if not os.path.exists(path): continue
        s = open(path, encoding='utf-8').read()
        s = re.sub(r'<\?xml[^>]*>\s*', '', s, flags=re.I)
        s = re.sub(r'<svg[^>]*>', '', s, count=1, flags=re.I)
        s = re.sub(r'</svg>', '', s, flags=re.I)
        if color:
            hex_color = '#%02x%02x%02x' % tuple(color)
            is_white = (re.search(r'(fill|stroke)(\s*:\s*)#fff(fff)?', s, flags=re.I))
            if not is_white:
                s = re.sub(r'(fill|stroke)=("|\')#[0-9a-fA-F]{3,6}("|\')', lambda m: f'{m.group(1)}="{hex_color}"', s, flags=re.I)
                def style_color(m):
                    style = m.group(0)
                    style = re.sub(r'(fill\s*:\s*)#[0-9a-fA-F]{3,6}', r'\1'+hex_color, style, flags=re.I)
                    style = re.sub(r'(stroke\s*:\s*)#[0-9a-fA-F]{3,6}', r'\1'+hex_color, style, flags=re.I)
                    return style
                s = re.sub(r'style=("[^"]*"|\'[^"]*\')', style_color, s, flags=re.I)
        svg.append(f'<g transform="translate({w*0.125},{h*0.125}) scale(0.75)">{s}</g>')
    svg.append('</svg>')
    return '\n'.join(svg)