class DigitalRain extends HTMLElement {
  static observedAttributes = ["class"];

  #animationRunning = false;

  /** @type {Array<Snake>} */
  #snakes = [];

  #ctx;

  #canvas;

  constructor() {
    super();

    this.#canvas = document.createElement("canvas");
    this.#ctx = this.#canvas.getContext("2d");
  }

  connectedCallback() {
    this.#canvas.style = "background-color: darkblue;";

    this.appendChild(this.#canvas);

    this.#createSnakes();
    this.#handleMotionReduce();
  }

  /**
   * @param {string} name
   * @param {string} _oldValue
   * @param {string} newValue
   */
  attributeChangedCallback(name, _oldValue, newValue) {
    if (name !== "class") return;
    this.#canvas.className = newValue ?? "";
  }

  #createSnakes = () => {
    if (this.#ctx == null) return;

    const fullHeight = this.#canvas.getBoundingClientRect().height;
    const fullWidth = this.#canvas.getBoundingClientRect().width;

    this.#canvas.width = fullWidth;
    this.#canvas.height = fullHeight;

    this.#ctx.font = "bold .8rem mono";
    this.#ctx.scale(-1, 1);

    const charBox = this.#ctx.measureText("0");
    const charHeight =
      charBox.fontBoundingBoxAscent + charBox.fontBoundingBoxDescent;
    const charWidth = charBox.width;

    for (let i = 0; i < 20; i++) {
      this.#snakes.push(
        new Snake({
          ctx: this.#ctx,
          charHeight,
          charWidth,
          fullHeight,
          fullWidth,
        }),
      );
    }
  };

  #handleMotionReduce = () => {
    const mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");

    if (!mediaQuery.matches) {
      this.#animationRunning = true;
      this.#render();
    }

    mediaQuery.addEventListener("change", () => {
      if (mediaQuery.matches && this.#animationRunning) {
        this.#animationRunning = false;
      } else if (!this.#animationRunning) {
        this.#animationRunning = true;
        this.#render();
      }
    });

    this.motionReduced = mediaQuery.matches;
  };

  #render = async () => {
    if (this.#ctx == null || !this.#animationRunning) return;

    const timeStart = performance.now();

    this.#ctx.clearRect(0, 0, -this.#canvas.width, this.#canvas.height);

    for (const snake of this.#snakes) {
      snake.refresh();
    }

    const timeDifference = performance.now() - timeStart;
    await new Promise((r) => setTimeout(r, 33 - timeDifference));
    requestAnimationFrame(this.#render);
  };
}

class Snake {
  /** @type {Array<string>} */
  #chars = [];

  /** @param {CanvasRenderingContext2D} ctx */
  #ctx;

  #maxLength = 0;
  #charHeight = 0;
  #charWidth = 0;
  #charCountHorizontal = 0;
  #fullWidth = 0;
  #fullHeight = 0;
  #posY = 0;
  #posX = 0;
  #startDelay = 0;
  #updateDelay = 0;
  #maxUpdateDelay = 0;

  /**
   * @param {{ctx: CanvasRenderingContext2D, charHeight: number, charWidth: number, fullWidth: number, fullHeight: number}} params
   */
  constructor({ ctx, charHeight, charWidth, fullHeight, fullWidth }) {
    this.#ctx = ctx;
    this.#charHeight = charHeight;
    this.#charWidth = charWidth;
    this.#fullHeight = fullHeight;
    this.#fullWidth = fullWidth;

    this.#charCountHorizontal = Math.ceil(this.#fullWidth / this.#charWidth);

    this.#reset();
    this.#headStart();
  }

