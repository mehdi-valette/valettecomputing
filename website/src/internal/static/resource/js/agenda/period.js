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

import { minutesToHours } from "../minutes-to-hours.js";

const template = document.createElement("template");
template.innerHTML = `
  <style>
    .container {
      align-items: center;
      background-color: blue;
      border: 1px solid white;
      box-sizing: border-box;
      color: white;
      cursor: move;
      display: flex;
      flex-direction: column;
      height: 100%;
      justify-content: center;
      position: absolute;
      user-select: none;
      width: 100%;

      .title {
        font-size: 1.2rem;
      }

      .time {
        vertical-align: middle;
      }

      .duration {
        color: lightblue;
        font-size: 0.8rem;
      }
    }
      
  </style>

  <div class="container">
      <div class="title"></div>
      <div>
        <span class="time"></span>
        <span class="duration"></span>
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

  /** @type {string} title of the period */
  #title;

  /** @type {number} the minute at which this period begins */
  #start;

  /** @type {number} the minute at which this period ends */
  #end;

  /** @type {number} the duration in minutes */
  #duration;

  /** @type {HTMLElement} the element that contains the title */
  #titleElement;

  /** @type {HTMLElement} the element that contains the start and end times */
  #timeElement;

  /** @type {HTMLElement} the element that contains the duration */
  #durationElement;

  constructor() {
    super();

    this.attachShadow({ mode: "open" }).append(
      document.importNode(template.content, true)
    );

    this.#container = this.shadowRoot?.querySelector(".container") ?? template;

    this.#titleElement = this.#container.querySelector(".title") ?? template;
    this.#timeElement = this.#container.querySelector(".time") ?? template;
    this.#durationElement =
      this.#container.querySelector(".duration") ?? template;

    this.#dragOffset = 0;
    this.#dragging = false;
    this.#parent = null;

    this.#title = "Massage Classique";
    this.#start = 0;
    this.#end = 90;
    this.#duration = this.#end - this.#start;
  }

  connectedCallback() {
    this.#container.addEventListener("mousedown", (evt) => {
      this.#dragOffset = evt.offsetY;
      this.#dragging = true;
      this.#parent?.appendChild(this);
    });

    this.#timeElement.innerHTML = this.#start.toString();

    document.addEventListener("mousemove", this.#mousemove);

    document.addEventListener("mouseup", () => {
      if (!this.#dragging) return;

      this.#dragOffset = 0;
      this.#dragging = false;
    });

    this.#updateText();
  }

  /** @param {Array<Element>} others */
  checkIntersection = (others) => {
    this.#container.style.backgroundColor = "blue";

    const filtered = others.filter(
      (other) => other !== this && other instanceof VsPeriod
    );

    filtered.forEach((other) => {
      if (!(other instanceof VsPeriod)) return;

      other.checkIntersection(filtered);

      const otherTop = other.#container.getBoundingClientRect().top;
      const otherBottom = other.#container.getBoundingClientRect().bottom;
      const thisTop = this.#container.getBoundingClientRect().top;
      const thisBottom = this.#container.getBoundingClientRect().bottom;

      if (thisBottom <= otherTop || thisTop >= otherBottom) return;

      this.#container.style.backgroundColor = "red";
    });
  };

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

    this.#updatePosition(newOffset, positions);

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
   * adjust the new position and map it to a new start time
   * @param {number} newOffset
   * @param {Positions} pos
   */
  #updatePosition = (newOffset, pos) => {
    if (this.#parent == null) return;

    // adjust the new offset by steps of 15 minutes
    const step = this.#parent.pixelStep * 15;
    newOffset = newOffset - (newOffset % step);

    // if the offset doesn't change, nothing changes
    if (newOffset === this.#container.offsetTop) return;

    // keep the box within the parent box
    if (newOffset < 0) newOffset = 0;

    if (pos.parentTop + newOffset + pos.selfHeight > pos.parentBottom)
      newOffset = pos.parentHeight - pos.selfHeight;

    // calculate the start and end times
    this.#start = Math.round(
      (this.#parent.totalMinutes * newOffset) / pos.parentHeight
    );

    this.#end = Math.round(this.#start + this.#duration);

    // move the element, dispatch the event and update the text
    this.#container.style.top = newOffset + "px";
    this.dispatchEvent(new CustomEvent("moved"));
    this.#updateText();
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

  #updateText = () => {
    this.#titleElement.innerHTML = this.#title;
    this.#timeElement.innerHTML = `${minutesToHours(
      this.#start
    )} - ${minutesToHours(this.#end)}`;

    this.#durationElement.innerHTML = `(${this.#duration}mn)`;
  };
}

customElements.define("vs-period", VsPeriod);
