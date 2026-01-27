/**
 * @description shows lines of 1s and 0s going from right to left, as if someone was typing
 */

export class VsRtlCode extends HTMLCanvasElement {
  /**
   * @type {Array<{text: string, posY: number}>}
   */
  lines = [];

  connectedCallback() {
    this.style = "position: fixed; width: 100%; height: 100%; z-index: -1";
    this.width = this.clientWidth;
    this.height = this.clientHeight;

    const ctx = this.getContext("2d");
    if (ctx == null) return;

    const characterCount = this.clientWidth / ctx.measureText("01").width / 2;

    ctx.font = "5rem mono";
    ctx.fillStyle = "lightblue";
    const lineHeight = ctx.measureText("01").emHeightAscent * 2;

    let currentY = lineHeight;

    for (let i = 0; currentY <= this.clientHeight + lineHeight; i++) {
      this.lines[i] = {
        text: this.createRandomLine(characterCount),
        posY: currentY,
      };

      currentY += lineHeight;

      setInterval(
        () => this.line(ctx, i, lineHeight),
        Math.floor(Math.random() * 50 + 300),
      );
    }
  }

  /**
   * @param {number} characterCount
   */
  createRandomLine = (characterCount) => {
    let newText = "";

    for (let i = 0; i < characterCount; i++) {
      newText += Math.random() > 0.5 ? "1" : "0";
    }

    return newText;
  };

  /**
   * @param {CanvasRenderingContext2D} ctx
   * @param {number} lineIndex
   * @param {number} lineHeight
   */
  line = (ctx, lineIndex, lineHeight) => {
    let { text, posY } = this.lines[lineIndex];

    this.lines[lineIndex].text =
      text.substring(1) + (Math.random() > 0.5 ? "1" : "0");

    requestAnimationFrame(() => {
      ctx.clearRect(0, posY - lineHeight, this.clientWidth, lineHeight + 10);
      ctx.fillText(text, 10, posY);
    });
  };
}

customElements.define("vs-rtl-code", VsRtlCode, { extends: "canvas" });