  #reset = () => {
    this.#chars = [];
    this.#maxLength = 10 + Math.floor(Math.random() * 10);
    this.#posX =
      (Math.floor(this.#charCountHorizontal) * Math.random() + 1) *
      this.#charWidth *
      -1;
    this.#posY = 0;
    this.#startDelay = Math.random() * 20;
    this.#maxUpdateDelay = 1 + Math.floor(Math.random() * 2);
    this.#updateDelay = 0;
  };

  #headStart = () => {
    this.#chars = Alphabet.generateArray(
      Math.ceil(Math.random() * this.#maxLength),
    );

    if (this.#chars.length === this.#maxLength)
      this.#posY = Math.floor(
        Math.floor(this.#fullHeight / this.#charHeight) * Math.random(),
      );

    this.#startDelay = 0;

    this.refresh();
  };

  refresh = () => {
    if (this.#posY > this.#fullHeight) this.#reset();

    if (this.#startDelay > 0) {
      this.#startDelay--;
      return;
    }

    this.#updateChars();
    this.#paintChars();
  };

  #updateChars = () => {
    if (this.#updateDelay < this.#maxUpdateDelay) {
      this.#updateDelay++;
      return;
    }

    this.#updateDelay = 0;

    if (this.#chars.length > this.#maxLength) {
      this.#chars.shift();
      this.#posY += this.#charHeight;
    }

    this.#chars.push(Alphabet.pickChar(""));
  };

  #paintChars = () => {
    const lastIndex = this.#chars.length - 1;

    for (const i in this.#chars) {
      const index = Number(i);
      const paintPosY = this.#posY + index * this.#charHeight;
      const currentChar = this.#chars[i];

      this.#ctx.fillStyle =
        index === lastIndex
          ? "white"
          : `rgb(200 200 255 / ${1 - (lastIndex - index) * (1 / this.#maxLength)})`;

      this.#ctx.fillText(currentChar, this.#posX, paintPosY);
    }
  };
}

class Alphabet {
  static #characters = [
    "｡",
    "｢",
    "｣",
    "､",
    "･",
    "ｦ",
    "ｧ",
    "ｨ",
    "ｩ",
    "ｪ",
    "ｫ",
    "ｬ",
    "ｭ",
    "ｮ",
    "ｯ",
    "ｰ",
    "ｱ",
    "ｲ",
    "ｳ",
    "ｴ",
    "ｵ",
    "ｶ",
    "ｷ",
    "ｸ",
    "ｹ",
    "ｺ",
    "ｻ",
    "ｼ",
    "ｽ",
    "ｾ",
    "ｿ",
    "ﾀ",
    "ﾁ",
    "ﾂ",
    "ﾃ",
    "ﾄ",
    "ﾅ",
    "ﾆ",
    "ﾇ",
    "ﾈ",
    "ﾉ",
    "ﾊ",
    "ﾋ",
    "ﾌ",
    "ﾍ",
    "ﾎ",
    "ﾏ",
    "ﾐ",
    "ﾑ",
    "ﾒ",
    "ﾓ",
    "ﾔ",
    "ﾕ",
    "ﾖ",
    "ﾗ",
    "ﾘ",
    "ﾙ",
    "ﾚ",
    "ﾛ",
    "ﾜ",
    "ﾝ",
    "A",
    "B",
    "C",
    "D",
    "E",
    "F",
    "G",
    "H",
    "I",
    "J",
    "K",
    "L",
    "M",
    "N",
    "O",
    "P",
    "Q",
    "R",
    "S",
    "T",
    "U",
    "V",
    "W",
    "X",
    "Y",
    "Z",
    "0",
    "1",
    "2",
    "3",
    "4",
    "5",
    "6",
    "7",
    "8",
    "9",
  ];

  /** @param {string} excludedCharacter */
  static pickChar = (excludedCharacter) => {
    /** @type {string} */
    let character;

    do {
      const index = Math.floor(Math.random() * this.#characters.length);
      character = this.#characters[index];
    } while (character === excludedCharacter);

    return character;
  };

  /** @param {number} size */
  static generateArray = (size) => {
    const chars = [];

    for (let i = 0; i < size; i++) {
      chars.push(this.pickChar(""));
    }

    return chars;
  };
}

customElements.define("vs-digital-rain", DigitalRain);
