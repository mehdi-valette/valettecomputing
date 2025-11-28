/**
 * @typedef Positions
 * @type {object}
 * @property {number} documentTop
 * @property {number} parentBottom
 * @property {number} parentHeight
 * @property {number} parentTop
 * @property {number} selfHeight
 */

class VsPeriod extends HTMLElement {
  /** @type {number} the offset between the cursor and the top of this custom element */
  #dragOffset;

  /** @type {boolean} is true when this custom element is being dragged */
  #dragging;

  /** @type {HTMLElement | null} the parent element of this custom element */
  #parent;

  /** @type {ShadowRoot} element's shadow root */
  #shadow;

  /** @type {HTMLElement} the main DIV for this custom element */
  #container;

  /** @type {number} the minute at which this period begins */
  #start;

  /** @type {number} the minute at which this period ends */
  #end;

  /** @type {number} the duration in minutes */
  #duration;

  /** @type {number} total number of minutes in the parent container */
  #totalMinutes;

  constructor() {
    super();

    this.#shadow = this.attachShadow({ mode: "open" });

    this.#container = document.createElement("div");
    this.#container.style.position = "absolute";
    this.#container.style.width = "100%";
    this.#container.style.cursor = "move";
    this.#container.style.userSelect = "none";

    this.#dragOffset = 0;
    this.#dragging = false;
    this.#parent = null;

    this.#start = 0;
    this.#end = 90;
    this.#duration = this.#end - this.#start;
    this.#totalMinutes = 1440;

    this.#shadow.append(this.#container);
  }

  connectedCallback() {
    this.#container.innerHTML = `
      <div style="display: flex; align-items: center; color: white; justify-content: center; background-color: blue; width: 100%; height: 100%;">
        <span data-start="">${this.#start}</span>
      </div>
    `;

    this.#container.addEventListener("mousedown", (evt) => {
      this.#dragOffset = evt.offsetY;
      this.#dragging = true;
    });

    document.addEventListener("mousemove", this.#mousemove);

    document.addEventListener("mouseup", () => {
      this.#dragOffset = 0;
      this.#dragging = false;
    });
  }

  /** @param {HTMLElement} parent */
  setParent = (parent) => {
    this.#parent = parent;

    this.#container.style.height =
      (this.#parent.clientHeight / this.#totalMinutes) * this.#duration + "px";
  };

  /** @param {MouseEvent} evt */
  #mousemove = (evt) => {
    if (!this.#dragging || this.#parent == null) return;

    const positions = this.#getPositions();

    let newOffset =
      evt.screenY +
      positions.documentTop -
      positions.selfHeight -
      positions.parentTop -
      this.#dragOffset;

    this.#updatePeriodBeginning(newOffset, positions);

    this.#scrollIntoView();
  };

  #scrollIntoView = () => {
    const selfBox = this.#container.getBoundingClientRect();

    if (selfBox.bottom > window.innerHeight) {
      const scroll = selfBox.bottom - window.innerHeight;
      document.documentElement.scrollBy({ top: scroll });
    }

    if (selfBox.top < 0) {
      document.documentElement.scrollBy({ top: selfBox.top });
    }
  };

  /**
   * @param {number} newOffset
   * @param {Positions} pos
   */
  #updatePeriodBeginning = (newOffset, pos) => {
    // adjust the new offset by step of 15 minutes
    const step = (pos.parentHeight * 15) / this.#totalMinutes;
    newOffset = newOffset - (newOffset % step);

    // offset doesn't change, nothing changes
    if (newOffset === this.#container.offsetTop) return;

    if (newOffset < 0) newOffset = 0;

    if (pos.parentTop + newOffset + pos.selfHeight > pos.parentBottom)
      newOffset = pos.parentHeight - pos.selfHeight;

    this.#start = Math.round(
      (this.#totalMinutes * newOffset) / pos.parentHeight
    );
    this.#end = Math.round(this.#start + this.#duration);
    const containerStart = this.#container.querySelector("[data-start]");

    if (containerStart != null)
      containerStart.innerHTML = this.#start.toString();

    this.#container.style.top = newOffset + "px";
  };

  /** @returns {Positions} */
  #getPositions = () => {
    if (this.#parent == null) throw new Error("the parent must be defined");

    const documentTop = document.documentElement.scrollTop;
    const parentTop = this.#parent.getBoundingClientRect().top + documentTop;
    const parentHeight = this.#parent.clientHeight;
    const parentBottom = parentTop + parentHeight;
    const selfHeight = this.#container.scrollHeight;

    return {
      documentTop,
      parentBottom,
      parentHeight,
      parentTop,
      selfHeight,
    };
  };
}

customElements.define("vs-period", VsPeriod);
