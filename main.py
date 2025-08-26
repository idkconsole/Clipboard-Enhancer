import pyperclip as pc, keyboard as kb, google.generativeai as g, time

g.configure(api_key="gemini api key")

def glow(x):
    r = g.GenerativeModel("gemini-1.5-flash").generate_content(
        f"fix grammar, spelling, make clear, smooth, pro, keep meaning:\n{x}"
    )
    return r.text.strip()

def go():
    t = pc.paste().strip()
    if not t: return print("empty")
    print(f"{t}\n...wait")
    e = glow(t)
    pc.copy(e)
    print(e)

print("run")
kb.add_hotkey('delete+end', go)
while not kb.is_pressed('esc'): time.sleep(.1)
