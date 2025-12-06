/**
 * @typedef Positions
 * @type {object}
 * @property {number} documentTop
 * @property {number} parentBottom
 * @property {number} parentHeight
 * @property {number} parentTop
 * @property {number} selfHeight
 */

/** @typedef {import ("./calendar-day.js").VsCalendarDay} VsCalendarDay */

const template = document.createElement("template");
template.innerHTML = `
  <style>
    .outer {
      position: absolute;
      width: 100%;
      cursor: move;
      user-select: none;
    }

    .inner {
      display: flex;
      align-items: center;
      color: white;
      justify-content: center;
      background-color: blue;
      width: 100%;
      height: 100%;
    }
  </style>

  <div class="outer">
    <div class="inner">
      <span data-start=""></span>
    </div>
  </div>
`;

export class VsPeriod extends HTMLElement {
  /** @type {number} the offset between the cursor and the top of this custom element */
  #dragOffset;

  /** @type {boolean} is true when this custom element is being dragged */
  #dragging;

  /** @type {VsCalendarDay | null} the parent element of this custom element */
  #parent;

  /** @type {HTMLElement} the main DIV for this custom element */
  #container;

  /** @type {HTMLElement} the element that contains the start time */
  #startElement;

  /** @type {number} the minute at which this period begins */
  #start;

  /** @type {number} the minute at which this period ends */
  #end;

  /** @type {number} the duration in minutes */
  #duration;

  constructor() {
    super();

    this.attachShadow({ mode: "open" }).append(
      document.importNode(template.content, true)
    );

    this.#container = this.shadowRoot?.querySelector(".outer") ?? template;

    this.#startElement =
      this.#container.querySelector("[data-start]") ?? template;

    this.#dragOffset = 0;
    this.#dragging = false;
    this.#parent = null;

    this.#start = 0;
    this.#end = 90;
    this.#duration = this.#end - this.#start;
  }

  connectedCallback() {
    this.#container.addEventListener("mousedown", (evt) => {
      this.#dragOffset = evt.offsetY;
      this.#dragging = true;
    });

    this.#startElement.innerHTML = this.#start.toString();

    document.addEventListener("mousemove", this.#mousemove);

    document.addEventListener("mouseup", () => {
      this.#dragOffset = 0;
      this.#dragging = false;
    });
  }

  /** @param {VsCalendarDay} parent */
  setParent = (parent) => {
    this.#parent = parent;

    this.#container.style.height =
      (this.#parent.height / this.#parent.totalMinutes) * this.#duration + "px";
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

    if (selfBox.top < 0) {
      document.documentElement.scrollBy({ top: selfBox.top });
    }

    if (selfBox.bottom > window.innerHeight) {
      const scrollOffset = selfBox.bottom - window.innerHeight;
      document.documentElement.scrollBy({ top: scrollOffset });
    }
  };

  /**
   * map the new position to a new start time
   * @param {number} newOffset
   * @param {Positions} pos
   */
  #updatePeriodBeginning = (newOffset, pos) => {
    if (this.#parent == null) return;

    // adjust the new offset by steps of 15 minutes
    const step = this.#parent.pixelStep * 15;
    newOffset = newOffset - (newOffset % step);

    // if the offset doesn't change, nothing changes
    if (newOffset === this.#container.offsetTop) return;

    // bound the box to the parent box
    if (newOffset < 0) newOffset = 0;

    if (pos.parentTop + newOffset + pos.selfHeight > pos.parentBottom)
      newOffset = pos.parentHeight - pos.selfHeight;

    // calculate and show the start and end times
    this.#start = Math.round(
      (this.#parent.totalMinutes * newOffset) / pos.parentHeight
    );

    this.#end = Math.round(this.#start + this.#duration);

    this.#startElement.innerHTML = this.#start.toString();

    // move the element
    this.#container.style.top = newOffset + "px";
  };

  /** @returns {Positions} */
  #getPositions = () => {
    if (this.#parent == null) throw new Error("the parent must be defined");

    const documentTop = document.documentElement.scrollTop;
    const parentTop = this.#parent.top;
    const parentHeight = this.#parent.height;
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
