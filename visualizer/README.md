# lem-in 3D visualizer

Run `lem-in`, copy its complete output, then open `index.html` in a browser and
paste the output into the text area. Press **Load colony** to render it.

The page is static and needs no local server. It loads Three.js and OrbitControls
from a CDN, so the first browser load requires an internet connection.

Run the parser and playback tests from the repository root with:

```sh
node --test visualizer/app.test.js
```
