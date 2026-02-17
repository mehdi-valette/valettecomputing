class BinaryGrid extends HTMLElement {
  #canvas;
  #ctx;

  /** @type {Painter | null} */
  painter = null;

  constructor() {
    super();

    this.#canvas = document.createElement("canvas");
    this.#canvas.style = "width: 100%; height: 100%;";

    this.#ctx = this.#canvas.getContext("2d");
  }

  connectedCallback() {
    if (this.#ctx == null) {
      console.error("cannot get the canva's context");
      return;
    }

    this.appendChild(this.#canvas);

    this.painter = new Painter(this.#ctx);

    this.#init();

    addEventListener("resize", this.#init);
  }

  #init = () => {
    if (this.painter == null) {
      console.error("the property draw is null");
      return;
    }

    const boundingBox = this.getBoundingClientRect();

    this.#canvas.width = boundingBox.width;
    this.#canvas.height = boundingBox.height;

    this.painter.init();
  };
}

/** paint the binary digits onto the canvas */
class Painter {
  static mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");

  #ctx;
  #charWidth = 0;
  #charHeight = 0;
  #cntCharsPerLine = 0;
  #cntLines = 0;
  #motionReduced = false;

  #animationRunning = false;

  /** @param {CanvasRenderingContext2D} ctx  */
  constructor(ctx) {
    this.#ctx = ctx;

    Painter.mediaQuery.addEventListener("change", this.init);
  }

  init = () => {
    const margin = 1.5;

    const fontSize = Math.ceil(this.#ctx.canvas.width / 40);

    this.#ctx.font = `${fontSize}px monospace`;
    this.#ctx.fillStyle = "#ccf";

    const textMetrics = this.#ctx.measureText("0");

    this.#charWidth = textMetrics.width * margin;

    this.#charHeight =
      (textMetrics.actualBoundingBoxAscent +
        textMetrics.actualBoundingBoxDescent) *
      margin;

    this.#cntCharsPerLine = Math.ceil(this.#ctx.canvas.width / this.#charWidth);
    this.#cntLines = Math.ceil(this.#ctx.canvas.height / this.#charHeight) + 1;

    this.#motionReduced = Painter.mediaQuery.matches;

    this.#drawAll();

    if (!this.#motionReduced && !this.#animationRunning) {
      this.#animationRunning = true;
      this.#animate();
    } else if (this.#motionReduced && this.#animationRunning) {
      this.#animationRunning = false;
    }
  };

  #drawAll = () => {
    this.#ctx.clearRect(0, 0, this.#ctx.canvas.width, this.#ctx.canvas.height);

    for (let y = 0; y < this.#cntLines; y++) {
      for (let x = 0; x < this.#cntCharsPerLine; x++) {
        this.#ctx.fillText(
          Math.random() > 0.5 ? "1" : "0",
          x * this.#charWidth,
          y * this.#charHeight,
        );
      }
    }
  };

  #animate = async () => {
    if (!this.#animationRunning) return;

    const timeStart = performance.now();

    for (let i = 0; i < 100; i++) {
      const x =
        Math.floor(Math.random() * this.#cntCharsPerLine) * this.#charWidth;
      const y = Math.floor(Math.random() * this.#cntLines) * this.#charHeight;

      this.#ctx.clearRect(x, y, this.#charWidth, -this.#charHeight);
      this.#ctx.fillText(Math.random() > 0.5 ? "0" : "1", x, y);
    }

    const timeDifference = performance.now() - timeStart;
    await new Promise((r) => setTimeout(r, 200 - timeDifference));

    requestAnimationFrame(this.#animate);
  };
}

customElements.define("vs-binary-grid", BinaryGrid);
