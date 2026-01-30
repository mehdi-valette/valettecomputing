/**
 * @description shows lines of 1s and 0s going from right to left, as if someone was typing
 */

class VsBinaryRtl extends HTMLCanvasElement {
  /** @type {Array<{text: string, posY: number; interval: number | null}>} */
  #lines = [];

  /** @type {CanvasRenderingContext2D | null} */
  #ctx = null;

  #lineHeight = 0;
  #charsPerLine = 0;
  #motionReduced = false;

  connectedCallback() {
    this.style = "position: fixed; width: 100%; height: 100%; z-index: -1;";
    this.#ctx = this.getContext("2d");

    this.#handleResize();
    this.#handleMotionReduce();
    this.#refreshAllLines();
  }

  #handleResize = () => {
    addEventListener("resize", () => {
      this.#setSizes();
      this.#refreshAllLines();
    });

    this.#setSizes();
  };

  #setSizes = () => {
    if (this.#ctx == null) return;

    this.width = this.clientWidth;
    this.height = this.clientHeight;

    // it looks like changing the width and height resets the context
    this.#ctx.font = "5rem mono";
    this.#ctx.fillStyle = "#ccf";

    this.#charsPerLine =
      (this.clientWidth / this.#ctx.measureText("01").width) * 2;

    this.#lineHeight = this.#ctx.measureText("01").emHeightAscent * 2;
  };

  #handleMotionReduce = () => {
    const mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");

    mediaQuery.addEventListener("change", () => {
      this.#motionReduced = mediaQuery.matches;
      this.#refreshAllLines();
    });

    this.#motionReduced = mediaQuery.matches;
  };

  #refreshAllLines = () => {
    if (this.#ctx == null) return;

    for (const l of this.#lines)
      if (l.interval != null) clearInterval(l.interval);

    this.#lines = [];

    let posY = this.#lineHeight;

    for (let i = 0; posY <= this.clientHeight + this.#lineHeight; i++) {
      this.#initLine(i, posY);

      posY += this.#lineHeight;
    }
  };

  /**
   * @param {number} lineIndex
   * @param {number} posY
   */
  #initLine = (lineIndex, posY) => {
    let text = "";

    for (let i = 0; i < this.#charsPerLine; i++) {
      text += Math.random() > 0.5 ? "1" : "0";
    }

    this.#lines[lineIndex] = { text, posY, interval: null };

    if (!this.#motionReduced) {
      this.#lines[lineIndex].interval = setInterval(
        () => this.#updateLine(lineIndex),
        Math.floor(Math.random() * 100 + 200),
      );
    }

    this.#drawLine(lineIndex);
  };

  /**
   * @param {number} lineIndex
   */
  #updateLine = (lineIndex) => {
    let { text, posY } = this.#lines[lineIndex];

    this.#lines[lineIndex].text =
      text.substring(1) + (Math.random() > 0.5 ? "1" : "0");

    this.#drawLine(lineIndex);
  };

  /**
   * @param {number} lineIndex
   */
  #drawLine = (lineIndex) => {
    if (this.#ctx == null) return;

    const { text, posY } = this.#lines[lineIndex];

    this.#ctx.clearRect(
      0,
      posY - this.#lineHeight,
      this.clientWidth,
      this.#lineHeight + 10,
    );

    this.#ctx.fillText(text, 10, posY);
  };
}

customElements.define("binary-rtl", VsBinaryRtl, { extends: "canvas" });
