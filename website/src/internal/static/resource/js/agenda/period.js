/**
 * @typedef Positions
 * @type {object}
 * @property {number} parentBottom
 * @property {number} parentHeight
 * @property {number} parentTop
 * @property {number} selfHeight
 */

/** @typedef {import ("./calendar-day.js").VsCalendarDay} VsCalendarDay */

/** @typedef {"up" | "down" | "none"} Direction */

import { minutesToHours } from "./minutes-to-hours.js";

const template = document.createElement("template");
template.innerHTML = `
  <style>
    .container {
      background-image: linear-gradient(90deg, lightblue 0%, lightblue 5%, blue 5%, blue 100%);
      border: 1px solid white;
      box-sizing: border-box;
      cursor: move;
      position: absolute;
      width: 100%;

      &.intersect {
        background-image: linear-gradient(90deg, red 0%, red 5%, darkred 5%, darkred 100%);
        opacity: .7;
        z-index: 10;
      }

      .content {      
        position: sticky;
        top: 0;
        display: flex;
        padding-left: 10%;
        align-items: center;
        color: white;
        flex-direction: column;
        justify-content: center;
        user-select: none;
        pointer-events: none;
      }

      .title {
        font-size: 1.2rem;
        pointer-events: none;
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
    <div class="content">
      <div class="title"></div>
      <div>
        <span class="time"></span>
        <span class="duration"></span>
      </div>
    </div>
  </div>
`;

export class VsPeriod extends HTMLElement {
  /** @type {number} the offset between the cursor and the top of this custom element */
  #dragOffset;

  /** @type {Direction} */
  #dragDirection;

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
      document.importNode(template.content, true),
    );

    this.#container = this.shadowRoot?.querySelector(".container") ?? template;

    this.#titleElement = this.#container.querySelector(".title") ?? template;
    this.#timeElement = this.#container.querySelector(".time") ?? template;
    this.#durationElement =
      this.#container.querySelector(".duration") ?? template;

    this.#dragOffset = 0;
    this.#dragDirection = "none";
    this.#dragging = false;
    this.#parent = null;

    this.#title = "";
    this.#start = 0;
    this.#end = 0;
    this.#duration = this.#end - this.#start;
  }

  static observedAttributes = ["start", "end", "title"];

  /**
   * @param {string} name
   * @param {string | boolean | null} oldValue
   * @param {string | boolean | null} newValue
   */
  attributeChangedCallback(name, oldValue, newValue) {
    if (oldValue === newValue || typeof newValue !== "string") return;

    switch (name) {
      case "end": {
        if (this.initialized) this.end = Number.parseInt(newValue);
        else {
          this.#end = Number.parseInt(newValue);
          this.#duration = this.#end - this.#start;
        }

        break;
      }
      case "start": {
        if (this.initialized) this.start = Number.parseInt(newValue);
        else {
          this.#start = Number.parseInt(newValue);
          this.#duration = this.#end - this.#start;
        }

        break;
      }
      case "title": {
        this.#title = newValue;
        break;
      }
    }

    this.#redraw();
  }

  connectedCallback() {
    this.#timeElement.innerHTML = this.#start.toString();
    this.#redraw();

    this.#container.addEventListener("mousedown", this.#mousedown);
    document.addEventListener("mousemove", this.#mousemove);
    document.addEventListener("mouseup", this.#mouseup);
  }

  disconnectedCallback() {
    this.#container.removeEventListener("mousedown", this.#mousedown);
    document.removeEventListener("mousemove", this.#mousemove);
    document.removeEventListener("mouseup", this.#mouseup);
  }

  /** @param {Array<Element>} others */
  checkIntersection = (others) => {
    this.#container.classList.remove("intersect");

    const filtered = others.filter(
      (other) => other !== this && other instanceof VsPeriod,
    );

    filtered.forEach((other) => {
      if (!(other instanceof VsPeriod)) return;

      other.checkIntersection(filtered);

      if (this.end <= other.start || this.start >= other.end) return;

      this.#container.classList.add("intersect");
    });
  };

  /** @param {VsCalendarDay} parent */
  init = (parent) => {
    this.#parent = parent;
    this.start = this.#start;
  };

  get end() {
    return this.#end;
  }

  set end(val) {
    this.#end = val;
    this.#duration = this.#end - this.#start;
    this.setAttribute("end", this.#end.toString());
  }

  get initialized() {
    return this.#parent != null && this.#duration > 0;
  }

  get start() {
    return this.#start;
  }

  set start(val) {
    this.#start = val;
    this.end = this.#start + this.#duration;
    this.setAttribute("start", this.#start.toString());
  }

  get title() {
    return this.#title;
  }

  set title(val) {
    this.#title = val;
    this.setAttribute("title", this.#title);
  }

  /** @param {MouseEvent} evt */
  #mousedown = (evt) => {
    this.#dragOffset = evt.offsetY;
    this.#dragging = true;
    this.#parent?.appendChild(this);
  };

  #mouseup = () => {
    if (!this.#dragging) return;

    this.#dragOffset = 0;
    this.#dragging = false;
  };

  /** @param {MouseEvent} evt */
  #mousemove = (evt) => {
    if (!this.#dragging || this.#parent == null) return;

    const positions = this.#getPositions();

    const oldPosition = this.#top;

    let newOffset =
      evt.pageY -
      Math.round(window.pageYOffset) -
      this.#parent.getBoundingClientRect().top -
      this.#dragOffset;

    this.#adjustStart(newOffset, positions);

    const newPosition = this.#top;

    if (newPosition > oldPosition) this.#dragDirection = "down";
    else if (newPosition < oldPosition) this.#dragDirection = "up";
    else this.#dragDirection = "none";

    this.#scrollIntoView();
  };

  #scrollIntoView = () => {
    const selfBox = this.#container.getBoundingClientRect();

    if (selfBox.top < 0 && this.#dragDirection == "up") {
      document.documentElement.scrollBy({ top: selfBox.top });
    }

    if (selfBox.bottom > window.innerHeight && this.#dragDirection == "down") {
      const scrollOffset = selfBox.bottom - window.innerHeight;
      document.documentElement.scrollBy({ top: scrollOffset });
    }
  };

  get #top() {
    return Number.parseInt(this.#container.style.top.replace("px", ""), 10);
  }

  /**
   * adjust the new position and map it to a new start time
   * @param {number} newOffset
   * @param {Positions} pos
   */
  #adjustStart = (newOffset, pos) => {
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

    this.start = Math.round(
      (this.#parent.totalMinutes * newOffset) / pos.parentHeight,
    );
  };

  /** @returns {Positions} */
  #getPositions = () => {
    if (this.#parent == null) throw new Error("the parent must be defined");

    const parentTop = this.#parent.top;
    const parentHeight = this.#parent.height;
    const parentBottom = parentTop + parentHeight;
    const selfHeight = this.#container.scrollHeight;

    return {
      parentBottom,
      parentHeight,
      parentTop,
      selfHeight,
    };
  };

  #redraw = () => {
    this.#titleElement.innerHTML = this.#title;

    this.#timeElement.innerHTML = `${minutesToHours(
      this.#start,
    )} - ${minutesToHours(this.#end)}`;

    this.#durationElement.innerHTML = `(${this.#duration}mn)`;

    requestAnimationFrame(() => {
      if (this.#parent == null) return;

      this.#container.style.top = this.#start * this.#parent.pixelStep + "px";

      this.#container.style.height =
        this.#duration * this.#parent.pixelStep + "px";
    });

    this.dispatchEvent(new CustomEvent("moved"));
  };
}

customElements.define("vs-period", VsPeriod);
